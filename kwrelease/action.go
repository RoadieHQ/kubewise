package kwrelease

type Action string

const (
	ActionPreInstall            Action = "PRE_INSTALL"
	ActionPreReplaceUpgrade     Action = "PRE_REPLACE-UPGRADE"
	ActionPreReplaceRollback    Action = "PRE_REPLACE-ROLLBACK"
	ActionPostInstall           Action = "POST_INSTALL"
	ActionPostReplace           Action = "POST_REPLACE"
	ActionPostReplaceSuperseded Action = "POST_REPLACE-SUPERSEDED"
	ActionPreUninstall          Action = "PRE_UNINSTALL"
)

func (a Action) String() string {
	return string(a)
}
