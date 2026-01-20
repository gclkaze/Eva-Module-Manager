package backend

import (
	"emm/internal/config"
	"fmt"

	"github.com/magiconair/properties"
)

type Backend struct {
	properties *properties.Properties
}

func NewBackend() *Backend {
	return &Backend{}
}

func (b *Backend) Init() error {
	return b.readCondiguration()
}

func (b *Backend) InitFromPath(path string) error {
	return b.readCondigurationFromPath(path)
}

func (b Backend) GetDefaultFileStorageLocation() string {
	return b.properties.GetString("save_folder_location", ".")
}

func (b *Backend) readCondiguration() error {
	config.Init()

	if config.TheConfigReader.IsOnError() {
		return config.TheConfigReader.GetError()
	}

	b.properties = config.TheConfigReader.GetProperties()
	return nil
}

func (b *Backend) readCondigurationFromPath(path string) error {
	config.InitWithPropertiesPath(path)

	if config.TheConfigReader.IsOnError() {
		return config.TheConfigReader.GetError()
	}

	b.properties = config.TheConfigReader.GetProperties()
	return nil
}
func (b Backend) GetServerURL() (string, error) {
	proto := b.properties.GetString("protocol", "")
	server := b.properties.GetString("server", "")
	port := b.properties.GetString("port", "")

	if proto == "" {
		return "", fmt.Errorf("no protocol was specified for the Module Server")
	}
	if server == "" {
		return "", fmt.Errorf("no server base URL was specified for the Module Server")
	}

	if proto == "https" {
		return fmt.Sprintf("%s://%s", proto, server), nil
	}

	if port == "" {
		return "", fmt.Errorf("no server base port was specified for the Module Server")
	}
	return fmt.Sprintf("%s://%s:%s", proto, server, port), nil
}
