package kwrelease

// Action describes the action that was taken by Helm. For example, it could be an install
// or an uninstall (delete) action. There are typically two phases to each action, pre and
// post. This can best be demonstrated by using helm ... --wait. Helm will create a release
// to mark the start of the action (pre) and update it to mark the end of the action (post).
type Action string

// All possible actions and their states. See helm/pkg/release/status.go for more details.
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
