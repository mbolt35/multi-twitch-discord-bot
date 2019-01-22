multi-twitch-discord-bot
======

A light-weight web application which can subscribe to "live" events from multiple Twitch users, the redestribute the event through a Discord webhook. This "bot" can be hosted on a free heroku web dyno + free heroku Postgres database add-on, and it does _not_ require connecting a single Twitch account or Discord account. 

----

## Why?
Members of a Discord server wanted to create a **#live-now** channel to support specific Twitch streams when they go live. We determined very quickly that there are ways to do this _for each user_. However, using available bots, there would've been a bot for each user (permission management hell), and the process for creation involved each user linking both their Twitch and Discord accounts. We all agreed that having a single Discord bot listen for events for a configured set of Twitch users would be optimal. 

----

## Simplicity is Bliss
The process:
* **Creation of Discord Bot**: Fortunately, Discord provides a very easy way to send messages to a specific channel using [Execute WebHook](https://discordapp.com/developers/docs/resources/webhook#execute-webhook). There is also a [Discord Guide](https://support.discordapp.com/hc/en-us/articles/228383668) on creating webooks via Server settings. 
* **Subscribe to Twitch Events**: Application should consume a list of Twitch user names, subscribe to stream updates for each, and use the Discord webhook API to send messages relating to Twitch events.
* **Determine Where/How to Host**: Because the "bot" needed hosting in order to receive notifications of Twitch events, we used [Heroku](https://www.heroku.com/) to host and deploy (using this GitHub repository) our application. 

----

## Guide: How to Setup

...
