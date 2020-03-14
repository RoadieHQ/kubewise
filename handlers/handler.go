package handlers

import (
	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/handlers/slack"
	"helm.sh/helm/v3/pkg/release"
)

type Handler interface {
	Init(c *config.Config) error
	ObjectCreated(currentRelease, previousRelease *release.Release)
	ObjectDeleted(currentRelease, previousRelease *release.Release)
	ObjectUpdated(currentRelease, previousRelease *release.Release)
}

var Map = map[string]interface{}{
	"default": &Default{},
	"slack":   &slack.Slack{},
}

type Default struct {
}

func (d *Default) Init(c *config.Config) error {
	return nil
}

func (d *Default) ObjectCreated(currentRelease, previousRelease *release.Release) {

}

func (d *Default) ObjectDeleted(currentRelease, previousRelease *release.Release) {

}

func (d *Default) ObjectUpdated(currentRelease, previousRelease *release.Release) {

}
