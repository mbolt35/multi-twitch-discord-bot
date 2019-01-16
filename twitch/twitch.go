package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	httputil "github.com/mbolt35/multi-twitch-discord-bot/util"
)

const (
	// TwitchUrl is the base url for Twitch
	TwitchUrl string = "http://twitch.tv"

	// TwitchUserNameToUserIdUrl is the API url for Twitch User Id Lookup
	TwitchUserNameToUserIdUrl string = "https://api.twitch.tv/kraken/users"

	// TwitchWebhookUrl is the webhook subscription api Twitch WebHooks Url
	TwitchWebhookUrl string = "https://api.twitch.tv/helix/webhooks/hub"

	// TwitchStreamsTopicUrl Twitch WebHook Topic Url
	TwitchStreamsTopicUrl string = "https://api.twitch.tv/helix/streams"

	// TwitchV5 Twitch V5 API type
	TwitchV5 string = "application/vnd.twitchtv.v5+json"

	// TwitchHubChallengeQueryParameter Webhook Challenge Query Parameter
	TwitchHubChallengeQueryParameter string = "hub.challenge"

	// TwitchHubLeaseQueryParameter Webhook Lease Time Query Parameter
	TwitchHubLeaseQueryParameter string = "hub.lease_seconds"

	// TwitchHubModeQueryParameter Webhook Mode Query Parameter
	TwitchHubModeQueryParameter string = "hub.mode"

	// TwitchHubTopicQueryParameter Webhook Topic Query Parameter
	TwitchHubTopicQueryParameter string = "hub.topic"

	// TwitchHubReasonQueryParameter Webhook Reason Query Parameter
	TwitchHubReasonQueryParameter string = "hub.reason"

	// TwitchModeDenied Twitch Subscribe Request denied
	TwitchModeDenied string = "denied"

	// TwitchModeSubscribe Mode option for twitch web hook
	TwitchModeSubscribe string = "subscribe"

	// TwitchModeUnsubscribe Mode option for twitch web hook
	TwitchModeUnsubscribe string = "unsubscribe"

	// TwitchUserNameQueryParameter User Name Url Query Parameter
	TwitchUserNameQueryParameter string = "user_login"

	// TwitchUserIdQueryParameter UserId Query Parameter
	TwitchUserIdQueryParameter string = "user_id"

	// TwitchUserNameToUserIdQueryParameter User Name to User Id Query Parameter
	TwitchUserNameToUserIdQueryParameter string = "login"

	// TwitchMaxLeaseSeconds Maximum Lease Time for Subscriptions
	TwitchMaxLeaseSeconds int = 864000
)

// Request Payload from Client to Twitch Requesting Notifications
type TwitchWebhookPayload struct {
	CallbackUrl  string `json:"hub.callback"`
	Mode         string `json:"hub.mode"`
	Topic        string `json:"hub.topic"`
	LeaseSeconds int    `json:"hub.lease_seconds,omitempty"`
	Secret       string `json:"hub.secret,omitempty"`
}

// Post Payload from Twitch to Client for a Notification
type TwitchNotification struct {
	Id           string   `json:"id,omitempty"`
	UserId       string   `json:"user_id,omitempty"`
	UserName     string   `json:"user_name,omitempty"`
	GameId       string   `json:"game_id,omitempty"`
	CommunityIds []string `json:"community_ids,omitempty"`
	Type         string   `json:"type,omitempty"`
	Title        string   `json:"title,omitempty"`
	ViewerCount  int      `json:"viewer_count,omitempty"`
	StartedAt    string   `json:"started_at,omitempty"`
	Language     string   `json:"language,omitempty"`
	ThumbnailUrl string   `json:"thumbnail_url,omitempty"`
	TagIds       []string `json:"tag_ids,omitempty"`
}

// Post Wrapper for Notification Payloads
type TwitchNotificationPayload struct {
	Notifications []TwitchNotification `json:"data"`
}

// TwitchUser representation from Querying user info endpoint
type TwitchUser struct {
	UserId      string `json:"_id"`
	UserName    string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LogoUrl     string `json:"logo"`
}

// Twitch user endpoint payload
type TwitchUsersPayload struct {
	Total int          `json:"_total"`
	Users []TwitchUser `json:"users"`
}

type twitch struct {
	userNameCache map[string]string
	clientId      string
}

// TwitchClient is the interface used to represent a client capable of communicating with twitch.tv apis.
type TwitchClient interface {
	FromUserId(userId string) string
	UserIdsFor(userNames []string) ([]string, error)
	SubscribeToStreams(notifyEndPoint string, userIds []string)
}

// NewTwitch creates a new TwitchClient implementation and returns it
func NewTwitch(clientId string) TwitchClient {
	instance := twitch{
		userNameCache: make(map[string]string),
		clientId:      clientId,
	}

	return &instance
}

// UserStreamUrl returns the url of the live stream for a specific user
func UserStreamUrl(userName string) string {
	return TwitchUrl + "/" + userName
}

// FromUserId looks up a single user id from internal cache
func (t *twitch) FromUserId(userId string) string {
	return t.userNameCache[userId]
}

// UserIdsFor converts user names into a comma delimited string of user ids
func (t *twitch) UserIdsFor(userNames []string) ([]string, error) {
	userIds := []string{}

	if len(userNames) == 0 {
		return userIds, errors.New("UserNames is Length: 0")
	}

	request, err := http.NewRequest(http.MethodGet, getUserConversionUrl(userNames), nil)
	if nil != err {
		return userIds, err
	}

	request.Header.Set(httputil.HttpAcceptHeader, TwitchV5)
	request.Header.Set(httputil.HttpClientIdHeader, t.clientId)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if nil != err {
		return userIds, err
	}

	var payload TwitchUsersPayload
	e := httputil.DecodeJson(resp.Body, &payload)
	if nil != e {
		return userIds, e
	}

	for _, twitchUser := range payload.Users {
		userIds = append(userIds, twitchUser.UserId)

		// Cache Display Name for User Id
		t.userNameCache[twitchUser.UserId] = twitchUser.DisplayName
	}

	return userIds, nil
}

// Sends a Subscribe Request for Go Live Events for the Provided Users
func (t *twitch) SubscribeToStreams(notifyEndPoint string, userIds []string) {
	if len(userIds) == 0 {
		return
	}

	for _, userId := range userIds {
		topicUrl := getStreamTopicUrl(userId)

		payload := TwitchWebhookPayload{
			CallbackUrl:  notifyEndPoint,
			Mode:         TwitchModeSubscribe,
			Topic:        topicUrl,
			LeaseSeconds: TwitchMaxLeaseSeconds,
		}

		jsonBytes, err := httputil.EncodeJson(payload)
		if err != nil {
			log.Fatalln(err)
		}

		request, err := http.NewRequest(http.MethodPost, TwitchWebhookUrl, bytes.NewBuffer(jsonBytes))
		if err != nil {
			log.Fatalln(err)
		}

		request.Header.Set(httputil.HttpContentTypeHeader, httputil.JsonContentType)
		request.Header.Set(httputil.HttpClientIdHeader, t.clientId)

		httpClient := &http.Client{}
		resp, err := httpClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
	}
}

// Gets the Stream Topic URL
func getStreamTopicUrl(userId string) string {
	u, _ := url.Parse(TwitchStreamsTopicUrl)
	q := u.Query()
	q.Add(TwitchUserIdQueryParameter, userId)
	u.RawQuery = q.Encode()

	return u.String()
}

// Gets the user name to user id conversion url
func getUserConversionUrl(userNames []string) string {
	users := strings.Join(userNames, ",")

	u, _ := url.Parse(TwitchUserNameToUserIdUrl)
	q := u.Query()
	q.Add(TwitchUserNameToUserIdQueryParameter, users)
	u.RawQuery = q.Encode()

	return u.String()
}
