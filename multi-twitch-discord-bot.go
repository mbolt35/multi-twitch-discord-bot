package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	// Content Type Request Header Key
	HttpContentTypeHeader string = "Content-Type"

	// Client Id Request Header Key
	HttpClientIdHeader string = "Client-ID"

	// Accept Request Header Key
	HttpAcceptHeader string = "Accept"

	// The base url of the host for constructing callback urls
	HostUrlEnvVar string = "HOST_URL"

	// Heroku injects the exposed port via environment variable
	HostPortEnvVar string = "PORT"

	// The Twitch App Client Identifier used when communicating with Twitch APIs
	ClientIdEnvVar string = "TWITCH_CLIENT_ID"

	// A comma delimited list of Twitch user names to subscribe to go live events for
	UsersEnvVar string = "TWITCH_USERS"

	// The default host url
	DefaultHostUrl string = "http://localhost"

	// Default Port Value when running locally
	DefaultPort string = "3001"

	// JSON Content-Type
	JsonContentType string = "application/json"

	// Twitch V5 API type
	TwitchV5 string = "application/vnd.twitchtv.v5+json"

	// The end point we'll bind to for receiving http requests
	NotifyEndPoint string = "notify"

	// Twitch User Id Lookup
	TwitchUserNameToUserIdUrl string = "https://api.twitch.tv/kraken/users"

	// Twitch WebHooks Url
	TwitchWebhookUrl string = "https://api.twitch.tv/helix/webhooks/hub"

	// Twitch WebHook Topic Url
	TwitchStreamsTopicUrl string = "https://api.twitch.tv/helix/streams"

	// Webhook Challenge Query Parameter
	TwitchHubChallengeQueryParameter string = "hub.challenge"

	// Webhook Lease Time Query Parameter
	TwitchHubLeaseQueryParameter string = "hub.lease_seconds"

	// Webhook Mode Query Parameter
	TwitchHubModeQueryParameter string = "hub.mode"

	// Webhook Topic Query Parameter
	TwitchHubTopicQueryParameter string = "hub.topic"

	// Webhook Reason Query Parameter
	TwitchHubReasonQueryParameter string = "hub.reason"

	// Twitch Subscribe Request denied
	TwitchModeDenied string = "denied"

	// Mode option for twitch web hook
	TwitchModeSubscribe string = "subscribe"

	// Mode option for twitch web hook
	TwitchModeUnsubscribe string = "unsubscribe"

	// User Name Url Query Parameter
	TwitchUserNameQueryParameter string = "user_login"

	// UserId Query Parameter
	TwitchUserIdQueryParameter string = "user_id"

	// User Name to User Id Query Parameter
	TwitchUserNameToUserIdQueryParameter string = "login"

	// Maximum Lease Time for Subscriptions
	TwitchMaxLeaseSeconds int = 864000
)

var (
	twitchClientId  string
	twitchUserNames []string
	hostUrl         string
	hostPort        string
)

// hub.callback       string  URL where notifications will be delivered.
// hub.mode           string  Type of request. Valid values: subscribe, unsubscribe.
// hub.topic          string  URL for the topic to subscribe to or unsubscribe from. topic maps to a new Twitch API endpoint.
// hub.lease_seconds  int     Number of seconds until the subscription expires. Default: 0. Maximum: 864000.
// hub.secret         string  Secret used to sign notification payloads.
type TwitchWebhookPayload struct {
	CallbackUrl  string `json:"hub.callback"`
	Mode         string `json:"hub.mode"`
	Topic        string `json:"hub.topic"`
	LeaseSeconds int    `json:"hub.lease_seconds,omitempty"`
	Secret       string `json:"hub.secret,omitempty"`
}

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

type TwitchNotificationPayload struct {
	Notifications []TwitchNotification `json:"data"`
}

type TwitchUser struct {
	UserId    string `json:"_id"`
	UserName  string `json:"name"`
	Type      string `json:"type"`
	Bio       string `json:"bio"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	LogoUrl   string `json:"logo"`
}

type TwitchUsersPayload struct {
	Total int          `json:"_total"`
	Users []TwitchUser `json:"users"`
}

// Debug Function to Dump All Environment Variables to stdout
func DumpEnvironmentVariables() {
	fmt.Println("--- ENV Vars ---")
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0] + " = " + os.Getenv(pair[0]))
	}
	fmt.Println("---------------")
}

// Gets the Host URL base
func GetHostUrl() string {
	host := os.Getenv(HostUrlEnvVar)

	if host == "" {
		host = DefaultHostUrl
	}

	return host
}

// Gets the port the web app should be hosted on
func GetHostPort() string {
	if "" != hostPort {
		return hostPort
	}

	hostPort = os.Getenv(HostPortEnvVar)

	if hostPort == "" {
		log.Println("$PORT not set. Defaulting to 3001")
		hostPort = DefaultPort
	}

	return hostPort
}

// Gets the Client Identifier Header for HTTP Requests
func GetClientId() string {
	if "" != twitchClientId {
		return twitchClientId
	}

	twitchClientId = os.Getenv(ClientIdEnvVar)
	return twitchClientId
}

// Gets the name of twitch users to listen for go live events
func GetUserNames() []string {
	if nil != twitchUserNames {
		return twitchUserNames
	}

	userNames := os.Getenv(UsersEnvVar)
	if "" != userNames {
		twitchUserNames = strings.Split(userNames, ",")
	} else {
		twitchUserNames = []string{}
	}

	return twitchUserNames
}

// Accepts a parameter and value and returns the full query parameter
func ToQueryParameter(p string, v string, first bool) string {
	var qp string
	if first {
		qp = "?"
	} else {
		qp = "&"
	}

	return qp + p + "=" + v
}

// Converts user names into a comma delimited string of user ids
func ToUserIds(userNames []string) []string {
	if len(userNames) == 0 {
		return []string{}
	}

	users := strings.Join(userNames, ",")
	url := TwitchUserNameToUserIdUrl + ToQueryParameter(TwitchUserNameToUserIdQueryParameter, users, true)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if nil != err {
		panic(err)
	}

	request.Header.Set(HttpAcceptHeader, TwitchV5)
	request.Header.Set(HttpClientIdHeader, GetClientId())

	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	var payload TwitchUsersPayload
	decoder := json.NewDecoder(resp.Body)
	e := decoder.Decode(&payload)
	if e != nil {
		panic(e)
	}

	userIds := []string{}
	for _, twitchUser := range payload.Users {
		log.Println("UserId: " + twitchUser.UserId)
		userIds = append(userIds, twitchUser.UserId)
	}
	return userIds
}

func GetStreamTopicUrl(userId string) string {
	return TwitchStreamsTopicUrl + ToQueryParameter(TwitchUserIdQueryParameter, userId, true)
}

func EncodeJson(obj interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(obj)

	return buffer.Bytes(), err
}

// Sends a Subscribe Request for Go Live Events for the Provided Users
func SubscribeToGoLiveEvents(users []string) {
	if len(users) == 0 {
		return
	}

	userIds := ToUserIds(users)
	for _, userId := range userIds {
		topicUrl := GetStreamTopicUrl(userId)

		payload := TwitchWebhookPayload{
			CallbackUrl:  GetHostUrl() + "/" + NotifyEndPoint,
			Mode:         TwitchModeSubscribe,
			Topic:        topicUrl,
			LeaseSeconds: TwitchMaxLeaseSeconds,
		}

		jsonBytes, err := EncodeJson(payload)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%s\n", jsonBytes)

		request, err := http.NewRequest(http.MethodPost, TwitchWebhookUrl, bytes.NewBuffer(jsonBytes))
		if err != nil {
			log.Fatalln(err)
		}

		request.Header.Set(HttpContentTypeHeader, JsonContentType)
		request.Header.Set(HttpClientIdHeader, GetClientId())

		httpClient := &http.Client{}
		resp, err := httpClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		log.Println(result)
	}
}

// Handles Incoming Twitch Notifications
func OnTwitchNotification(rw http.ResponseWriter, request *http.Request) {
	if http.MethodGet == request.Method {
		fmt.Println("Received GET")

		q := request.URL.Query()

		mode := q.Get(TwitchHubModeQueryParameter)
		topic := q.Get(TwitchHubTopicQueryParameter)

		if TwitchModeDenied == mode {
			reason := q.Get(TwitchHubReasonQueryParameter)
			fmt.Println("Failed to Subscribe to Webhook: " + reason)
			rw.WriteHeader(http.StatusOK)
			return
		}

		challenge := q.Get(TwitchHubChallengeQueryParameter)
		lease := q.Get(TwitchHubLeaseQueryParameter)

		fmt.Printf("Challenge: %s\nLease: %s\nMode: %s\nTopic: %s\n", challenge, lease, mode, topic)

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(challenge))
		return
	}

	if http.MethodPost == request.Method {
		fmt.Println("Received POST")
		decoder := json.NewDecoder(request.Body)
		payload := TwitchNotificationPayload{}
		err := decoder.Decode(&payload)
		if err != nil {
			panic(err)
		}

		for _, notification := range payload.Notifications {
			fmt.Println("Notification [Name=" + notification.UserName + ", Status: " + notification.Type + ", Title: " + notification.Title + "]")
		}
	}
}

// Initializes HTTP End Points
func InitializeEndPoints() {
	http.HandleFunc("/"+NotifyEndPoint, OnTwitchNotification)

	go func() {
		err := http.ListenAndServe(":"+GetHostPort(), nil)
		if nil != err {
			panic(err)
		}
	}()
}

func main() {
	DumpEnvironmentVariables()
	InitializeEndPoints()

	users := GetUserNames()
	SubscribeToGoLiveEvents(users)

	select {}
}
