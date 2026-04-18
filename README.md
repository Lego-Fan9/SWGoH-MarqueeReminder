# SWGoH Marquee Reminder

This is a marquee reminder for SWGoH.
It is in use in my [Discord server](https://discord.gg/cmZjsRBwTY)

# Self Hosting

This project is built to use [Comlink](https://github.com/swgoh-utils/swgoh-comlink) for the `/getEvents` endpoint it provides. It does not use the `/data` or `/localization` endpoints, and instead uses [the gamedata repo](https://github.com/swgoh-utils/gamedata) for that.

It optionally can use [SWGoH AE2](https://github.com/swgoh-utils/swgoh-ae2) for images in the embeds. My own [SWGoH AssetAPI](https://github.com/Lego-Fan9/swgoh-assetapi) may be used as well. The environment variable will be the same.

## Environment Variables

The **REQUIRED** environment variables are:
* DISCORD_WEBHOOK
* COMLINK_URL
  * See [Comlink](https://github.com/swgoh-utils/swgoh-comlink) for more details on how to set this up. It should include the protocol (usually http://)
* PING_ROLE
  * The Discord role ID that will be used to ping users. It should not include <@& or >, that will be inserted for you. Just the role ID

The **OPTIONAL** environment variables are:
* SWGOH_AE_URL
  * If not provided no image will be added.
* AVATAR_URL
  * If not provided the SWGoH Updates logo will be used
* DISCORD_USERNAME
  * If not provided "Marquee Reminder" will be the name

The **PROBABLY NOT NEEDED** environment variables are:
* ENV_PATH
  * If set this is what will be used to load the .env file. By default it will look besides the binary. If set to NONE it will not open any file.
* DOCKER
  * If set to "Y" it will disable env file loading and any future configuration or save files will be in a good place for docker (This env will be handled by any Docker image I distribute)
* CUSTOM_FORMAT
  * This will override the template used to format messages to discord when SWGOH_AE_URL is not supplied. See [Custom Messages](#custom-messages) for more details
* CUSTOM_FORMAT_IMG
 * This will override the template used to format messages to discord when there an image. See [Custom Messages](#custom-messages) for more details

## Running
You can get started by running this (assuming your env file is .env)
```
docker pull ghcr.io/lego-fan9/swgoh-marqueereminder:latest
docker run --name marquee-reminder -d --restart unless-stopped --env-file .env ghcr.io/lego-fan9/swgoh-marqueereminder
```

# Custom Messages
This tool uses Go's `text/template` library for formatting messages for Discord. This is done in `src/env/template.go`.
To modify messages you can change CUSTOM_FORMAT and CUSTOM_FORMAT_IMG.
The template must be a valid Discord webhook message. 
This is the struct containing the values you can use
```go
type MarqueeTemplateData struct {
	Role     string
	NameKey  string
	Filename string // Will only be set if there is an image
	Username string
	Avatar   string
}
```