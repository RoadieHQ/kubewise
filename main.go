package main

import (
	"log"
	"os"

	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/controller"
	"github.com/larderdev/kubewise/handlers"
	"github.com/larderdev/kubewise/handlers/logs"
	"github.com/larderdev/kubewise/handlers/slack"
)

func createLogFile() *os.File {
	file, err := os.OpenFile("log/events.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Unable to create file", err)
	}
	return file
}

func createFileLoggingEventHandler(file *os.File) handlers.Handler {
	logger := log.New(file, "", 0644)
	eventHandler := &logs.Logs{Logger: logger}
	return eventHandler
}

func main() {
	// file := createLogFile()
	// defer file.Close()
	// eventHandler := createFileLoggingEventHandler(file)

	eventHandler := new(slack.Slack)

	channel := "#general"
	if value, ok := os.LookupEnv("KW_SLACK_CHANNEL"); ok {
		channel = value
	}

	var token string
	if value, ok := os.LookupEnv("KW_SLACK_TOKEN"); ok {
		token = value
	} else {
		log.Fatalln("Missing environment variable KW_SLACK_TOKEN")
	}

	namespace := ""
	if value, ok := os.LookupEnv("KW_NAMESPACE"); ok {
		namespace = value
	}

	conf := config.Config{
		Namespace:    namespace,
		SlackChannel: channel,
		SlackToken:   token,
	}
	eventHandler.Init(&conf)

	controller.Start(&conf, eventHandler)
}
