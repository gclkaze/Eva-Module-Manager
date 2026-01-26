package services

import (
	"bytes"
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/models/userinput"
	"emm/internal/output"
	"emm/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/magiconair/properties"
)

type ModuleService struct {
	backend       *backend.Backend
	output        output.Printer
	props         *properties.Properties
	authService   *AuthService
	maxUploadSize int64
}

func NewModuleService(auth *AuthService) *ModuleService {
	return &ModuleService{authService: auth}
}

func (inst *ModuleService) SetProperties(props *properties.Properties) {
	inst.props = props

	if props != nil {
		inst.maxUploadSize = props.GetInt64("maxUploadSize", 100)
	}
}

func (inst *ModuleService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *ModuleService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst *ModuleService) UploadModule(token string, paths []string, params *userinput.UploadParams) error {
	return inst.validateAndUploadAll(token, paths, params)
}

func (inst *ModuleService) UpdateModule(token string, paths []string, params *userinput.UploadModuleUpdateParams) error {
	return inst.validateAndUpdateModule(token, paths, params)
}

func (inst ModuleService) findUserModule(token string, moduleName string) (*models.Module, error) {
	list, err := inst.GetUserModulesList(token)
	if err != nil {
		return nil, err
	}

	for i := range list {
		if list[i].RepoName == moduleName {
			return &list[i], nil
		}
	}
	err = fmt.Errorf("couldn't find user module: %s", moduleName)
	return nil, err
}

func (inst *ModuleService) SuggestModuleRelease(token string, params *userinput.ModuleReleaseSuggestionParams) error {
	m, err := inst.findUserModule(token, params.ModuleRepr)
	if err != nil {
		inst.output.Error(err)
		return nil
	}

	if m == nil {
		inst.output.Error(fmt.Errorf("coulnd't find user module: '%s'", params.ModuleRepr))
		return nil
	}
	cb := func(theToken string) (*http.Response, error) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		_ = writer.WriteField("modId", utils.UintToString(m.ID))
		_ = writer.WriteField("version", params.Version)
		writer.Close()

		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, ModuleSuggestEndpoint), &body)
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
			return nil
		}
		inst.output.Error(fmt.Errorf("%s", response.Details))
		return nil
	}

	inst.output.Info(fmt.Sprintf("Module release %s@%s was suggested successfully!", params.ModuleRepr, params.Version))
	return nil
}

func (inst ModuleService) GetUserModulesList(token string) ([]models.Module, error) {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, GetUserModulesEndpoint), nil)
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
		var response models.RequestResult[[]models.Module]
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return response.Value, nil
	}
	var response models.ErrorResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%s", response.Details)

}

func (inst ModuleService) GetUserModules(token string) error {
	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, GetUserModulesEndpoint), nil)
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
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var response models.RequestResult[[]models.Module]
		err = json.Unmarshal(body, &response)
		if err != nil {
			inst.output.Error(err)
			return err
		}
		inst.output.PrintModules(response.Value)
	} else {
		var response models.ErrorResult
		err = json.Unmarshal(body, &response)
		if err != nil {
			inst.output.Error(err)
			return err
		}
		inst.output.Error(fmt.Errorf("%s", response.Details))
		err = nil
	}

	return err
}

func (inst *ModuleService) validateAndUpdateModule(token string, paths []string, params *userinput.UploadModuleUpdateParams) error {
	m, err := inst.findUserModule(token, params.ModuleRepr)
	if err != nil {
		inst.output.Error(err)
		return nil
	}
	var total int64

	// Step 1: expand folders into files
	files, err := inst.expandPaths(paths)
	if err != nil {
		return err
	}

	oldSize := len(files)
	files = utils.UniqueStrings(files)

	if oldSize != len(files) {
		inst.output.Warn("Ommitting duplicate file upload..")
	}

	// Step 2: sum total size
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return err
		}
		total += info.Size()
	}

	// Step 3: add safety margin for multipart overhead
	const overhead = 256 * 1024 // 256 KB
	if total+overhead > inst.maxUploadSize {
		return fmt.Errorf(
			"total upload size %.2f MB exceeds backend limit %.2f MB",
			float64(total)/1024/1024,
			float64(inst.maxUploadSize)/1024/1024,
		)
	}

	cb := func(theToken string) (*http.Response, error) {

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		_ = writer.WriteField("title", params.Title)
		_ = writer.WriteField("repr", params.Repr)
		_ = writer.WriteField("tags", params.Tags)
		_ = writer.WriteField("description", params.Description)
		_ = writer.WriteField("modId", utils.UintToString(m.ID))

		for _, f := range files {
			if err := inst.addFile(writer, "file", f); err != nil {
				return nil, err
			}
		}

		writer.Close()

		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, ModuleUpdateEndpoint), &body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+theToken)
		// Step 5: POST once
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
		inst.output.Error(fmt.Errorf("%s", response.Details))
		return fmt.Errorf("upload failed: %s", string(b))
	}

	inst.output.Info(fmt.Sprintf("✅ Uploaded %d files successfully", len(files)))
	inst.output.Info(fmt.Sprintf("Module %s was updated successfully!", params.ModuleRepr))
	return nil
}

func (inst *ModuleService) validateAndUploadAll(token string, paths []string, params *userinput.UploadParams) error {
	var total int64

	// Step 1: expand folders into files
	files, err := inst.expandPaths(paths)
	if err != nil {
		return err
	}

	oldSize := len(files)
	files = utils.UniqueStrings(files)

	if oldSize != len(files) {
		inst.output.Warn("Ommitting duplicate file upload..")
	}

	// Step 2: sum total size
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return err
		}
		total += info.Size()
	}

	// Step 3: add safety margin for multipart overhead
	const overhead = 256 * 1024 // 256 KB
	if total+overhead > inst.maxUploadSize {
		return fmt.Errorf(
			"total upload size %.2f MB exceeds backend limit %.2f MB",
			float64(total)/1024/1024,
			float64(inst.maxUploadSize)/1024/1024,
		)
	}

	// Step 4: build multipart request with all files
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("title", params.Title)
	_ = writer.WriteField("repr", params.Repr)
	_ = writer.WriteField("tags", params.Tags)
	_ = writer.WriteField("description", params.Description)

	for _, f := range files {
		if err := inst.addFile(writer, "file", f); err != nil {
			return err
		}
	}

	writer.Close()

	cb := func(theToken string) (*http.Response, error) {
		url, err := inst.backend.GetServerURL()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, ModuleUploadEndpoint), &body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+theToken)
		// Step 5: POST once
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
		//inst.output.Error(fmt.Errorf("%s", response.Details))
		return fmt.Errorf("upload failed: %s", response.Details)
	}

	inst.output.Info(fmt.Sprintf("✅ Uploaded %d files successfully", len(files)))
	return nil
}

func (inst *ModuleService) addFile(
	writer *multipart.Writer,
	fieldName string,
	filePath string,
) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	return err
}

func (inst *ModuleService) expandPaths(paths []string) ([]string, error) {
	var files []string

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("invalid path %s: %w", p, err)
		}

		if info.IsDir() {
			err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			files = append(files, p)
		}
	}

	return files, nil
}
