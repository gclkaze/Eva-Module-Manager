package models

import "time"

type KeywordDTO struct {
	ID    uint
	Label string
}

type ReleaseDTO struct {
	ID          uint       `json:"id" binding:"required"`
	Version     string     `json:"version" binding:"required"`
	ReleasedAt  *time.Time `json:"released_at" binding:"required"`
	Description string     `json:"description" binding:"required"`
	DiskSize    int64      `json:"diskSize" binding:"required"`

	Status   string       `json:"status" binding:"required"`
	Keywords []KeywordDTO `json:"keywords" binding:"required"`
}

type ModuleEnrichedDTO struct {
	ID          uint   `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Repr        string `json:"repr" binding:"required"`
	Description string `json:"description" binding:"required"`
	//OwnerName   string       `json:"ownerName" binding:"required"`
	RepoName    string       `json:"repoName" binding:"required"`
	Tags        []string     `json:"tags" binding:"required"`
	ReleaseInfo []ReleaseDTO `json:"releases_info" binding:"required"`
}
