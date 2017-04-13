# isitup
Tool written in Go, which checks the Availability of defined Services or Resources - sends a Message via a Telegram Bot, if a service is offline

Sends only one Message in 10 Minutes if a Service is down at multiple scans. When its up again, you'll receive a Message.

# Usage
```
Usage of ./isitup:
  -config string
    	Sets the location of the config file (default "./settings.toml")
  -interval string
    	Sets the Scan interval in seconds - WARNING: Overrides the Setting in the config file. (default "60")
  -logfile string
    	Sets the location of the logfile. (default "./isitup.log")
  -servicefile string
    	Sets the location of the service file - WARNING: Overrides the Setting in the config file. (default "./service.isitup")


```

# Setup

- cd /opt
- git clone https://github.com/alxndr13/isitup
- chmod +x install.sh && ./install.sh

### settings.toml:
- Enter the Access Token of your Telegram Bot
- Enter the Chat ID of the person, who should get notifications (You have to figure out, how to get this yourself. :p - You have to Message the Bot first [/start], before you can receive Messages.)

### service.isitup
- Enter IP Adress as well as the Port, line after line
- Service Names can be added after a single dash (-)

Example:
```
192.168.0.100:80-Local Webserver
8.8.8.8:53-Google DNS
```

## Startup

- ``` systemctl enable isitup && systemctl start isitup ```

### ToDo:
- multiple receivers
- udp?
- ..?

