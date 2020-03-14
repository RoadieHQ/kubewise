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

func GetPreviousRelease(releaseSecret *api_v1.Secret) *release.Release {
	previousReleaseSecretName := inferNameOfPreviousReleaseSecret(releaseSecret.Name)
	log.Println("Finding previous release with name:", previousReleaseSecretName)
	if previousReleaseSecretName == "" {
		return nil
	}

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
	}

	previousRelease, err := driver.DecodeRelease(string(previousReleaseSecret.Data["release"]))

	if err != nil {
		log.Println("Error decoding previous Helm release", previousReleaseSecret)
		return nil
	}

	return previousRelease
}
