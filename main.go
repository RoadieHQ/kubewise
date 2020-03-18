package main

import (
	"log"
	"os"

	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/controller"
	"github.com/larderdev/kubewise/handlers/slack"
)

func main() {
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
	err := eventHandler.Init(&conf)

	if err != nil {
		log.Fatalln("Error initializing eventHandler", err)
	}

	controller.Start(&conf, eventHandler)
}
