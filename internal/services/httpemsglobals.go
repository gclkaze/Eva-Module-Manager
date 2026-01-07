package services

const (
	APIGroup         = "/api"
	AuthGroup        = "/auth"
	LoginEndpoint    = "/login"
	RefreshEndpoint  = "/refresh"
	RegisterEndpoint = "/register"

	ModulesGroup          = "/modules"
	ModuleSearchEndpoint  = "/search"
	ModuleGetInfoEndpoint = "/info"
	ModuleDeleteEndpoint  = "/delete"
	ModuleUploadEndpoint  = "/upload"
	ModuleUpdateEndpoint  = "/update"
	ModuleSuggestEndpoint = "/suggest"

	ReleasesGroup         = "/releases"
	ReleaseEndpoint       = "release"
	ReleaseSearchEndpoint = "/search"
	ReleaseDeleteEndpoint = "/delete"
	ReleaseCancelEndpoint = "/cancel"

	DownloadGroup           = "/download"
	DownloadReleaseEndpoint = "/release"

	SuperviseGroup                      = "/supervise"
	SuperviseDownloadAnyReleaseEndpoint = "/download/release"
	SuperviseRejectReleaseEndpoint      = "/reject/release"
	SuperviseAcceptReleaseEndpoint      = "/accept/release"
	SuperviseCancelReleaseEndpoint      = "/cancel/release"
	SupervisePendingReleaseEndpoint     = "/pending/release"
	SuperviseBanUserEndpoint            = "/ban"
	SuperviseUnbanUserEndpoint          = "/unban"
)
