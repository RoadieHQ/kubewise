package kwrelease

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/larderdev/kubewise/utils"
	rspb "helm.sh/helm/v3/pkg/release"
	helmdriver "helm.sh/helm/v3/pkg/storage/driver"
)

func inferNameOfPreviousReleaseSecret(currentReleaseSecretName string) string {
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

// GetRelease retrieves a release object from the Kubernetes Secret store. It delegates to
// the Helm Driver for this operation in order to reduce the possibility of breaking changes.
func (e *Event) GetRelease(secretName string) *rspb.Release {
	kubeClient := utils.GetClient()
	secrets := helmdriver.NewSecrets(kubeClient.CoreV1().Secrets(e.CurrentReleaseSecret.Namespace))
	result, err := secrets.Get(secretName)

	if err != nil {
		log.Println("Error finding release secret with name:", secretName)
		return nil
	}

	return result
}

// GetPreviousRelease locates previous Helm releases based off the name of a given secret.
// Helm 3 releases have secret names like: sh.helm.release.v1.zookeeper.v1
// The last `v1` is incremented every time an upgrade occurs. This way, the previous releases
// of a package can be located by looking for secrets with matching names.
//
// When Helm upgrades a package, it leaves secrets from the previous releases in the cluster. They
// can be located and used to determine if the current operation is an install or an upgrade. It
// is also useful to inform the user of the appVersion being upgraded from.
func (e *Event) getPreviousRelease() *rspb.Release {
	previousReleaseSecretName := inferNameOfPreviousReleaseSecret(e.CurrentReleaseSecret.Name)
	if previousReleaseSecretName == "" {
		return nil
	}

	log.Println("Finding previous release with name:", previousReleaseSecretName)
	return e.GetRelease(previousReleaseSecretName)
}

// ListActiveReleases lists releases which have not been superseded by an upgrade, rollback or
// other operation.
func ListActiveReleases() []*rspb.Release {
	kubeClient := utils.GetClient()

	namespace := ""
	if value, ok := os.LookupEnv("KW_NAMESPACE"); ok {
		namespace = value
	}

	secrets := helmdriver.NewSecrets(kubeClient.CoreV1().Secrets(namespace))
	results, err := secrets.List(func(r *rspb.Release) bool {
		return r.Info.Status != rspb.StatusSuperseded
	})

	if err != nil {
		log.Println("Error finding release secrets")
		return nil
	}

	return results
}
