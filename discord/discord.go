package discord

import (
	"bytes"
	"github.com/mbolt35/multi-twitch-discord-bot/settings"
	"github.com/mbolt35/multi-twitch-discord-bot/util"
	"log"
	"net/http"
	"strings"
)

const (
	// Discord WebHook Base Url
	DiscordWebHookUrl string = "https://discordapp.com/api/webhooks"
)

// Discord WebHook Request Payload
type DiscordWebHookMessage struct {
	Message string `json:"content"`
}

// Gets the discord webbhook base url
func GetDiscordWebHookUrl() string {
	return strings.Join([]string{DiscordWebHookUrl, settings.GetDiscordHookId(), settings.GetDiscordHookToken()}, "/")
}

// Sends a message to the discord server/channel using the webhook
func SendDiscordMessage(message string) {
	discordMessage := DiscordWebHookMessage{
		Message: message,
	}

	jsonBytes, err := httputil.EncodeJson(discordMessage)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s\n", jsonBytes)

	resp, err := http.Post(GetDiscordWebHookUrl(), httputil.JsonContentType, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
