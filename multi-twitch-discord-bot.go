package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mbolt35/multi-twitch-discord-bot/discord"
	"github.com/mbolt35/multi-twitch-discord-bot/settings"
	"github.com/mbolt35/multi-twitch-discord-bot/twitch"
)

const (
	// The end point we'll bind to for receiving http requests
	NotifyEndPoint string = "notify"

	// The end point used for keep alive
	PingEndPoint string = "ping"
)

var twitchClient twitch.TwitchClient

// Handles Incoming Twitch Notifications
func OnTwitchNotification(rw http.ResponseWriter, request *http.Request) {
	if http.MethodGet == request.Method {
		fmt.Println("Received GET")

		q := request.URL.Query()

		mode := q.Get(twitch.TwitchHubModeQueryParameter)
		topic := q.Get(twitch.TwitchHubTopicQueryParameter)

		if twitch.TwitchModeDenied == mode {
			reason := q.Get(twitch.TwitchHubReasonQueryParameter)
			fmt.Println("Failed to Subscribe to Webhook: " + reason)
			rw.WriteHeader(http.StatusOK)
			return
		}

		challenge := q.Get(twitch.TwitchHubChallengeQueryParameter)
		lease := q.Get(twitch.TwitchHubLeaseQueryParameter)

		fmt.Printf("Challenge: %s\nLease: %s\nMode: %s\nTopic: %s\n", challenge, lease, mode, topic)

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(challenge))
		return
	}

	if http.MethodPost == request.Method {
		fmt.Println("Received POST")
		decoder := json.NewDecoder(request.Body)
		payload := twitch.TwitchNotificationPayload{}
		err := decoder.Decode(&payload)
		if err != nil {
			panic(err)
		}

		for _, notification := range payload.Notifications {
			displayName := twitchClient.FromUserId(notification.UserId)

			fmt.Println("Notification [UserId: " + notification.UserId + ", Name=" + displayName + ", Status: " + notification.Type + ", Title: " + notification.Title + "]")
			notificationMessage := strings.Replace(displayName, "_", "\\_", -1) + " is now live! http://twitch.tv/" + displayName
			discord.SendDiscordMessage(notificationMessage)
		}
	}
}

// Initialze
func Initialize() {
	settings.DumpEnvironmentVariables()
	twitchClient = twitch.NewTwitch(settings.GetClientId())
	InitializeEndPoints()

	go func() {
		err := http.ListenAndServe(":"+settings.GetHostPort(), nil)
		if nil != err {
			log.Fatalln(err)
		}
	}()
}

// Initializes HTTP End Points
func InitializeEndPoints() {
	http.HandleFunc("/"+NotifyEndPoint, OnTwitchNotification)
}

func main() {
	Initialize()

	hostUrl := settings.GetHostUrl() + "/" + NotifyEndPoint
	users := settings.GetUserNames()

	userIds, err := twitchClient.UserIdsFor(users)
	if nil != err {
		log.Fatalln(err)
	}

	twitchClient.SubscribeToStreams(hostUrl, userIds)

	select {}
}
