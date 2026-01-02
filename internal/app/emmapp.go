package app

import (
	"emm/internal/output"
	"emm/internal/services"

	"github.com/magiconair/properties"
)

type EMMApp struct {
	properties    *properties.Properties
	output        output.Printer
	searchService *services.ModuleSearchService
}

func NewEMMApp(searchService *services.ModuleSearchService, output output.Printer) *EMMApp {
	return &EMMApp{searchService: searchService, output: output}
}

func (inst EMMApp) SearchByComponents(name []string, description []string, tags []string) error {
	return nil
}

func (inst EMMApp) SearchByQuery(query string) error {
	return nil
}
