package userinput

import "time"

type ReleaseFilterParams struct {
	Status        []string
	Versions      []string
	Tags          []string
	ModuleName    []string
	RepoName      []string
	CreatedAfter  *time.Time
	ReleasedAfter time.Time
	Description   []string
	Creator       []string
	CreatorEmail  []string
}

func NewReleaseFilterParams() *ReleaseFilterParams {
	return &ReleaseFilterParams{}
}

func (inst ReleaseFilterParams) IsEmpty() bool {
	if len(inst.Status) > 0 ||
		len(inst.Versions) > 0 ||
		len(inst.Tags) > 0 ||
		len(inst.ModuleName) > 0 ||
		len(inst.RepoName) > 0 ||
		len(inst.Description) > 0 ||
		len(inst.Creator) > 0 ||
		len(inst.CreatorEmail) > 0 {
		return false
	}

	// time filters
	if inst.CreatedAfter != nil {
		return false
	}
	if !inst.ReleasedAfter.IsZero() {
		return false
	}

	return true
}
