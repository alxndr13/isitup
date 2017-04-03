package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Globals - is there a better way for this?
var receiver string
var token string
var interval string
var service string
var logPath string
var confPath string

func main() {
	// Init
	initFlags()
	enablelogging()
	readConfig()

	// Check if there is a servicefile and let it rip
	if servicefileAvailable() {
		for {
			Info.Println("Starting Scanning...")
			go checkupness()
			seconds, _ := strconv.Atoi(interval)
			time.Sleep(time.Second * time.Duration(seconds))
		}
	}
}

func initFlags() {
	// All that command line arguments
	logPtr := flag.String("logfile", "./isitup.log", "Sets the location of the logfile.")
	svcPtr := flag.String("servicefile", "./service.isitup", "Sets the location of the service file - WARNING: Overrides the Setting in the config file.")
	confPtr := flag.String("config", "./settings.toml", "Sets the location of the config file")
	intPtr := flag.String("interval", "60", "Sets the Scan interval in seconds - WARNING: Overrides the Setting in the config file.")
	flag.Parse()
	// to the globals
	logPath = *logPtr
	service = *svcPtr
	confPath = *confPtr
	interval = *intPtr

}

func readConfig() {
	viper.SetConfigFile(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		Error.Println("Could not open config file.. Are you sure it is there?")
	} else {
		Info.Println("Reading Config..")
		token = viper.GetString("TOKEN")
		Info.Println("TOKEN is: " + token)
		receiver = viper.GetString("RECEIVER")
		Info.Println("Receiver ID is: " + receiver)
		// Only get the interval value, if the command line argument isnt set.
		if len(interval) == 0 {
			interval = viper.GetString("INTERVAL")
			Info.Println("Scan Interval is: " + interval)
		}
		if len(service) == 0 {
			service = viper.GetString("SERVICEFILE")
			Info.Println("Servicefile path is: " + service)
		}

	}

}

func enablelogging() {
	// Initializing Logging
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Warning.Println("Could not open logfile..")
	} else {
		InitLogging(ioutil.Discard, file, file, file)
		// Startup Message
		Info.Println("#####################################")
		Info.Println("Starting Isitup")
		Info.Println("#####################################")
		// Logging Initialization finished
		Info.Println("Logging initialized.")
	}
}

func servicefileAvailable() bool {
	file, err := os.Open(service)
	defer file.Close()
	if err != nil {
		Error.Println("Could not open service file. Exiting.")
		return false
	}
	Info.Println("Opened service file.")
	return true

}

func sendMessage(mes string) {
	url := "https://api.telegram.org/bot" + token + "/"
	// Send message with a get request
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
	//Check Service via TCP
	file, err1 := os.Open(service)
	if err1 != nil {
		Error.Println("could not open Servicefile")
	}
	defer file.Close()

	//Read the File and check the Services inside
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue // if it is a comment, then check the next line
		}
		input := strings.Split(scanner.Text(), "-")
		// This is where the magic happens
		_, err2 := net.Dial("tcp", input[0])
		Info.Println("Checking " + input[1] + "..")
		// The usual error handling.
		if err2 != nil {
			Warning.Println("#####################################")
			Warning.Println("Service " + input[1] + " on " + input[0] + " is down.")
			Warning.Println("#####################################")
			sendMessage("Service " + input[1] + " on " + input[0] + " is down.")
		} else {
			Info.Println("Service " + input[1] + " on " + input[0] + " is up.")
		}
	}
}
