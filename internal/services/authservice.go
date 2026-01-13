package services

import (
	"bytes"
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/models/userinput"
	"emm/internal/output"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/magiconair/properties"
)

type AuthService struct {
	keyringService *KeyringService
	props          *properties.Properties
	output         output.Printer
	currentUser    *models.User
	backend        *backend.Backend
}

func NewAuthService() *AuthService {
	return &AuthService{keyringService: NewKeyringService()}
}

func (inst *AuthService) SetCurrentUser() {

}

func (inst *AuthService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *AuthService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst *AuthService) SetProperties(props *properties.Properties) {
	inst.props = props
	inst.keyringService.SetProperties(props)
}

func (inst *AuthService) Register(creds *userinput.RegistrationCreds) (bool, error) {
	url, err := inst.backend.GetServerURL()
	if err != nil {
		return false, err
	}

	email := creds.Email
	if inst.keyringService.KeyExists(email) {
		inst.output.Info(fmt.Sprintf("User %s is already registered.", creds.Email))
		return true, nil
	}

	var buf *bytes.Buffer
	if creds != nil {
		b, _ := json.Marshal(creds)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s", url, APIGroup, AuthGroup, RegisterEndpoint), buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var response models.RequestResult[models.LoginResponse]
		err = json.Unmarshal(body, &response)
		if err != nil {
			inst.output.Error(err)
			return false, err
		}
		inst.keyringService.SaveToken(creds.Email, inst.combineToken(response.Value))
		inst.output.Info(fmt.Sprintf("User %s is registered successfully.", creds.Email))
		return true, nil
	}
	var response models.ErrorResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}
	err = fmt.Errorf("%s", response.Error)
	return false, err
}

func (inst AuthService) combineToken(resp models.LoginResponse) string {
	return fmt.Sprintf("%s %s", resp.AccessToken, resp.RefreshToken)
}

func (inst *AuthService) Login(email string, pwd string) (string, error) {
	return "", nil
}

func (inst *AuthService) Logout(email string) (string, error) {
	return "", nil
}

/*func (inst *AuthService) SetActiveUser(email string) (string, error) {
	return "", nil
}*/

func (inst *AuthService) SwitchActiveUser(email string) (string, error) {
	return "", nil
}

func (inst *AuthService) GetStatus() error {
	return nil
}
