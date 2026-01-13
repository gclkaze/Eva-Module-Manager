package userinput

import (
	"emm/pkg/utils"
	"fmt"
)

type RegistrationCreds struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Handle    string
}

func NewRegistrationCreds() *RegistrationCreds {
	return &RegistrationCreds{}
}

func (inst RegistrationCreds) AllInformationProvided() bool {
	if len(inst.Email) == 0 || len(inst.Password) == 0 || len(inst.FirstName) == 0 || len(inst.LastName) == 0 || len(inst.Handle) == 0 {
		return false
	}
	return true
}

func (inst RegistrationCreds) AllInformationProvidedExceptPassword() bool {
	if len(inst.Email) == 0 || len(inst.FirstName) == 0 || len(inst.LastName) == 0 || len(inst.Handle) == 0 {
		return false
	}
	return true
}

func (inst RegistrationCreds) AreValid() error {
	//check email
	if !utils.IsValidEmail(inst.Email) {
		return fmt.Errorf("invalid email provided: %s", inst.Email)
	}
	//check password
	err := utils.ValidatePassword(inst.Password, inst.Email)
	if err != nil {
		return err
	}
	//check firstname
	err = utils.IsValidNameWithError("first name", inst.FirstName)
	if err != nil {
		return err
	}
	//check lastname
	err = utils.IsValidNameWithError("last name", inst.LastName)
	if err != nil {
		return err
	}

	//check handle
	err = utils.IsValidHandleWithError("handle", inst.Handle)
	if err != nil {
		return err
	}
	return nil
}
