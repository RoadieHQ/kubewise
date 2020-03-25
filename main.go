package main

import (
	"log"
	"os"

	"github.com/larderdev/kubewise/controller"
	"github.com/larderdev/kubewise/handlers"
	"github.com/larderdev/kubewise/handlers/googlechat"
	"github.com/larderdev/kubewise/handlers/slack"
	"github.com/larderdev/kubewise/handlers/webhook"
)

func main() {
	if _, ok := os.LookupEnv("KW_HANDLER"); !ok {
		log.Fatalln("KW_HANDLER environment variable is required.")
	}

	var eventHandler handlers.Handler
	switch os.Getenv("KW_HANDLER") {
	case "googlechat":
		eventHandler = new(googlechat.GoogleChat)
	case "webhook":
		eventHandler = new(webhook.Webhook)
	// Slack is the default for backwards compatibility reasons. It was the first handler.
	default:
		eventHandler = new(slack.Slack)
	}

	eventHandler.Init()
	controller.Start(eventHandler)
}
