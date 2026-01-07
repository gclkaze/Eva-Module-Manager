package services

import (
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/output"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ModuleSearchService struct {
	backend *backend.Backend
	output  output.Printer
}

func NewModuleSearchService() *ModuleSearchService {
	return &ModuleSearchService{}
}

func (inst *ModuleSearchService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *ModuleSearchService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst ModuleSearchService) SearchByComponents(name []string, description []string, tags []string) error {
	url, err := inst.backend.GetServerURL()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, ModuleSearchEndpoint), nil)
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

	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
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
		err = fmt.Errorf("couldn't fetch the module information")
		inst.output.Error(err)
	}
	/*	fmt.Println(resp.Status)
		fmt.Println(string(body))*/

	return err
}

func (inst ModuleSearchService) GetModuleInfo(moduleName string) error {
	url, err := inst.backend.GetServerURL()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s%s", url, APIGroup, ModulesGroup, ModuleGetInfoEndpoint), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	//req.Header.Set("Authorization", "Bearer YOUR_TOKEN")
	q := req.URL.Query()
	q.Add("moduleName", moduleName)

	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var response models.RequestResult[models.ModuleEnrichedInformation]
		err = json.Unmarshal(body, &response)
		if err != nil {
			inst.output.Error(err)
			return err
		}
		inst.output.PrintModuleInfo(response.Value)
	} else {
		//err = fmt.Errorf("couldn't fetch the module information")
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
