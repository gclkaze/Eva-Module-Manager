package services

import (
	"github.com/magiconair/properties"
	"github.com/zalando/go-keyring"
)

type KeyringService struct {
	props       *properties.Properties
	serviceName string
}

func NewKeyringService() *KeyringService {
	return &KeyringService{}
}

func (inst *KeyringService) SetProperties(props *properties.Properties) {
	inst.props = props
	inst.serviceName = props.GetString("app_name", "emm")
}

func (inst KeyringService) SaveToken(key string, token string) error {
	return keyring.Set(inst.serviceName, key, token)
}

func (inst KeyringService) LoadToken(key string) (string, error) {
	return keyring.Get(inst.serviceName, key)
}

func (inst KeyringService) KeyExists(key string) bool {
	_, err := inst.LoadToken(key)
	return err == nil
}

func (inst KeyringService) DeleteToken(key string) error {
	return keyring.Delete(inst.serviceName, key)
}
