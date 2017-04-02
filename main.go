package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var receiver = ""
var token = ""
var interval = ""
var hosts = "hosts.isitup"

func main() {
	enablelogging()
	read_config()

	// Check if there is a Hostsfile
	if hostsfile_available() {
		Info.Println("Starting Scanning...")
		for {
			go checkupness()
			seconds, _ := strconv.Atoi(interval)
			time.Sleep(time.Second * time.Duration(seconds))
		}
	}
}

func read_config() {
	viper.SetConfigFile("settings.toml")
	err := viper.ReadInConfig()
	if err != nil {
		Error.Println("Could not open config file.. Are you sure it is there?")
	} else {
		Info.Println("Reading Config..")
		token = viper.GetString("TOKEN")
		Info.Println("TOKEN is: " + token)
		receiver = viper.GetString("RECEIVER")
		Info.Println("Receiver ID is: " + receiver)
		interval = viper.GetString("INTERVAL")
		Info.Println("Scan Interval is: " + interval)

	}

}

func enablelogging() {
	// Initializing Logging
	file, err := os.OpenFile("isitup.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Could not open Log File..")
	} else {
		InitLogging(ioutil.Discard, file, file, file)
		Info.Println("Logging initialized.")
	}
}

func hostsfile_available() bool {
	f, err := os.Open(hosts)
	defer f.Close()
	if err != nil {
		Error.Println("Could not open hostsfile. Exiting.")
		return false
	} else {
		Info.Println("Opened Hostsfile.")
		return true
	}

}

func send_message(mes string) {
	// i need a another way to store the Token
	url := "https://api.telegram.org/bot" + token + "/"
	// Send message
	res, err := http.Get(url + "sendMessage?text=" + mes + "&chat_id=" + receiver)
	if err != nil {
		Error.Println("couldnt send Message to Telegram.")
	} else {
		if res.StatusCode != 200 {
			Warning.Println("Could not send message. Status Code: " + res.Status)
		} else {
			Info.Println("Send Message to: " + receiver)
		}

	}

}

func checkupness() {
	//Check Host and Port via TCP
	file, err1 := os.Open(hosts)
	if err1 != nil {
		Error.Println("could not open Hostsfile")
	}
	defer file.Close()

	//Read the File and check the Hosts inside
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err2 := net.Dial("tcp", scanner.Text())
		Info.Println("Checking " + scanner.Text() + "..")
		// The usual error handling.
		if err2 != nil {
			Warning.Println("#####################################")
			Warning.Println("Host " + scanner.Text() + " is down.")
			Warning.Println("#####################################")
			send_message("Host " + scanner.Text() + " is down.")
		} else {
			Info.Println("Host " + scanner.Text() + " is up.")
		}
	}
}
