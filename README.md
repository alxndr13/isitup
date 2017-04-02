# isitup
Tool written in Go, which checks the Availability of defined Services or Resources - sends a Message via a Telegram Bot, if a service is offline

# Setup

settings.toml:
- Enter the Access Token of your Telegram Bot
- Enter the Chat ID of the person, who should get notifications (You have to figure out, how to get this yourself. :p)

hosts.isitup
- Enter IP Adress as well as the Port, line after line

## Example:
`192.168.0.100:80

8.8.8.8:53`
