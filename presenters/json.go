package presenters

import (
	"encoding/json"
	"os"

	"github.com/larderdev/kubewise/kwrelease"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Some fields are being denormalized here (such as UpdatedAtTimestamp being taken out of
// Secret MetaData) because it makes more sense for users of the webhooks. The user wants to
// know what time the event occurred as a first class concept in the Json.
type EventJSON struct {
	AppName              string       `json:"appName"`
	AppVersion           string       `json:"appVersion"`
	Namespace            string       `json:"namespace"`
	PreviousAppVersion   string       `json:"previousAppVersion"`
	Action               string       `json:"action"`
	AppDescription       string       `json:"appDescription"`
	InstallNotes         string       `json:"installNotes"`
	MessagePrefix        string       `json:"messagePrefix"`
	CreatedAt            meta_v1.Time `json:"createdAt"`
	UpdatedAt            meta_v1.Time `json:"updatedAt"`
	SecretUID            types.UID    `json:"secretUid"`
	ChartVersion         string       `json:"chartVersion"`
	PreviousChartVersion string       `json:"previousChartVersion"`
}

func ReleaseEventToJSON(e *kwrelease.Event) ([]byte, error) {
	event := EventJSON{
		AppName:              e.GetAppName(),
		AppVersion:           e.GetAppVersion(),
		Namespace:            e.GetNamespace(),
		Action:               e.GetAction().String(),
		InstallNotes:         e.GetNotes(),
		AppDescription:       e.GetDescription(),
		CreatedAt:            e.GetSecretCreationTimestamp(),
		SecretUID:            e.GetSecretUID(),
		ChartVersion:         e.GetChartVersion(),
		PreviousChartVersion: e.GetPreviousChartVersion(),
	}

	if value := e.GetLabelsModifiedAtTimestamp(); !value.IsZero() {
		event.UpdatedAt = value
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
