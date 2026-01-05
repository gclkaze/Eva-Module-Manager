package app

import (
	"emm/internal/backend"
	"emm/internal/output"
	"emm/internal/services"
	"strings"
)

type EMMApp struct {
	output        output.Printer
	searchService *services.ModuleSearchService
	backend       *backend.Backend
}

func NewEMMApp(searchService *services.ModuleSearchService, output output.Printer) *EMMApp {
	backend := backend.NewBackend()
	return &EMMApp{searchService: searchService, output: output, backend: backend}
}

func (inst *EMMApp) Init() error {
	err := inst.backend.Init()
	if err != nil {
		return err
	}

	inst.searchService.SetBackend(inst.backend)
	return nil
}

func (inst *EMMApp) InitFromPath(path string) error {
	return inst.backend.InitFromPath(path)
}

func (inst EMMApp) SearchByComponents(name []string, description []string, tags []string) error {
	err := inst.searchService.SearchByComponents(name, description, tags)
	return err
}

func (inst EMMApp) SearchByQuery(query string) error {
	l := strings.Split(query, ",")
	return inst.SearchByComponents(l, l, l)
}
