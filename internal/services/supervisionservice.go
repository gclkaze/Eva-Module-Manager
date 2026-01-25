package services

import (
	"bytes"
	"context"
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/models/dto"
	"emm/internal/output"
	"emm/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/magiconair/properties"
)

type SupervisionService struct {
	authService *AuthService
	props       *properties.Properties
	output      output.Printer
	currentUser *models.User
	backend     *backend.Backend

	currentUserKey string
}

func NewSupervisionService(authService *AuthService) *SupervisionService {
	return &SupervisionService{authService: authService}
}
func (inst *SupervisionService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *SupervisionService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst *SupervisionService) SetProperties(props *properties.Properties) {
	inst.props = props
}

func (inst *SupervisionService) UserBanByID(ctx context.Context, token string, userID uint) error {
	res, err := inst.applyBan(ctx, token, userID, SuperviseBanUserEndpoint)
	if err != nil {
		inst.output.Error(err)
		inst.output.Error(fmt.Errorf("couldn't ban user"))
		return nil
	}
	if !res {
		inst.output.Error(fmt.Errorf("couldn't ban user"))
		return nil
	}

	inst.output.Info(fmt.Sprintf("User %d was banned successfully", userID))
	return nil
}

func (inst *SupervisionService) UserUnbanByID(ctx context.Context, token string, userID uint) error {
	res, err := inst.applyBan(ctx, token, userID, SuperviseUnbanUserEndpoint)
	if err != nil {
		inst.output.Error(err)
		inst.output.Error(fmt.Errorf("couldn't unban user"))
		return nil
	}
	if !res {
		inst.output.Error(fmt.Errorf("couldn't unban user"))
		return nil
	}

	inst.output.Info(fmt.Sprintf("User %d was unbanned successfully", userID))
	return nil
}

func (inst *SupervisionService) UserBan(ctx context.Context, token string, email string) error {
	user, err := inst.GetUser(ctx, token, email)
	if err != nil {
		inst.output.Error(err)
		return fmt.Errorf("couldn't find user '%s'", email)
	}

	if user == nil {
		inst.output.Error(err)
		return nil
	}

	if user.IsBanned {
		inst.output.Info(fmt.Sprintf("User %s is already banned", email))
		return nil
	}

	res, err := inst.applyBan(ctx, token, user.ID, SuperviseBanUserEndpoint)
	if err != nil {
		inst.output.Error(err)
		inst.output.Error(fmt.Errorf("couldn't ban user"))
		return nil
	}
	if !res {
		inst.output.Error(fmt.Errorf("couldn't ban user"))
		return nil
	}

	inst.output.Info(fmt.Sprintf("User %s was banned successfully", email))
	return nil
}

func (inst *SupervisionService) UserUnban(ctx context.Context, token string, email string) error {
	user, err := inst.GetUser(ctx, token, email)
	if err != nil {
		inst.output.Error(err)
		return fmt.Errorf("couldn't find user '%s'", email)
	}

	if user == nil {
		inst.output.Error(err)
		return nil
	}

	if !user.IsBanned {
		inst.output.Info(fmt.Sprintf("User %s is not banned", email))
		return nil
	}

	res, err := inst.applyBan(ctx, token, user.ID, SuperviseUnbanUserEndpoint)
	if err != nil {
		inst.output.Error(err)
		inst.output.Error(fmt.Errorf("couldn't unban user"))
		return nil
	}
	if !res {
		inst.output.Error(fmt.Errorf("couldn't unban user"))
		return nil
	}

	inst.output.Info(fmt.Sprintf("User %s was unbanned successfully", email))
	return nil
}

func (inst *SupervisionService) GetUsers(ctx context.Context, token string) error {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s%s", url, APIGroup, SuperviseGroup, SuperviseGetUsersEndpoint), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/octet-stream")
		req.Header.Set("Authorization", "Bearer "+theToken)

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	resp, err := inst.authService.PerformSafeCall(token, cb)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var response models.ErrorResult
		err := json.Unmarshal(body, &response)
		if err != nil {
			return err
		}
		return fmt.Errorf("users info fetch failed: %s", response.Details)
	}

	var response models.RequestResult[[]dto.DeveloperDTO]
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	myEmail, err := inst.authService.GetActiveUser()
	if err != nil {
		return err
	}
	inst.output.PrintDevelopers(response.Value, myEmail)
	return nil
}

func (inst *SupervisionService) GetUser(ctx context.Context, token string, email string) (*dto.SimpleUserAccountDTO, error) {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s%s/%s", url, APIGroup, SuperviseGroup, SuperviseGetEndpoint, email), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/octet-stream")
		req.Header.Set("Authorization", "Bearer "+theToken)

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	resp, err := inst.authService.PerformSafeCall(token, cb)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var response models.ErrorResult
		err := json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("user info fetch failed: %s", response.Details)
	}

	var response models.RequestResult[dto.SimpleUserAccountDTO]
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response.Value, nil
}

func (inst *SupervisionService) applyBan(ctx context.Context, token string, toBeBanned uint, restEndpoint string) (bool, error) {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		_ = writer.WriteField("userId", utils.UintToString(toBeBanned))
		writer.Close()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s%s%s/%d", url, APIGroup, SuperviseGroup, restEndpoint, toBeBanned), &body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/octet-stream")
		req.Header.Set("Authorization", "Bearer "+theToken)

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	resp, err := inst.authService.PerformSafeCall(token, cb)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var response models.ErrorResult
		err := json.Unmarshal(body, &response)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("user activity change failed: %s", response.Details)
	}

	var response models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(body, &response); err != nil {
		return false, err
	}

	return response.Result, nil
}
