package services

import (
	"bytes"
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/models/userinput"
	"emm/internal/output"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/magiconair/properties"
)

type AuthService struct {
	keyringService *KeyringService
	props          *properties.Properties
	output         output.Printer
	currentUser    *models.User
	backend        *backend.Backend

	currentUserKey string
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
	inst.currentUserKey = props.GetString("currentUserKey", "EMMCurrentUser")
	inst.keyringService.SetProperties(props)
}

func (inst *AuthService) Register(creds *userinput.RegistrationCreds) (bool, error) {
	email := creds.Email
	if inst.keyringService.KeyExists(email) {
		inst.output.Info(fmt.Sprintf("User %s is already registered.", creds.Email))
		return true, nil
	}

	url, err := inst.backend.GetServerURL()
	if err != nil {
		return false, err
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

func (inst *AuthService) Login(creds *userinput.LoginCreds) (bool, error) {
	email := creds.Email
	if inst.keyringService.KeyExists(email) {
		inst.output.Info(fmt.Sprintf("User %s is already logged on.", creds.Email))
		inst.switchActiveUser(creds.Email)
		return true, nil
	}

	url, err := inst.backend.GetServerURL()
	if err != nil {
		return false, err
	}

	var buf *bytes.Buffer
	if creds != nil {
		b, _ := json.Marshal(creds)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s", url, APIGroup, AuthGroup, LoginEndpoint), buf)
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
		inst.output.Info(fmt.Sprintf("User %s is logged on successfully.", creds.Email))

		inst.switchActiveUser(creds.Email)

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

func (inst AuthService) getRefreshToken(email string) (string, error) {
	tk, err := inst.keyringService.LoadToken(email)
	if err != nil {
		return "", nil
	}

	toks := strings.Split(tk, " ")
	if len(toks) != 2 {
		return "", errors.New("incorrect token format")
	}
	return toks[1], nil
}

func (inst AuthService) combineToken(resp models.LoginResponse) string {
	return fmt.Sprintf("%s %s", resp.AccessToken, resp.RefreshToken)
}

func (inst *AuthService) refreshToken(email string) (bool, error) {

	url, err := inst.backend.GetServerURL()
	if err != nil {
		return false, err
	}

	reft, err := inst.getRefreshToken(email)
	if err != nil {
		return false, err
	}

	if reft == "" {
		return false, errors.New("empty ref token")
	}

	var buf *bytes.Buffer
	creds := models.NewRefreshRequest(reft)

	if creds != nil {
		b, _ := json.Marshal(creds)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s", url, APIGroup, AuthGroup, RefreshEndpoint), buf)
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
		inst.keyringService.SaveToken(email, inst.combineToken(response.Value))
		inst.switchActiveUser(email)
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

func (inst *AuthService) Logout(email string) error {
	if email == "" {
		//logout current user
		active, err := inst.GetActiveUser()
		if err != nil {
			return fmt.Errorf("no current active user")
		}

		//logout user
		err = inst.logoutUser(active)
		if err != nil {
			return err
		}

		//we need to clear the currentUser
		return inst.deleteActiveUser()
	}

	err := inst.logoutUser(email)
	if err != nil {
		return err
	}
	active, err := inst.GetActiveUser()
	if err != nil {
		return nil
	}

	if active == email {
		return inst.deleteActiveUser()
	}

	return nil
}

func (inst *AuthService) logoutUser(email string) error {
	if !inst.keyringService.KeyExists(email) {
		return fmt.Errorf("user with email %s is not logged in", email)
	}
	inst.output.Info(fmt.Sprintf("Logging out User %s", email))
	return inst.keyringService.DeleteToken(email)
}

/*func (inst *AuthService) SetActiveUser(email string) (string, error) {
	return "", nil
}*/

func (inst AuthService) switchActiveUser(email string) {
	inst.keyringService.SaveToken(inst.currentUserKey, email)
}

func (inst AuthService) deleteActiveUser() error {
	return inst.keyringService.DeleteToken(inst.currentUserKey)
}

func (inst AuthService) deleteActiveUserIfExists() {
	if inst.keyringService.KeyExists(inst.currentUserKey) {
		inst.keyringService.DeleteToken(inst.currentUserKey)
		return
	}
}

func (inst AuthService) SwitchCurrentActiveUser(email string) error {
	active, err := inst.GetActiveUser()
	if err == nil {
		if active == email {
			inst.output.Info(fmt.Sprintf("The current active User is %s already.", email))
			return nil
		}
	}

	//need to check if we know this user
	s, err := inst.keyringService.LoadToken(email)
	if err != nil {
		return fmt.Errorf("unknown user %s", email)
	}
	if s == "" {
		return fmt.Errorf("unknown user %s", email)
	}

	inst.switchActiveUser(email)
	inst.output.Info(fmt.Sprintf("The current active User is %s", email))
	return nil
}

func (inst AuthService) GetActiveUser() (string, error) {
	return inst.keyringService.LoadToken(inst.currentUserKey)
}

func (inst AuthService) ShowCurrentUser() error {
	active, err := inst.GetActiveUser()
	if err != nil {
		return fmt.Errorf("no current active user")
	}
	inst.output.Info(active)
	return nil
}

func (inst *AuthService) GetStatus() error {
	return nil
}
