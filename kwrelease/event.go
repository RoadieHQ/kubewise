package kwrelease

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"helm.sh/helm/v3/pkg/chartutil"
	rspb "helm.sh/helm/v3/pkg/release"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kbtypes "k8s.io/apimachinery/pkg/types"
)

// This file contains most of the simple Getters and Setters for the Event struct. There are
// other, more complex methods elsewhere.

// Event marks a release event. Upgrading a Helm chart from one version to another
// is an example of an event. Installing a chart is a different type of event.
type Event struct {
	// Describes the action that happened to the secret in order to trigger this event. Events
	// occur when secrets are created, updated or deleted. What was the action that led to this
	// particular event.
	SecretAction         string
	CurrentReleaseSecret *api_v1.Secret
	currentRelease       *rspb.Release
	previousRelease      *rspb.Release
}

// Init pre-loads data for the event.
// - Brand new installs will only have e.currentRelease.
// - Upgrades, rollbacks and uninstalls will have e.currentRelease and e.previousRelease (unless
//   they have been deleted by the user or something)
func (e *Event) Init() error {
	// Fetching the release from the secret store is unnecessary except for the fact that we need
	// to decode it and the safest way to do that is to let the Helm lib do it. THe Helm code has
	// a function called driver.decodeRelease but it is private and cannot be accessed directly.
	// By Getting the release from the store we get it back decoded for us.
	e.currentRelease = e.GetRelease(e.CurrentReleaseSecret.Name)
	e.previousRelease = e.getPreviousRelease()

	return nil
}

// GetAppName returns the name of the application being installed by the Helm chart.
func (e *Event) GetAppName() string {
	return e.currentRelease.Name
}

// GetAppVersion returns the version of the application being installed by the Helm chart.
func (e *Event) GetAppVersion() string {
	return e.currentRelease.Chart.AppVersion()
}

// GetPreviousAppVersion returns the version of the application being superseeded by the
// current upgrade event. It can be used to send notifications like "upgrading X from version 1
// to version 2"
func (e *Event) GetPreviousAppVersion() string {
	if e.previousRelease != nil {
		return e.previousRelease.Chart.AppVersion()
	}
	return ""
}

// GetNamespace returns the namespace that the event is occurring in.
func (e *Event) GetNamespace() string {
	return e.currentRelease.Namespace
}

// GetAppDescription returns the description of the application being installed or upgraded.
func (e *Event) GetAppDescription() string {
	return e.currentRelease.Chart.Metadata.Description
}

// GetReleaseDescription returns the Description of the release. This is primarily useful
// because it contains information about the cause of failure when a failure occurs.
func (e *Event) GetReleaseDescription() string {
	return e.currentRelease.Info.Description
}

// GetNotes returns the install notes of the Helm package being installed or upgraded.
func (e *Event) GetNotes() string {
	return e.currentRelease.Info.Notes
}

// GetSecretUID returns the UID of the release secret. Note that it returns a custom struct
// which is defined in the Kubernetes library. It's not a string or suchlike.
func (e *Event) GetSecretUID() kbtypes.UID {
	return e.CurrentReleaseSecret.GetUID()
}

// GetSecretCreationTimestamp returns the time that the release secret was created by Helm.
// Note that it returns a custom Time object which is defined in the Kubernetes meta_v1 API. It's
// not an instance of the standard go Time struct.
func (e *Event) GetSecretCreationTimestamp() meta_v1.Time {
	return e.CurrentReleaseSecret.GetCreationTimestamp()
}

// GetLabelsModifiedAtTimestamp returns the modifiedAt time in the release Meta. It's stored in the
// release as an Int and thus must be converted to a meta_v1.Time so it can be more easily
// compared with the creation timestamp.
func (e *Event) GetLabelsModifiedAtTimestamp() meta_v1.Time {
	labels := e.CurrentReleaseSecret.GetObjectMeta().GetLabels()

	// This has happened in regular use.
	if labels["modifiedAt"] == "" {
		return meta_v1.Time{}
	}

	i, err := strconv.ParseInt(labels["modifiedAt"], 10, 64)
	if err != nil {
		log.Println("Failed to ParseInt secret.GetObjectMeta().GetLabels()['modifiedAt']:", labels["modifiedAt"])
		return meta_v1.Time{}
	}

	return meta_v1.Unix(i, 0)
}

// GetChartVersion returns the version of the Helm chart being installed. This is different than
// the version of the application being installed. The same application version may span multiple
// chart versions.
func (e *Event) GetChartVersion() string {
	return e.currentRelease.Chart.Metadata.Version
}

// GetPreviousChartVersion returns the version of the previously installed Helm chart. This is
// useful during an upgrade when we wish to see how significant the upgrade is.
func (e *Event) GetPreviousChartVersion() string {
	if e.previousRelease != nil {
		return e.previousRelease.Chart.Metadata.Version
	}
	return ""
}

// IsAppVersionChanged makes it easy to tell if the application is upgraded when upgrading from
// one Helm Chart version to another.
func (e *Event) IsAppVersionChanged() bool {
	return e.GetAppVersion() != e.GetPreviousAppVersion()
}

// GetConfigDiffYAML is useful to show what has changed during a chart upgrade or rollback. It
// will show a diff of the values file in the Slack message or other notification. Be careful
// with secrets.
func (e *Event) GetConfigDiffYAML() string {
	currentReleaseConfigYAML, err := chartutil.Values(e.currentRelease.Config).YAML()

	if err != nil {
		// Do NOT log the values. Could leak sensitive data.
		log.Println("Error rendering current release config to YAML for application:", e.GetAppName())
		return ""
	}

	if e.previousRelease == nil {
		return fmt.Sprintf("%s\n", currentReleaseConfigYAML)
	}

	previousReleaseConfigYAML, err := chartutil.Values(e.previousRelease.Config).YAML()

	if err != nil {
		// Do NOT log the values. Could leak sensitive data.
		log.Println("Error rendering previous release config to YAML for application:", e.GetAppName())
		return ""
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(previousReleaseConfigYAML),
		B:        difflib.SplitLines(currentReleaseConfigYAML),
		FromFile: "Old Values",
		ToFile:   "New Values",
		Context:  3,
	}

	diffText, err := difflib.GetUnifiedDiffString(diff)

	if err != nil {
		log.Println("Error diffing chart for application:", e.GetAppName())
		return ""
	}

	return diffText
}

// GetAction returns the action which is being performed in this Event. It may be an install,
// upgrade or other Event.
func (e *Event) GetAction() Action {
	if e.currentRelease.Info.Status == rspb.StatusPendingInstall {
		return ActionPreInstall
	} else if e.currentRelease.Info.Status == rspb.StatusPendingUpgrade {
		return ActionPreUpgrade
	} else if e.currentRelease.Info.Status == rspb.StatusPendingRollback {
		return ActionPreRollback
	} else if e.currentRelease.Info.Status == rspb.StatusDeployed {
		if e.previousRelease == nil {
			return ActionPostInstall
		} else if strings.HasPrefix(e.currentRelease.Info.Description, "Rollback") {
			return ActionPostRollback
		} else if strings.HasPrefix(e.currentRelease.Info.Description, "Upgrade") {
			return ActionPostUpgrade
		}
		return ActionPostReplace
	} else if e.currentRelease.Info.Status == rspb.StatusFailed {
		if e.previousRelease == nil {
			return ActionFailedInstall
		}
		// There is no way to differentiate between an upgrade and a rollback.
		return ActionFailedReplace
	} else if e.currentRelease.Info.Status == rspb.StatusSuperseded {
		return ActionPostReplaceSuperseded
	}
	return ActionPreUninstall
}
