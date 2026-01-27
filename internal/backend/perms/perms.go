package perms

type Permission string

const (
	CreateMyModule        Permission = "CreateMyModule"
	DeleteModules         Permission = "DeleteModules"
	DeleteMyModule        Permission = "DeleteMyModule"
	SuggestMyModule       Permission = "SuggestMyModule"
	UpdateModules         Permission = "UpdateModules"
	DeleteReleases        Permission = "DeleteReleases"
	DeleteMyRelease       Permission = "DeleteMyRelease"
	UpdateReleases        Permission = "UpdateReleases"
	ChangeReleaseStatuses Permission = "ChangeReleaseStatuses"
	RejectReleases        Permission = "RejectReleases"
	AcceptReleases        Permission = "AcceptReleases"
	CancelReleases        Permission = "CancelReleases"
	BanUsers              Permission = "BanUsers"
	UnbanUsers            Permission = "UnbanUsers"
)

var permissionStrings = [...]string{
	"CreateMyModule",
	"DeleteModules",
	"DeleteMyModule",
	"SuggestMyModule",
	"UpdateModules",
	"DeleteReleases",
	"DeleteMyRelease",
	"UpdateReleases",
	"ChangeReleaseStatuses",
	"RejectReleases",
	"AcceptReleases",
	"CancelReleases",
	"BanUsers",
	"UnbanUsers",
}

func ContainsPerms(constraints []Permission, userPerms []Permission) bool {
	found := len(constraints)
	for i := range constraints {
		for j := range userPerms {
			if constraints[i] == userPerms[j] {
				found--
				break
			}
		}
	}

	return found == 0
}

func ContainsMappedPerms(constraints []Permission, userPerms map[string]bool) bool {
	for i := range constraints {
		_, ok := userPerms[string(constraints[i])]
		if !ok {
			return false
		}
	}
	return true
}
