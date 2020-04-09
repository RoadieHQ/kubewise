package presenters

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/larderdev/kubewise/kwrelease"
	"github.com/olekukonko/tablewriter"
	"helm.sh/helm/v3/pkg/release"
)

func getChangeInAppVersion(releaseEvent *kwrelease.Event) string {
	var appVersion string
	if releaseEvent.IsAppVersionChanged() {
		appVersion = fmt.Sprintf("App version will be *%s*, up from %s",
			releaseEvent.GetAppVersion(),
			releaseEvent.GetPreviousAppVersion(),
		)
	} else {
		appVersion = fmt.Sprintf("App version will be *%s* (unchanged)",
			releaseEvent.GetAppVersion(),
		)
	}
	return appVersion
}

func getConfigDiff(releaseEvent *kwrelease.Event) string {
	showDiff, err := strconv.ParseBool(os.Getenv("KW_CHART_VALUES_DIFF_ENABLED"))

	if err != nil {
		log.Println("Invalid value passed for environment variable KW_CHART_VALUES_DIFF_ENABLED. Boolean required.")
	}

	configDiffYAML := releaseEvent.GetConfigDiffYAML()
	if showDiff && configDiffYAML != "" {
		return fmt.Sprintf("\n```%s```", configDiffYAML)
	}

	return ""
}

// PrepareMsg prepares a short, markdown-like message which is suitable for sending to chat
// applications like Slack. Formatting like *text* us used to add emphasis. This is supported by
// both Slack and Google Chat. Emoji are also used liberally.
func PrepareMsg(releaseEvent *kwrelease.Event) string {
	msg := initializeServerStartupMsg()

	switch releaseEvent.GetAction() {
	case kwrelease.ActionPreInstall:
		msg += fmt.Sprintf("📀 Installing *%s* version *%s* into namespace *%s* via Helm. ⏳\n\nApp version: *%s*\n%s",
			releaseEvent.GetAppName(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetAppVersion(),
			releaseEvent.GetAppDescription(),
		)

	case kwrelease.ActionPreUpgrade:
		msg += fmt.Sprintf("⏫ Upgrading *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳\n%s",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousChartVersion(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
			getChangeInAppVersion(releaseEvent),
		)

		if configDiff := getConfigDiff(releaseEvent); configDiff != "" {
			msg += configDiff
		}

	case kwrelease.ActionPreRollback:
		msg += fmt.Sprintf("⏬ Rolling back *%s* from version %s to version *%s* in namespace *%s* via Helm. ⏳\n%s",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousChartVersion(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
			getChangeInAppVersion(releaseEvent),
		)

		if configDiff := getConfigDiff(releaseEvent); configDiff != "" {
			msg += configDiff
		}

	case kwrelease.ActionPreUninstall:
		msg += fmt.Sprintf("🧼 Uninstalling *%s* from namespace *%s* via Helm. ⏳",
			releaseEvent.GetAppName(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPostInstall:
		msg += fmt.Sprintf("📀 Installed *%s* version *%s* into namespace *%s* via Helm. ✅\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetNotes(),
		)

	case kwrelease.ActionPostUpgrade:
		msg += fmt.Sprintf("⏫ Upgraded *%s* from version %s to version *%s* in namespace *%s* via Helm. ✅",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousChartVersion(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPostRollback:
		msg += fmt.Sprintf("⏬ Rolled back *%s* from version %s to version *%s* in namespace *%s* via Helm. ✅",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousChartVersion(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionPostReplace:
		msg += fmt.Sprintf("Replaced *%s* version %s with version *%s* in namespace *%s* via Helm. ✅",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousChartVersion(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
		)

	case kwrelease.ActionFailedInstall:
		msg += fmt.Sprintf("❌ Installation of *%s* version *%s* in namespace *%s* has FAILED. ❌\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
			// This has the cause of the failure.
			releaseEvent.GetReleaseDescription(),
		)

	case kwrelease.ActionFailedReplace:
		msg += fmt.Sprintf("❌ Replacing *%s* version %s with version *%s* in namespace *%s* has FAILED. ❌\n\n```%s```",
			releaseEvent.GetAppName(),
			releaseEvent.GetPreviousChartVersion(),
			releaseEvent.GetChartVersion(),
			releaseEvent.GetNamespace(),
			releaseEvent.GetReleaseDescription(),
		)
	}

	return msg
}

func initializeServerStartupMsg() string {
	var msg string

	if value, ok := os.LookupEnv("KW_MESSAGE_PREFIX"); ok {
		msg += value
	}

	return msg
}

// PrepareServerStartupMsg prepares a message which is suitable for sending to a chat application
// like Slack on server startup. The message will contain information about the Helm charts that
// are installed in the cluster at the time of install. They will be presented in a monospaced
// table.
func PrepareServerStartupMsg(releases []*release.Release) string {
	msg := initializeServerStartupMsg()
	numberOfReleases := len(releases)
	msg += "👋 KubeWise initialized."

	if numberOfReleases == 1 {
		msg += " There is *1* Helm chart installed."
	} else {
		msg += fmt.Sprintf(" There are *%s* Helm charts installed.",
			strconv.Itoa(numberOfReleases),
		)
	}

	if numberOfReleases > 0 {
		msg += fmt.Sprintf("```%s```", renderTableShowingInstalledCharts(releases))
	}

	return msg
}

func renderTableShowingInstalledCharts(releases []*release.Release) string {
	data := make([][]string, len(releases))
	for i, release := range releases {
		data[i] = []string{
			release.Name,
			release.Chart.AppVersion(),
			release.Chart.Metadata.Version,
		}
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"App Name", "App Version", "Chart Version"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return tableString.String()
}
