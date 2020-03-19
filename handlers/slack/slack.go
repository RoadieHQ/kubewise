package slack

import (
	"log"
	"os"

	"github.com/larderdev/kubewise/presenters"
	"github.com/slack-go/slack"
	"helm.sh/helm/v3/pkg/release"
)

type Slack struct {
	Token   string
	Channel string
}

func (s *Slack) Init() error {
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

	s.Token = token
	s.Channel = channel
	return nil
}

func (s *Slack) ObjectCreated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectCreatedMsg(currentRelease, previousRelease); msg != "" {
		sendMessage(s, msg)
	}
}

func (s *Slack) ObjectDeleted(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectDeletedMsg(currentRelease, previousRelease); msg != "" {
		sendMessage(s, msg)
	}
}

func (s *Slack) ObjectUpdated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectUpgradedMsg(currentRelease, previousRelease); msg != "" {
		sendMessage(s, msg)
	}
}

func sendMessage(s *Slack, msg string) {
	api := slack.New(s.Token)
	text := slack.MsgOptionText(msg, false)
	asUser := slack.MsgOptionAsUser(true)

	channelID, timestamp, err := api.PostMessage(s.Channel, text, asUser)

	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
