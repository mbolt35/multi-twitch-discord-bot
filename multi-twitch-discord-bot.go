package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mbolt35/multi-twitch-discord-bot/discord"
	"github.com/mbolt35/multi-twitch-discord-bot/settings"
	"github.com/mbolt35/multi-twitch-discord-bot/twitch"

	httputil "github.com/mbolt35/multi-twitch-discord-bot/util"
)

// NotifyEndPoint The end point we'll bind to for receiving http requests
const NotifyEndPoint string = "notify"

var (
	twitchClient  twitch.TwitchClient
	discordClient discord.DiscordClient
	done          chan bool
)

// escapeUnderscore escapes any underscore characters in the string
func escapeUnderscore(s string) string {
	return strings.Replace(s, "_", "\\_", -1)
}

// logNotification outputs the twitch notification to stdout
func logNotification(notification twitch.TwitchNotification) {
	log.Printf(
		"TwitchNotification[\n  UserId: %s,\n  DisplayName: %s\n  Type: %s\n  Title: %s\n  GameId: %s\n  StartedAt: %s\n  ViewerCount: %d\n]\n",
		notification.UserId,
		twitchClient.FromUserId(notification.UserId),
		notification.Type,
		notification.Title,
		notification.GameId,
		notification.StartedAt,
		notification.ViewerCount)
}

// newTwitchLiveMessage returns the message to send to the discord channel for a user going live.
func newTwitchLiveMessage(userId string) string {
	userName := twitchClient.FromUserId(userId)
	return escapeUnderscore(userName) + " is now live! " + twitch.UserStreamUrl(userName)
}

// OnTwitchNotification Handles Incoming Twitch Notifications
func OnTwitchNotification(rw http.ResponseWriter, request *http.Request) {
	// The GET occurs after the subscription to the stream update is made
	// The main purpose is to provide twitch a way to validate the endpoint
	if http.MethodGet == request.Method {
		fmt.Println("Received GET")

		q := request.URL.Query()
		log.Println(q)

		mode := q.Get(twitch.TwitchHubModeQueryParameter)
		//topic := q.Get(twitch.TwitchHubTopicQueryParameter)

		if twitch.TwitchModeDenied == mode {
			reason := q.Get(twitch.TwitchHubReasonQueryParameter)
			log.Println("Failed to Subscribe to Webhook: " + reason)
			rw.WriteHeader(http.StatusOK)
			return
		}

		challenge := q.Get(twitch.TwitchHubChallengeQueryParameter)
		//lease := q.Get(twitch.TwitchHubLeaseQueryParameter)

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(challenge))
		return
	}

	// The POST occurs when the actual event of going live occurs, we'll need to decode
	// the payload into notification objects, then send a discord webhook message
	if http.MethodPost == request.Method {
		fmt.Println("Received POST")

		var payload twitch.TwitchNotificationPayload
		err := httputil.DecodeJson(request.Body, &payload)
		if err != nil {
			panic(err)
		}

		// Iterate through all notifications, send discord message for live streams
		for _, notification := range payload.Notifications {
			logNotification(notification)
			discordClient.SendDiscordMessage(newTwitchLiveMessage(notification.UserId))
		}
	}
}

// Initialze
func Initialize() {
	settings.DumpEnvironmentVariables()

	// Create twitch and discord clients
	twitchClient = twitch.NewTwitch(settings.GetClientId())
	discordClient = discord.NewDiscord(settings.GetDiscordHookId(), settings.GetDiscordHookToken())

	InitializeEndPoints()
	go StartWebServer(settings.GetHostPort())
}

// InitializeEndPoints Initializes HTTP End Points
func InitializeEndPoints() {
	http.HandleFunc("/"+NotifyEndPoint, OnTwitchNotification)
}

// StartWebServer starts running the web server for receiving requests from twitch
func StartWebServer(port string) {
	err := http.ListenAndServe(":"+port, nil)
	if nil != err {
		log.Fatalln(err)
	}
	done <- true
}

// main Entry Point
func main() {
	Initialize()

	// Acquire Host URL and Twitch Users to Watch for Live Events
	hostUrl := settings.GetHostUrl() + "/" + NotifyEndPoint
	users := settings.GetUserNames()

	// Convert Twitch User Names to User Ids
	userIds, err := twitchClient.UserIdsFor(users)
	if nil != err {
		log.Fatalln(err)
	}

	// Subscribe to Stream Live Events
	twitchClient.SubscribeToStreams(hostUrl, userIds)

	// Wait for the web server to complete
	<-done
}
