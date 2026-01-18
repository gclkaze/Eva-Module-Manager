package userinput

import (
	"emm/pkg/utils"
	"fmt"
)

type UploadParams struct {
	Title       string
	Repr        string
	Tags        string
	Description string
}

func NewUploadParams() *UploadParams {
	return &UploadParams{}
}

func (inst UploadParams) AllValid() error {
	if len(inst.Title) == 0 {
		return fmt.Errorf("the module title is empty")
	}

	if len(inst.Repr) == 0 {
		return fmt.Errorf("the module representation name is empty")
	}

	if !utils.IsValidModuleName(inst.Repr) {
		return fmt.Errorf("the module representation name is not correct...need to be between %d and %d characters, only digits and letters are accepted", utils.ModuleReprMin, utils.ModuleReprMax)
	}

	return nil
}
