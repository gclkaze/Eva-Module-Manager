package app

import (
	"emm/internal/output"
	"emm/internal/services"
	"fmt"
	"io"
	"net/http"

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
	//http.MethodGet, fmt.Sprintf("%s%s%s", services.APIGroup, services.ModulesGroup, services.ModuleSearchEndpoint)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080%s%s%s", services.APIGroup, services.ModulesGroup, services.ModuleSearchEndpoint), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	//req.Header.Set("Authorization", "Bearer YOUR_TOKEN")
	q := req.URL.Query()
	for _, n := range name {
		q.Add("name", n)
	}
	for _, d := range description {
		q.Add("description", d)
	}
	for _, tag := range tags {
		q.Add("tags", tag)
	}

	// IMPORTANT: assign back
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(body))
	return err
}

func (inst EMMApp) SearchByQuery(query string) error {
	return nil
}
