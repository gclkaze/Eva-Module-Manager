package eva

type EvaProject struct {
	SchemaVersion int                      `json:"schemaVersion"`
	ModulesFolder string                   `json:"modulesFolder"`
	Modules       map[string]EvaModuleInfo `json:"modules"` // key: "module-name@version"
}

type EvaModuleInfo struct {
	ModuleName         string `json:"moduleName"`
	Version            string `json:"version"`
	InstallationFolder string `json:"installationFolder"` // relative, e.g. "eva-modules/my-module/1.1.1"
}
