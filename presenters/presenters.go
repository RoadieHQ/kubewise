package presenters

import (
	"fmt"

	"helm.sh/helm/v3/pkg/release"
)

func PrepareObjectCreatedMsg(currentRelease, previousRelease *release.Release) string {
	var msg string

	if currentRelease.Info.Status == release.StatusPendingInstall {
		msg = fmt.Sprintf(
			"💽 Installing *%s* version *%s* into namespace *%s* via Helm. ⏳\n\n%s",
			currentRelease.Name,
			currentRelease.Chart.AppVersion(),
			currentRelease.Namespace,
			currentRelease.Info.Description,
		)
	} else if currentRelease.Info.Status == release.StatusPendingUpgrade {
		if previousRelease == nil {
			msg = fmt.Sprintf(
				"⏫ Upgrading *%s* to version *%s* in namespace *%s* via Helm. ⏳",
				currentRelease.Name,
				currentRelease.Chart.AppVersion(),
				currentRelease.Namespace,
			)
		} else {
			msg = fmt.Sprintf(
				"⏫ Upgrading *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳",
				currentRelease.Name,
				previousRelease.Chart.AppVersion(),
				currentRelease.Chart.AppVersion(),
				currentRelease.Namespace,
			)
		}
	} else if currentRelease.Info.Status == release.StatusPendingRollback {
		if previousRelease == nil {
			msg = fmt.Sprintf(
				"⏬ Rolling back *%s* from version %s in namespace *%s* via Helm. ⏳",
				currentRelease.Name,
				currentRelease.Chart.AppVersion(),
				currentRelease.Namespace,
			)
		} else {
			msg = fmt.Sprintf(
				"⏫ Rolling back *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳",
				currentRelease.Name,
				previousRelease.Chart.AppVersion(),
				currentRelease.Chart.AppVersion(),
				currentRelease.Namespace,
			)
		}
	}

	return msg
}

func PrepareObjectDeletedMsg(currentRelease, previousRelease *release.Release) string {
	return ""
}

func PrepareObjectUpgradedMsg(currentRelease, previousRelease *release.Release) string {
	var msg string

	if currentRelease.Info.Status == release.StatusDeployed {
		if previousRelease == nil {
			msg = fmt.Sprintf(
				"💽 Successfully installed *%s* version *%s* into namespace *%s* via Helm. ✅\n\n```%s```",
				currentRelease.Name,
				currentRelease.Chart.AppVersion(),
				currentRelease.Namespace,
				currentRelease.Info.Notes,
			)
		} else {
			// There is no way to tell if this is an upgrade or a rollback. The previous release
			// status will be changed to "superseeded" and the new release will have the status
			// "deployed". Versions are arbitrary strings so we can't compare them against each
			// other.
			// Best I can do is use neutral language like "replaced".
			msg = fmt.Sprintf(
				"⏫ Successfully replaced *%s* version %s with version *%s* in namespace *%s* via Helm. ✅\n\n```%s```",
				currentRelease.Name,
				previousRelease.Chart.AppVersion(),
				currentRelease.Chart.AppVersion(),
				currentRelease.Namespace,
				currentRelease.Info.Notes,
			)
		}
	} else if currentRelease.Info.Status == release.StatusUninstalling {
		msg = fmt.Sprintf(
			"🧼 Uninstalling *%s* version *%s* from namespace %s via Helm.",
			currentRelease.Name,
			currentRelease.Chart.AppVersion(),
			currentRelease.Namespace,
		)
	}

	return msg
}
