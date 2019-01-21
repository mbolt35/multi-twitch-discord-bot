package discord

import (
	"bytes"
	"net/http"
	"strings"

	httputil "github.com/mbolt35/multi-twitch-discord-bot/util/http"
)

// Discord WebHook Base Url
const DiscordWebHookUrl string = "https://discordapp.com/api/webhooks"

// Discord WebHook Request Payload
type DiscordWebHookMessage struct {
	Message string `json:"content"`
}

// DiscordClient implementation
type discord struct {
	webHookId    string
	webHookToken string
}

// Interface representing a Discord client
type DiscordClient interface {
	SendDiscordMessage(message string) error
}

// Discord Client Factory
func NewDiscord(webHookId string, webHookToken string) DiscordClient {
	instance := discord{
		webHookId:    webHookId,
		webHookToken: webHookToken,
	}

	return &instance
}

// Gets the discord webbhook base url
func getDiscordWebHookUrl(hookId string, hookToken string) string {
	return strings.Join([]string{DiscordWebHookUrl, hookId, hookToken}, "/")
}

// Sends a message to the discord server/channel using the webhook
func (d *discord) SendDiscordMessage(message string) error {
	discordMessage := DiscordWebHookMessage{
		Message: message,
	}

	jsonBytes, err := httputil.EncodeJson(discordMessage)
	if err != nil {
		return err
	}

	webHookUrl := getDiscordWebHookUrl(d.webHookId, d.webHookToken)
	resp, err := http.Post(webHookUrl, httputil.JsonContentType, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
