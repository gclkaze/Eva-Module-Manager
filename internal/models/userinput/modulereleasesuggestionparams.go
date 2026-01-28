package userinput

import (
	"emm/pkg/utils"
	"fmt"
	"strings"
)

type ModuleReleaseSuggestionParams struct {
	ModuleRepr        string
	Version           string
	Description       string
	TagsCSV           string
	InheritModuleTags bool
}

func NewModuleReleaseSuggestionParams() *ModuleReleaseSuggestionParams {
	return &ModuleReleaseSuggestionParams{}
}
func (inst *ModuleReleaseSuggestionParams) AllValid() error {
	err := utils.IsValidRepoName(inst.ModuleRepr)
	if err != nil {
		return err
	}

	if !utils.IsValidVersion(inst.Version) {
		return fmt.Errorf("invalid module release version suggestion: '%s'", inst.Version)
	}
	// description (empty allowed)
	if err := utils.IsValidDescription(inst.Description); err != nil {
		return err
	}

	// tags (empty allowed)
	tags, err := utils.ParseTagsCSV(inst.TagsCSV)
	if err != nil {
		return err
	}

	// normalize CSV so what you send is consistent (lowercase, trimmed, deduped)
	inst.TagsCSV = strings.Join(tags, ",")
	return nil
}
