package kwrelease

type Action string

const (
	ActionPreInstall   Action = "PRE_INSTALL"
	ActionPreRollback  Action = "PRE_ROLLBACK"
	ActionPreUpgrade   Action = "PRE_UPGRADE"
	ActionPostInstall  Action = "POST_INSTALL"
	ActionPostUpgrade  Action = "POST_UPGRADE"
	ActionPostRollback Action = "POST_ROLLBACK"
	// Because string matching is a dangerous game, we should have a fallback for times we
	// can't tell if the operation was an Upgrade or a Rollback.
	ActionPostReplace           Action = "POST_REPLACE"
	ActionPostReplaceSuperseded Action = "POST_REPLACE-SUPERSEDED"
	ActionPreUninstall          Action = "PRE_UNINSTALL"
	ActionFailedInstall         Action = "FAILED_INSTALL"
	ActionFailedReplace         Action = "FAILED_REPLACE"
)

func (a Action) String() string {
	return string(a)
}
