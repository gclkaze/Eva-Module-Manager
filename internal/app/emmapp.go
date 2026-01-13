package app

import (
	"emm/internal/backend"
	"emm/internal/models/userinput"
	"emm/internal/output"
	"emm/internal/services"
	"strings"
)

type EMMApp struct {
	output        output.Printer
	searchService *services.ModuleSearchService
	authService   *services.AuthService
	backend       *backend.Backend
}

func NewEMMApp(searchService *services.ModuleSearchService, authService *services.AuthService, output output.Printer) *EMMApp {
	backend := backend.NewBackend()
	return &EMMApp{searchService: searchService, authService: authService, output: output, backend: backend}
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

	inst.authService.SetPrinter(inst.output)
	inst.searchService.SetPrinter(inst.output)

	return nil
}

func (inst *EMMApp) InitFromPath(path string) error {
	return inst.backend.InitFromPath(path)
}

func (inst EMMApp) SearchByComponents(name []string, description []string, tags []string) error {
	err := inst.searchService.SearchByComponents(name, description, tags)
	return err
}

func (inst EMMApp) UserRegister(creds *userinput.RegistrationCreds) error {
	_, err := inst.authService.Register(creds)
	return err
}

func (inst EMMApp) UserLogin(creds *userinput.LoginCreds) error {
	_, err := inst.authService.Login(creds)
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
