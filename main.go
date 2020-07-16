package main

import (
	"log"
	"os"

	"github.com/RoadieHQ/kubewise/controller"
	"github.com/RoadieHQ/kubewise/handlers"
	"github.com/RoadieHQ/kubewise/handlers/googlechat"
	"github.com/RoadieHQ/kubewise/handlers/slack"
	"github.com/RoadieHQ/kubewise/handlers/webhook"
	"github.com/RoadieHQ/kubewise/kwrelease"
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
	eventHandler.HandleServerStartup(kwrelease.ListActiveReleases())
	// This is a blocking call. Code placed after this won't run until teardown.
	controller.Start(eventHandler)
}
