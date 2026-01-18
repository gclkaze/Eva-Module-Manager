package models

type Module struct {
	ID          uint
	RepoName    string
	Title       string
	Repr        string
	Description string
	Releases    int
	Tags        []string
}

type ModuleEnrichedInformation struct {
	ID          uint   `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Repr        string `json:"repr" binding:"required"`
	Description string `json:"description" binding:"required"`
	//OwnerName   string       `json:"ownerName" binding:"required"`
	RepoName    string    `json:"repoName" binding:"required"`
	Tags        []string  `json:"tags" binding:"required"`
	ReleaseInfo []Release `json:"releases_info" binding:"required"`
}
