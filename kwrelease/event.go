package kwrelease

import (
	"log"
	"strconv"
	"strings"

	rspb "helm.sh/helm/v3/pkg/release"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kbtypes "k8s.io/apimachinery/pkg/types"
)

// This file contains most of the simple Getters and Setters for the Event struct. There are
// other, more complex methods elsewhere.

type Event struct {
	// Describes the action that happened to the secret in order to trigger this event. Events
	// occur when secrets are created, updated or deleted. What was the action that led to this
	// particular event.
	SecretAction         string
	CurrentReleaseSecret *api_v1.Secret
	currentRelease       *rspb.Release
	previousRelease      *rspb.Release
}

func (e *Event) Init() error {
	// Fetching the release from the secret store is unnecessary except for the fact that we need
	// to decode it and the safest way to do that is to let the Helm lib do it. THe Helm code has
	// a function called driver.decodeRelease but it is private and cannot be accessed directly.
	// By Getting the release from the store we get it back decoded for us.
	e.currentRelease = e.GetRelease(e.CurrentReleaseSecret.Name)
	e.previousRelease = e.getPreviousRelease()

	return nil
}

func (e *Event) GetAppName() string {
	return e.currentRelease.Name
}

func (e *Event) GetAppVersion() string {
	return e.currentRelease.Chart.AppVersion()
}

func (e *Event) GetPreviousAppVersion() string {
	if e.previousRelease != nil {
		return e.previousRelease.Chart.AppVersion()
	}
	return ""
}

func (e *Event) GetNamespace() string {
	return e.currentRelease.Namespace
}

func (e *Event) GetAppDescription() string {
	return e.currentRelease.Chart.Metadata.Description
}

func (e *Event) GetReleaseDescription() string {
	return e.currentRelease.Info.Description
}

func (e *Event) GetNotes() string {
	return e.currentRelease.Info.Notes
}

func (e *Event) GetSecretUID() kbtypes.UID {
	return e.CurrentReleaseSecret.GetUID()
}

func (e *Event) GetSecretCreationTimestamp() meta_v1.Time {
	return e.CurrentReleaseSecret.GetCreationTimestamp()
}

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

func (e *Event) GetChartVersion() string {
	return e.currentRelease.Chart.Metadata.Version
}

func (e *Event) GetPreviousChartVersion() string {
	if e.previousRelease != nil {
		return e.previousRelease.Chart.Metadata.Version
	}
	return ""
}

func (e *Event) IsAppVersionChanged() bool {
	return e.GetAppVersion() != e.GetPreviousAppVersion()
}

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
