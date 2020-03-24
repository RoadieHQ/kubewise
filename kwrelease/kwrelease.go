package kwrelease

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/larderdev/kubewise/driver"
	"github.com/larderdev/kubewise/utils"
	"helm.sh/helm/v3/pkg/release"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Event struct {
	// Describes the action that happened to the secret in order to trigger this event. Events
	// occur when secrets are created, updated or deleted. What was the action that led to this
	// particular event.
	SecretAction         string
	CurrentReleaseSecret *api_v1.Secret
	currentRelease       *release.Release
	previousRelease      *release.Release
}

func (e *Event) Init() error {
	currentRelease, err := driver.DecodeRelease(string(e.CurrentReleaseSecret.Data["release"]))

	if err != nil {
		log.Fatalln("Error getting releaseData from secret", e.CurrentReleaseSecret)
		return err
	}
	e.currentRelease = currentRelease
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

func (e *Event) GetDescription() string {
	return e.currentRelease.Chart.Metadata.Description
}

func (e *Event) GetNotes() string {
	return e.currentRelease.Info.Notes
}

func (e *Event) GetSecretUID() types.UID {
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
	if e.currentRelease.Info.Status == release.StatusPendingInstall {
		return ActionPreInstall
	} else if e.currentRelease.Info.Status == release.StatusPendingUpgrade {
		return ActionPreUpgrade
	} else if e.currentRelease.Info.Status == release.StatusPendingRollback {
		return ActionPreRollback
	} else if e.currentRelease.Info.Status == release.StatusDeployed {
		if e.previousRelease == nil {
			return ActionPostInstall
		} else if strings.HasPrefix(e.currentRelease.Info.Description, "Rollback") {
			return ActionPostRollback
		} else if strings.HasPrefix(e.currentRelease.Info.Description, "Upgrade") {
			return ActionPostUpgrade
		}
		return ActionPostReplace
	} else if e.currentRelease.Info.Status == release.StatusSuperseded {
		return ActionPostReplaceSuperseded
	}
	return ActionPreUninstall
}

func inferNameOfPreviousReleaseSecret(currentReleaseSecretName string) string {
	log.Println("Finding previous releases of release name", currentReleaseSecretName)
	// e.g. [sh helm release v1 airflow v2]
	currentReleaseVersion := strings.Split(currentReleaseSecretName, ".")
	previousReleaseVersion := make([]string, len(currentReleaseVersion))
	copy(previousReleaseVersion, currentReleaseVersion)
	// e.g. 2
	currentReleaseNumber, atoiErr := strconv.Atoi(currentReleaseVersion[len(currentReleaseVersion)-1][1:])

	if currentReleaseNumber <= 1 {
		log.Println("First release of", currentReleaseSecretName, "in Helm history. Not checking for previous releases.")
		return ""
	}

	if atoiErr != nil {
		log.Println("Error parsing current release number from string:", atoiErr)
		return ""
	}

	// e.g. "v1"
	previousReleaseVersion[len(previousReleaseVersion)-1] = fmt.Sprintf("v%v", currentReleaseNumber-1)
	// e.g. sh.helm.release.v1.airflow.v1
	return strings.Join(previousReleaseVersion, ".")
}

// GetPreviousRelease locates previous Helm releasees based off the name of a given secret.
// Helm 3 releases have secret names like: sh.helm.release.v1.zookeeper.v1
// The last `v1` is incremented every time an upgrade occurs. This way, the previous releases
// of a package can be located by looking for secrets with matching names.
//
// When Helm upgrades a package, it leaves secrets from the previous releases in the cluster. They
// can be located and used to determine if the current operation is an install or an upgrade. It
// is also useful to inform the user of the appVersion being upgraded from.
func (e *Event) getPreviousRelease() *release.Release {
	previousReleaseSecretName := inferNameOfPreviousReleaseSecret(e.CurrentReleaseSecret.Name)
	if previousReleaseSecretName == "" {
		return nil
	}

	log.Println("Finding previous release with name:", previousReleaseSecretName)

	kubeClient := utils.GetClient()
	getOptions := meta_v1.GetOptions{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
	}
	previousReleaseSecret, err := kubeClient.CoreV1().Secrets(e.CurrentReleaseSecret.Namespace).Get(previousReleaseSecretName, getOptions)

	if err != nil {
		log.Println("Error finding previous release secret with name", previousReleaseSecretName)
		return nil
	}

	previousRelease, err := driver.DecodeRelease(string(previousReleaseSecret.Data["release"]))

	if err != nil {
		log.Println("Error decoding previous Helm release", previousReleaseSecret)
		return nil
	}

	return previousRelease
}
