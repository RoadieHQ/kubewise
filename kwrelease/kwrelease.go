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
)

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
func GetPreviousRelease(releaseSecret *api_v1.Secret) *release.Release {
	previousReleaseSecretName := inferNameOfPreviousReleaseSecret(releaseSecret.Name)
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
	previousReleaseSecret, err := kubeClient.CoreV1().Secrets(releaseSecret.Namespace).Get(previousReleaseSecretName, getOptions)

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
