package settings

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// The base url of the host for constructing callback urls
	HostUrlEnvVar string = "HOST_URL"

	// Heroku injects the exposed port via environment variable
	HostPortEnvVar string = "PORT"

	// Heroku injects the attached database host via environment variable
	DatabaseHostEnvVar string = "DATABASE_URL"

	// The Twitch App Client Identifier used when communicating with Twitch APIs
	ClientIdEnvVar string = "TWITCH_CLIENT_ID"

	// A comma delimited list of Twitch user names to subscribe to go live events for
	UsersEnvVar string = "TWITCH_USERS"

	// The discord web hook id environment variable
	DiscordWebHookIdEnvVar string = "DISCORD_WEBHOOK_ID"

	// The discord web hook token environment variable
	DiscordWebHookTokenEnvVar string = "DISCORD_WEBHOOK_TOKEN"

	// The default host url
	DefaultHostUrl string = "http://localhost"

	// Default Port Value when running locally
	DefaultPort string = "3001"
)

var (
	twitchClientId      string
	twitchUserNames     []string
	hostUrl             string
	hostPort            string
	discordWebHookId    string
	discordWebHookToken string
	databaseHost        string
)

// DumpEnvironmentVariables is a Debug Function to Dump All Environment Variables to stdout
func DumpEnvironmentVariables() {
	fmt.Println("--- ENV Vars ---")
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0] + " = " + os.Getenv(pair[0]))
	}
	fmt.Println("---------------")
}

// GetHostUrl Gets the Host URL base
func GetHostUrl() string {
	host := os.Getenv(HostUrlEnvVar)

	if host == "" {
		host = DefaultHostUrl
	}

	return host
}

// GetHostPort Gets the port the web app should be hosted on
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

// GetDatabaseHost Gets the database host url
func GetDatabaseHost() string {
	if "" != databaseHost {
		return databaseHost
	}

	databaseHost = os.Getenv(DatabaseHostEnvVar)

	if "" == databaseHost {
		log.Println("$DATABASE_URL not set.")
	}

	return databaseHost
}

// GetClientId Gets the Client Identifier Header for HTTP Requests
func GetClientId() string {
	if "" != twitchClientId {
		return twitchClientId
	}

	twitchClientId = os.Getenv(ClientIdEnvVar)
	return twitchClientId
}

// GetUserNames Gets the name of twitch users to listen for go live events
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

// GetDiscordHookId gets the Discord WebHook Id
func GetDiscordHookId() string {
	if "" != discordWebHookId {
		return discordWebHookId
	}

	discordWebHookId = os.Getenv(DiscordWebHookIdEnvVar)
	return discordWebHookId
}

// GetDiscordHookToken gets the Discord WebHook Token
func GetDiscordHookToken() string {
	if "" != discordWebHookToken {
		return discordWebHookToken
	}

	discordWebHookToken = os.Getenv(DiscordWebHookTokenEnvVar)
	return discordWebHookToken
}
