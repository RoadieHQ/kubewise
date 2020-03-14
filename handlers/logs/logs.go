package logs

import (
	"log"

	"github.com/larderdev/kubewise/config"
	"github.com/larderdev/kubewise/presenters"
	"helm.sh/helm/v3/pkg/release"
)

type Logs struct {
	Logger *log.Logger
}

func (lh *Logs) Init(c *config.Config) error {
	lh.Logger.Println("Logger initialized")
	log.Println("Logger initialized")
	return nil
}

func (lh *Logs) ObjectCreated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectCreatedMsg(currentRelease, previousRelease); msg != "" {
		lh.Logger.Printf(msg)
	}
}

func (lh *Logs) ObjectDeleted(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectDeletedMsg(currentRelease, previousRelease); msg != "" {
		lh.Logger.Printf(msg)
	}
}

func (lh *Logs) ObjectUpdated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectCreatedMsg(currentRelease, previousRelease); msg != "" {
		lh.Logger.Printf(msg)
	}
}
