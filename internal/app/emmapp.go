package app

import (
	"context"
	"emm/internal/backend"
	"emm/internal/models/userinput"
	"emm/internal/output"
	"emm/internal/services"
	"fmt"
	"strings"
)

type EMMApp struct {
	cwd           string
	output        output.Printer
	searchService *services.ModuleSearchService
	authService   *services.AuthService
	backend       *backend.Backend

	moduleService  *services.ModuleService
	releaseService *services.ModuleReleaseService

	bookKeepingService *services.ProjectBookkeepingService
	saveLocation       string
	onError            bool
}

func NewEMMApp(cwd string, searchService *services.ModuleSearchService, authService *services.AuthService, moduleService *services.ModuleService, releaseService *services.ModuleReleaseService, bookKeepingService *services.ProjectBookkeepingService, output output.Printer) *EMMApp {
	backend := backend.NewBackend()
	return &EMMApp{cwd: cwd, searchService: searchService, authService: authService, output: output, backend: backend, moduleService: moduleService, releaseService: releaseService, bookKeepingService: bookKeepingService, onError: false}
}

func (inst EMMApp) GetCurrentWorkingDirector() string {
	return inst.cwd
}

func (inst EMMApp) GetDefaultFileStorageLocation() string {
	return inst.saveLocation
}

func (inst EMMApp) VerifyEvaProjectFile(p string) (string, error) {
	return inst.bookKeepingService.VerifyExisting(p)
}

func (inst EMMApp) IsOnError() bool {
	return inst.onError
}

func (inst *EMMApp) SetOnError() {
	inst.onError = true
}

func (inst EMMApp) GetPrinter() output.Printer {
	return inst.output
}
func (inst *EMMApp) Init() error {
	err := inst.backend.Init()
	if err != nil {
		return err
	}

	inst.authService.SetBackend(inst.backend)
	inst.searchService.SetBackend(inst.backend)
	inst.moduleService.SetBackend(inst.backend)
	inst.releaseService.SetBackend(inst.backend)
	inst.bookKeepingService.SetBackend(inst.backend)

	inst.authService.SetPrinter(inst.output)
	inst.searchService.SetPrinter(inst.output)
	inst.moduleService.SetPrinter(inst.output)
	inst.releaseService.SetPrinter(inst.output)
	inst.bookKeepingService.SetPrinter(inst.output)

	inst.saveLocation = inst.backend.GetDefaultFileStorageLocation()

	return nil
}

func (inst *EMMApp) InitFromPath(path string) error {
	return inst.backend.InitFromPath(path)
}

func (inst EMMApp) SearchByComponents(name []string, description []string, tags []string) error {
	err := inst.searchService.SearchByComponents(name, description, tags)
	return err
}

func (inst EMMApp) UploadModule(token string, paths []string,
	params *userinput.UploadParams) error {
	return inst.moduleService.UploadModule(token, paths, params)
}

func (inst EMMApp) SuggestModuleRelease(token string,
	params *userinput.ModuleReleaseSuggestionParams) error {
	return inst.moduleService.SuggestModuleRelease(token, params)
}

func (inst EMMApp) UpdateModule(token string, paths []string,
	params *userinput.UploadModuleUpdateParams) error {
	return inst.moduleService.UpdateModule(token, paths, params)
}

func (inst EMMApp) GetUserModules(token string) error {
	return inst.moduleService.GetUserModules(token)
}

func (inst EMMApp) UserRegister(creds *userinput.RegistrationCreds) error {
	_, err := inst.authService.Register(creds)
	return err
}

func (inst EMMApp) IsCurrentUserAuthorized() error {
	_, err := inst.authService.GetActiveUser()
	if err != nil {
		return fmt.Errorf("current user needs to be logged on first")
	}
	return nil
}

func (inst EMMApp) GetCurrentUserToken() (string, error) {
	return inst.authService.GetCurrentUserToken()
}

func (inst *EMMApp) SwitchCurrentUser(email string) error {
	return inst.authService.SwitchCurrentActiveUser(email)
}

func (inst EMMApp) ShowCurrentUser() error {
	return inst.authService.ShowCurrentUser()
}

func (inst EMMApp) UserLogin(creds *userinput.LoginCreds) error {
	_, err := inst.authService.Login(creds)
	return err
}

func (inst EMMApp) UserLogout(email string) error {
	err := inst.authService.Logout(email)
	return err
}

func (inst EMMApp) SearchBySearchQuery(q *userinput.ModuleSearchQuery) error {
	err := inst.searchService.SearchByComponents(q.Name, q.Description, q.Tags)
	return err
}

func (inst EMMApp) SearchByQuery(query string) error {
	l := strings.Split(query, ",")
	return inst.SearchByComponents(l, l, l)
}

func (inst EMMApp) GetModuleInfo(query string) error {
	return inst.searchService.GetModuleInfo(query)
}
func (inst EMMApp) AcceptRelease(token string, module string, version string) error {
	return inst.releaseService.AcceptRelease(token, module, version)
}

func (inst EMMApp) CancelRelease(token string, module string, version string) error {
	return inst.releaseService.CancelRelease(token, module, version)
}

func (inst EMMApp) LowerRelease(token string, module string, version string) error {
	return inst.releaseService.LowerRelease(token, module, version)
}

func (inst EMMApp) RejectRelease(token string, module string, version string) error {
	return inst.releaseService.RejectRelease(token, module, version)
}

func (inst EMMApp) ReleaseDump(token string, filter *userinput.ReleaseFilterParams, view string) error {
	return inst.releaseService.ReleaseDump(token, filter, view)
}

func (inst EMMApp) DownloadRelease(ctx context.Context, token string, module string, version string, saveLocation string) error {
	return inst.releaseService.DownloadRelease(ctx, token, module, version, saveLocation)
}
