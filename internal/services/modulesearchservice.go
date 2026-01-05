package services

import (
	"emm/internal/backend"
	"fmt"
	"io"
	"net/http"
)

type ModuleSearchService struct {
	backend *backend.Backend
}

func NewModuleSearchService() *ModuleSearchService {
	return &ModuleSearchService{}
}

func (inst *ModuleSearchService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
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
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(body))
	return err
}
