package presenters

import (
	"os"

	"github.com/larderdev/kubewise/kwrelease"
	rspb "helm.sh/helm/v3/pkg/release"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ReleaseEventForJSON is a container for a subset of the properties available in a typical
// Release Event. Release events contain all of the files present in the Helm package which
// was installed. They can be massive, larger than your typical JSON payload. They can also
// contain sensitive data. For this reason, we use ReleaseEventForJSON to whitelist properties.
//
// Some fields may be denormalized because it makes more sense for users of the webhooks.
// The user wants to know what time the event occurred as a first class concept in the Json.
type ReleaseEventForJSON struct {
	AppName              string       `json:"appName"`
	AppVersion           string       `json:"appVersion"`
	Namespace            string       `json:"namespace"`
	PreviousAppVersion   string       `json:"previousAppVersion,omitempty"`
	Action               string       `json:"action"`
	AppDescription       string       `json:"appDescription"`
	InstallNotes         string       `json:"installNotes"`
	MessagePrefix        string       `json:"messagePrefix,omitempty"`
	CreatedAt            meta_v1.Time `json:"createdAt"`
	UpdatedAt            meta_v1.Time `json:"updatedAt"`
	SecretUID            types.UID    `json:"secretUid"`
	ChartVersion         string       `json:"chartVersion"`
	PreviousChartVersion string       `json:"previousChartVersion"`
	ReleaseDescription   string       `json:"releaseDescription"`
}

// ToReleaseEventForJSON takes a release Event and turns it into a ReleaseEventForJSON. It holds
// knowledge such as when to omit certain empty time fields and where to find the message prefix.
func ToReleaseEventForJSON(e *kwrelease.Event) *ReleaseEventForJSON {
	event := ReleaseEventForJSON{
		AppName:              e.GetAppName(),
		AppVersion:           e.GetAppVersion(),
		Namespace:            e.GetNamespace(),
		Action:               e.GetAction().String(),
		InstallNotes:         e.GetNotes(),
		AppDescription:       e.GetAppDescription(),
		CreatedAt:            e.GetSecretCreationTimestamp(),
		SecretUID:            e.GetSecretUID(),
		ChartVersion:         e.GetChartVersion(),
		PreviousChartVersion: e.GetPreviousChartVersion(),
		ReleaseDescription:   e.GetReleaseDescription(),
		PreviousAppVersion:   e.GetPreviousAppVersion(),
	}

	if value := e.GetLabelsModifiedAtTimestamp(); !value.IsZero() {
		event.UpdatedAt = value
	}

	if value, ok := os.LookupEnv("KW_MESSAGE_PREFIX"); ok {
		event.MessagePrefix = value
	}

	return &event
}

// ExistingReleasesForJSON is used to marshal Helm release objects so they can be sent to an API.
//
// There are two problems which just directly marshaling Helm release objects.
//   1. They may contain sensitive data which should not leave the cluster.
//   2. They are huge when marshalled because all the templates are stored within.
//
// By implementing a custom struct we effectively whitelist the properties which should be
// send to any API.
type ExistingReleasesForJSON struct {
	MessagePrefix string `json:"messagePrefix,omitempty"`
	// Do not use omitempty on existingReleases. Doing so requires the API to have a null check
	// before mapping over the existingReleases and generally makes it more likely that bugs will
	// occur for the users.
	ExistingReleases []*ExistingReleaseForJSON `json:"existingReleases"`
}

// ToExistingReleasesForJSON takes a release Event and turns it into a ToExistingReleasesForJSON.
// It holds knowledge such where to find the message prefix environment variable.
func ToExistingReleasesForJSON(releases []*rspb.Release) *ExistingReleasesForJSON {
	container := ExistingReleasesForJSON{}

	if value, ok := os.LookupEnv("KW_MESSAGE_PREFIX"); ok {
		container.MessagePrefix = value
	}

	// Using make here ensures that the empty state is an empty slice rather than null. It's the
	// difference between receiving {"existingReleases":[]} at the API vs. {"existingReleases":null}
	existingReleases := make([]*ExistingReleaseForJSON, 0)
	for _, release := range releases {
		existingReleases = append(existingReleases, toExistingReleaseForJSON(release))
	}
	container.ExistingReleases = existingReleases

	return &container
}

// ExistingReleaseForJSON is a single release Event which may be contained in a
// ExistingReleasesForJSON object.
type ExistingReleaseForJSON struct {
	AppName            string `json:"appName"`
	AppVersion         string `json:"appVersion"`
	Namespace          string `json:"namespace"`
	AppDescription     string `json:"appDescription"`
	InstallNotes       string `json:"installNotes"`
	ChartVersion       string `json:"chartVersion"`
	ReleaseDescription string `json:"releaseDescription"`
}

func toExistingReleaseForJSON(r *rspb.Release) *ExistingReleaseForJSON {
	return &ExistingReleaseForJSON{
		AppName:            r.Name,
		AppVersion:         r.Chart.AppVersion(),
		Namespace:          r.Namespace,
		InstallNotes:       r.Info.Notes,
		AppDescription:     r.Chart.Metadata.Description,
		ChartVersion:       r.Chart.Metadata.Version,
		ReleaseDescription: r.Info.Description,
	}
}
