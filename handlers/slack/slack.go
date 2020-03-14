package slack

import (
	"fmt"
	"log"

	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/presenters"
	"github.com/slack-go/slack"
	"helm.sh/helm/v3/pkg/release"
)

var slackErrMsg = `
%s

You need to set both slack token and channel for slack notify,
using "--token/-t" and "--channel/-c", or using environment variables:

export KW_SLACK_TOKEN=slack_token
export KW_SLACK_CHANNEL=slack_channel

Command line flags will override environment variables

`

// Slack handler implements handler.Handler interface,
// Notify event to slack channel
type Slack struct {
	Token   string
	Channel string
}

func (s *Slack) Init(c *config.Config) error {
	s.Token = c.SlackToken
	s.Channel = c.SlackChannel

	return checkMissingSlackVars(s)
}

func (s *Slack) ObjectCreated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectCreatedMsg(currentRelease, previousRelease); msg != "" {
		notifySlack(s, msg, "created")
	}
}

func (s *Slack) ObjectDeleted(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectDeletedMsg(currentRelease, previousRelease); msg != "" {
		notifySlack(s, msg, "created")
	}
}

func (s *Slack) ObjectUpdated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectUpgradedMsg(currentRelease, previousRelease); msg != "" {
		notifySlack(s, msg, "created")
	}
}

func notifySlack(s *Slack, msg, action string) {
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

func checkMissingSlackVars(s *Slack) error {
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf(slackErrMsg, "Missing slack token or channel")
	}

	return nil
}
