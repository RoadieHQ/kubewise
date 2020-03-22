package presenters

import (
	"fmt"

	"github.com/larderdev/kubewise/kwrelease"
)

func PrepareMsg(releaseEvent *kwrelease.Event) string {
	switch releaseEvent.GetAction() {
	case kwrelease.ActionPreInstall:
		return fmt.Sprintf("💽 Installing *%s* version *%s* into namespace *%s* via Helm. ⏳\n\n%s",
			releaseEvent.GetAppName(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetDescription(),
		)

	case kwrelease.ActionPreReplaceUpgrade:
		return fmt.Sprintf("⏫ Upgrading *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPreReplaceRollback:
		return fmt.Sprintf("⏬ Rolling back *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPreUninstall:
		return fmt.Sprintf("🧼 Uninstalling *%s* from namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPostInstall:
		return fmt.Sprintf("Installed *%s* version *%s* into namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)

	case kwrelease.ActionPostReplace:
		return fmt.Sprintf("Replaced *%s* version %s with version *%s* in namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)

	default:
		return ""
	}
}
