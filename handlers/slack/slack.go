package slack

import (
	"log"

	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/presenters"
	"github.com/slack-go/slack"
	"helm.sh/helm/v3/pkg/release"
)

type Slack struct {
	Token   string
	Channel string
}

func (s *Slack) Init(c *config.Config) error {
	s.Token = c.SlackToken
	s.Channel = c.SlackChannel
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
