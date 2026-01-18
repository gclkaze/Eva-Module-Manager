package userinput

import (
	"emm/pkg/utils"
	"fmt"
)

type UploadModuleUpdateParams struct {
	Title       string
	Repr        string
	Tags        string
	Description string

	ModuleRepr string
}

func NewUploadModuleUpdateParams() *UploadModuleUpdateParams {
	return &UploadModuleUpdateParams{}
}

func (inst UploadModuleUpdateParams) AllValid() error {
	if len(inst.ModuleRepr) == 0 {
		return fmt.Errorf("the module title is empty")
	}
	if len(inst.Title) == 0 {
		return fmt.Errorf("the module title is empty")
	}

	if len(inst.Repr) == 0 {
		return fmt.Errorf("the module representation name is empty")
	}

	err := utils.IsValidRepoName(inst.ModuleRepr)
	if err != nil {
		return err
	}

	return nil
}
