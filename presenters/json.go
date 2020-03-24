package presenters

import (
	"encoding/json"
	"os"

	"github.com/larderdev/kubewise/kwrelease"
)

type EventJSON struct {
	AppName            string `json:"appName"`
	AppVersion         string `json:"appVersion"`
	Namespace          string `json:"namespace"`
	PreviousAppVersion string `json:"previousAppVersion"`
	Action             string `json:"action"`
	AppDescription     string `json:"appDescription"`
	InstallNotes       string `json:"installNotes"`
	MessagePrefix      string `json:"messagePrefix"`
}

func ReleaseEventToJSON(e *kwrelease.Event) ([]byte, error) {
	event := EventJSON{
		AppName:        e.GetAppName(),
		AppVersion:     e.GetAppVersion(),
		Namespace:      e.GetNamespace(),
		Action:         e.GetAction().String(),
		InstallNotes:   e.GetNotes(),
		AppDescription: e.GetDescription(),
	}

	previousAppVersion := e.GetPreviousAppVersion()

	// Prevents an empty string value in the JSON.
	if previousAppVersion != "" {
		event.PreviousAppVersion = previousAppVersion
	}

	if value, ok := os.LookupEnv("KW_MESSAGE_PREFIX"); ok && value != "" {
		event.MessagePrefix = value
	}

	return json.Marshal(event)
}
