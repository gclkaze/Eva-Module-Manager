package userinput

import (
	"emm/pkg/utils"
	"fmt"
)

type LoginCreds struct {
	Email    string
	Password string
}

func NewLoginCreds() *LoginCreds {
	return &LoginCreds{}
}

func (inst LoginCreds) AllInformationProvided() bool {
	if len(inst.Email) == 0 || len(inst.Password) == 0 {
		return false
	}
	return true
}

func (inst LoginCreds) AllInformationProvidedExceptPassword() bool {
	return len(inst.Email) != 0
}

func (inst LoginCreds) AreValid() error {
	//check email
	if !utils.IsValidEmail(inst.Email) {
		return fmt.Errorf("invalid email provided: %s", inst.Email)
	}
	//check password
	err := utils.ValidatePassword(inst.Password, inst.Email)
	if err != nil {
		return err
	}
	return nil
}
