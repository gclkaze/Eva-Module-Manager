package userinput

import (
	"emm/pkg/utils"
	"fmt"
)

type ModuleReleaseSuggestionParams struct {
	ModuleRepr string
	Version    string
}

func NewModuleReleaseSuggestionParams() *ModuleReleaseSuggestionParams {
	return &ModuleReleaseSuggestionParams{}
}
func (inst ModuleReleaseSuggestionParams) AllValid() error {
	err := utils.IsValidRepoName(inst.ModuleRepr)
	if err != nil {
		return err
	}

	if !utils.IsValidVersion(inst.Version) {
		return fmt.Errorf("invalid module release version suggestion: '%s'", inst.Version)
	}

	return nil
}
