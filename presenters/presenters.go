package presenters

import (
	"fmt"
	"os"

	"github.com/larderdev/kubewise/kwrelease"
)

func PrepareMsg(releaseEvent *kwrelease.Event) string {
	var msg string

	if value, ok := os.LookupEnv("KW_MESSAGE_PREFIX"); ok {
		msg += value
	}

	switch releaseEvent.GetAction() {
	case kwrelease.ActionPreInstall:
		msg += fmt.Sprintf("📀 Installing *%s* version *%s* into namespace *%s* via Helm. ⏳\n\n%s",
			releaseEvent.GetAppName(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetDescription(),
		)

	case kwrelease.ActionPreUpgrade:
		msg += fmt.Sprintf("⏫ Upgrading *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPreRollback:
		msg += fmt.Sprintf("⏬ Rolling back *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPreUninstall:
		msg += fmt.Sprintf("🧼 Uninstalling *%s* from namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPostInstall:
		msg += fmt.Sprintf("📀 Installed *%s* version *%s* into namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)

	case kwrelease.ActionPostUpgrade:
		msg += fmt.Sprintf("⏫ Upgraded *%s* from version %s to version *%s* in namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)

	case kwrelease.ActionPostRollback:
		msg += fmt.Sprintf("⏬ Rolled back *%s* from version %s to version *%s* in namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)

	case kwrelease.ActionPostReplace:
		msg += fmt.Sprintf("Replaced *%s* version %s with version *%s* in namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousAppVersion(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)
	}

	return msg
}
