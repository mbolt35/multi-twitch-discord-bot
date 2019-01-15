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
	// The base url of the host for constructing callback urls
	HostUrlEnvVar string = "HOST_URL"

	// Heroku injects the exposed port via environment variable
	HostPortEnvVar string = "PORT"

	// The Twitch App Client Identifier used when communicating with Twitch APIs
	ClientIdEnvVar string = "TWITCH_CLIENT_ID"

	// A comma delimited list of Twitch user names to subscribe to go live events for
	UsersEnvVar string = "TWITCH_USERS"

	// Default Port Value when running locally
	DefaultPort string = "3001"

	// JSON Content-Type
	JsonContentType string = "application/json"

	// The end point we'll bind to for receiving http requests
	NotifyEndPoint string = "notify"

	// Twitch WebHooks Url
	TwitchWebhookUrl string = "https://api.twitch.tv/helix/webhooks/hub"

	// Twitch WebHook Topic Url
	TwitchStreamsTopicUrl string = "https://api.twitch.tv/helix/streams"

	// Mode option for twitch web hook
	TwitchModeSubscribe string = "subscribe"

	// Mode option for twitch web hook
	TwitchModeUnsubscribe string = "unsubscribe"

	// User Name Url Query Parameter
	TwitchUserNameQueryParameter string = "user_login"
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

type TwitchNotificationPayload struct {
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
		host = "http://localhost"
	}

	return host
}

// Gets the port the web app should be hosted on
func GetHostPort() string {
	port := os.Getenv(HostPortEnvVar)

	if port == "" {
		log.Println("$PORT not set. Defaulting to 3001")
		port = DefaultPort
	}

	return port
}

// Gets the name of twitch users to listen for go live events
func GetUserNames() []string {
	userNames := os.Getenv(UsersEnvVar)

	if "" == userNames {
		return []string{}
	}

	return strings.Split(userNames, ",")
}

func GetStreamTopicUrl(users []string) string {
	topicUrl := TwitchStreamsTopicUrl
	firstUser := true

	for _, user := range users {
		log.Println("Listening for User: " + user)

		if firstUser {
			topicUrl += "?"
			firstUser = false
		} else {
			topicUrl += "&"
		}

		topicUrl += TwitchUserNameQueryParameter + "=" + user
	}

	log.Println("Topic URL: " + topicUrl)
	return topicUrl
}

// Sends a Subscribe Request for Go Live Events for the Provided Users
func SubscribeToGoLiveEvents(users []string) {
	if len(users) == 0 {
		return
	}

	topicUrl := GetStreamTopicUrl(users)
	payload := TwitchWebhookPayload{
		CallbackUrl: GetHostUrl() + "/" + NotifyEndPoint,
		Mode:        TwitchModeSubscribe,
		Topic:       topicUrl,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s\n", jsonBytes)

	resp, err := http.Post(TwitchWebhookUrl, JsonContentType, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
}

// Handles Incoming Twitch Notifications
func OnTwitchNotification(responseWriter http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var payload []TwitchNotificationPayload
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	for _, notification := range payload {
		log.Println("Notification [Name=" + notification.UserName + ", Status: " + notification.Type + ", Title: " + notification.Title + "]")
	}
}

// Initializes HTTP End Points
func InitializeEndPoints() {
	http.HandleFunc("/"+NotifyEndPoint, OnTwitchNotification)

	go func() {
		err := http.ListenAndServe(":"+GetHostPort(), nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
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
