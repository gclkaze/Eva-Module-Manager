package services

import (
	"bytes"
	"context"
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/models/userinput"
	"emm/internal/output"
	"emm/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/magiconair/properties"
	"github.com/schollz/progressbar/v3"
)

type ModuleReleaseService struct {
	backend     *backend.Backend
	output      output.Printer
	authService *AuthService
	props       *properties.Properties
}

func NewModuleReleaseService(auth *AuthService) *ModuleReleaseService {
	return &ModuleReleaseService{authService: auth}
}

func (inst *ModuleReleaseService) SetProperties(props *properties.Properties) {
	inst.props = props
}

func (inst *ModuleReleaseService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *ModuleReleaseService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst *ModuleReleaseService) AcceptRelease(token string, module string, version string) error {
	rel, err := inst.FindRelease(token, module, version)
	if err != nil {
		inst.output.Error(err)
		return nil
	}

	if rel == nil {
		inst.output.Error(fmt.Errorf("couldn't find module release %s@%s", module, version))
		return nil
	}

	err = inst.ApplyRelease(token, rel.ID, SuperviseAcceptReleaseEndpoint)
	if err == nil {
		inst.output.Info(fmt.Sprintf("Release %s@%s was accepted successfully", module, version))
	} else {
		inst.output.Info(fmt.Sprintf("Release %s@%s couldn't be accepted successfully", module, version))
	}
	return nil
}

func (inst *ModuleReleaseService) CancelRelease(token string, module string, version string) error {
	rel, err := inst.FindRelease(token, module, version)
	if err != nil {
		inst.output.Error(err)
		return nil
	}

	if rel == nil {
		inst.output.Error(fmt.Errorf("couldn't find module release %s@%s", module, version))
		return nil
	}

	err = inst.ApplyRelease(token, rel.ID, SuperviseCancelReleaseEndpoint)
	if err == nil {
		inst.output.Info(fmt.Sprintf("Release %s@%s was cancelled successfully", module, version))
	} else {
		inst.output.Info(fmt.Sprintf("Release %s@%s couldn't be cancelled successfully", module, version))
	}
	return nil
}

func (inst *ModuleReleaseService) LowerRelease(token string, module string, version string) error {
	rel, err := inst.FindRelease(token, module, version)
	if err != nil {
		inst.output.Error(err)
		return nil
	}

	if rel == nil {
		inst.output.Error(fmt.Errorf("couldn't find module release %s@%s", module, version))
		return nil
	}

	err = inst.ApplyRelease(token, rel.ID, SupervisePendingReleaseEndpoint)
	if err == nil {
		inst.output.Info(fmt.Sprintf("Release %s@%s being set to pending successfully", module, version))
	} else {
		inst.output.Info(fmt.Sprintf("Release %s@%s couldn't being set to pending successfully", module, version))
	}

	return nil
}

func (inst *ModuleReleaseService) RejectRelease(token string, module string, version string) error {
	rel, err := inst.FindRelease(token, module, version)
	if err != nil {
		inst.output.Error(err)
		return nil
	}

	if rel == nil {
		inst.output.Error(fmt.Errorf("couldn't find module release %s@%s", module, version))
		return nil
	}

	err = inst.ApplyRelease(token, rel.ID, SuperviseRejectReleaseEndpoint)
	if err == nil {
		inst.output.Info(fmt.Sprintf("Release %s@%s was rejected successfully", module, version))
	} else {
		inst.output.Info(fmt.Sprintf("Release %s@%s couldn't being reject successfully", module, version))
	}

	return nil
}

func (inst ModuleReleaseService) DownloadRelease(ctx context.Context, token string, module string, version string, saveLocation string) error {
	var err error
	if token == "" {
		err = inst.downloadPublicRelease(ctx, module, version, saveLocation)
	} else {
		err = inst.downloadRelease(ctx, token, module, version, saveLocation)
	}
	if err != nil {
		inst.output.Error(err)
	} else {
		inst.output.Info(fmt.Sprintf("✅ Module %s was downloaded successfully at '%s'.", inst.GetModuleTarBallName(module, version), saveLocation))
	}
	return err
}

func (inst ModuleReleaseService) GetLatestModuleRelease(ctx context.Context, token string, module string) (*models.ReleaseDTO, error) {
	url, err := inst.backend.GetServerURL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s%s%s", url, APIGroup, ReleasesGroup, "/latest/", module), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/octet-stream")

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var response models.RequestResult[models.ReleaseDTO]
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}
		return &response.Value, nil
	}

	var errResp models.ErrorResult
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%s", errResp.Details)

}

func (inst ModuleReleaseService) downloadPublicRelease(ctx context.Context, module string, version string, saveLocation string) error {
	url, err := inst.backend.GetServerURL()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s%s%s@%s", url, APIGroup, DownloadGroup, "/release/", module, version), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/octet-stream")

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	saveAs := inst.GetModuleTarBallName(module, version)
	return inst.StoreRelease(resp, saveLocation, saveAs)
}

func (inst ModuleReleaseService) GetModuleTarBallName(module string, version string) string {
	return fmt.Sprintf("%s-%s.tar.gz", module, version)
}

/*
download.GET("/release/:release", router.downloadHandler.DownloadPublicRelease)
download.GET("/auth/release/:release", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.downloadHandler.AuthUserDownloadSpecificRelease)
*/
func (inst ModuleReleaseService) downloadRelease(ctx context.Context, token string, module string, version string, saveLocation string) error {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s%s%s%s@%s", url, APIGroup, DownloadGroup, "/auth", "/release/", module, version), nil)
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

	saveAs := inst.GetModuleTarBallName(module, version)
	return inst.StoreRelease(resp, saveLocation, saveAs)
}

func (inst ModuleReleaseService) StoreRelease(resp *http.Response, destPath string, saveAs string) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var response models.ErrorResult
		b, _ := io.ReadAll(resp.Body)
		err := json.Unmarshal(b, &response)
		if err != nil {
			return err
		}
		return fmt.Errorf("download failed: %s", response.Details)
	}

	if !utils.FolderExists(filepath.Dir(destPath)) {
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}
	}

	theFile := filepath.Join(destPath, saveAs)
	out, err := os.Create(theFile)
	if err != nil {
		return err
	}
	defer out.Close()

	total := resp.ContentLength // -1 if unknown

	var bar *progressbar.ProgressBar

	if total > 0 {
		bar = progressbar.DefaultBytes(resp.ContentLength, "downloading")
	} else {
		bar = progressbar.NewOptions(
			-1,
			progressbar.OptionSetDescription("downloading"),
			progressbar.OptionShowBytes(true),
		)
	}

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return err
	}

	_ = bar.Finish()
	return nil
}

func (inst ModuleReleaseService) ReleaseDump(token string, filter *userinput.ReleaseFilterParams, view string) error {
	modules, err := inst.GetFilteredModuleReleases(token, filter)
	if err != nil {
		inst.output.Error(err)
		return nil
	}

	if view == "rows" {
		inst.output.PrintReleaseRows(modules)
	} else {
		inst.output.PrintDetailedModuleReleaseInfo(modules)

	}
	return nil
}

func (inst ModuleReleaseService) GetFilteredModuleReleases(token string, p *userinput.ReleaseFilterParams) ([]models.ModuleEnrichedDTO, error) {

	cb := func(theToken string) (*http.Response, error) {
		baseURL, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		// Build URL with query params
		u, err := url.Parse(fmt.Sprintf(
			"%s%s%s%s",
			baseURL,
			APIGroup,
			SuperviseGroup,
			SuperviseGetFilterReleaseEndpoint,
		))
		if err != nil {
			return nil, err
		}

		q := u.Query()

		// []string filters (repeatable)
		addSlice := func(key string, vals []string) {
			for _, v := range vals {
				q.Add(key, v)
			}
		}

		if p != nil {
			addSlice("status", p.Status)
			addSlice("versions", p.Versions)
			addSlice("tags", p.Tags)
			addSlice("module", p.ModuleName)
			addSlice("repo", p.RepoName)
			addSlice("description", p.Description)
			addSlice("creator", p.Creator)
			addSlice("creator-email", p.CreatorEmail)

			if p.CreatedAfter != nil {
				q.Set("created-after", p.CreatedAfter.Format(time.RFC3339))
			}
			if !p.ReleasedAfter.IsZero() {
				q.Set("released-after", p.ReleasedAfter.Format(time.RFC3339))
			}
		}

		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+theToken)

		client := &http.Client{}
		return client.Do(req)
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

	if resp.StatusCode == http.StatusOK {
		var response models.RequestResult[[]models.ModuleEnrichedDTO]
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}
		return response.Value, nil
	}

	var errResp models.ErrorResult
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%s", errResp.Details)
}

func (inst ModuleReleaseService) FindRelease(token string, module string, version string) (*models.Release, error) {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s%s/%s/%s", url, APIGroup, SuperviseGroup, SuperviseFindReleaseEndpoint, module, version), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+theToken)

		client := &http.Client{}
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

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var response models.RequestResult[models.Release]
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return &response.Value, nil
	}

	var response models.ErrorResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%s", response.Details)
}

func (inst ModuleReleaseService) ApplyRelease(token string, releaseID uint, verb string) error {

	cb := func(theToken string) (*http.Response, error) {

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		_ = writer.WriteField("releaseId", utils.UintToString(releaseID))
		writer.Close()

		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s/%d", url, APIGroup, SuperviseGroup, verb, releaseID), &body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+theToken)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := http.DefaultClient.Do(req)
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

	if resp.StatusCode != http.StatusOK {
		var response models.ErrorResult
		b, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(b, &response)
		if err != nil {
			inst.output.Error(err)
			return err
		}
		err = fmt.Errorf("%s", response.Details)
		inst.output.Error(err)
		return err
	}
	return err
}
