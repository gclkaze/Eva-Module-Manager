package userinput

import "emm/pkg/utils"

type ModuleSearchQuery struct {
	Tags        []string
	Name        []string
	Description []string
}

func NewModuleSearchQuery() *ModuleSearchQuery {
	return &ModuleSearchQuery{}
}

func (inst ModuleSearchQuery) IsEmpty() bool {
	return len(inst.Description) == 0 && len(inst.Name) == 0 && len(inst.Tags) == 0
}
func (inst *ModuleSearchQuery) Normalize() {
	inst.Tags = utils.NormalizeTags(inst.Tags)
	inst.Name = utils.NormalizeText(inst.Name)
	inst.Description = utils.NormalizeText(inst.Description)
}

// IsValid validates AFTER normalization
func (inst ModuleSearchQuery) IsValid() error {
	if err := utils.ValidateTags(inst.Tags); err != nil {
		return err
	}
	if err := utils.ValidateNames(inst.Name); err != nil {
		return err
	}
	if err := utils.ValidateDescriptions(inst.Description); err != nil {
		return err
	}
	return nil
}
