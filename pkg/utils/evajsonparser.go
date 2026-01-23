package utils

import (
	"emm/internal/backend"
	"emm/internal/models/eva"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
)

var (
	ErrInvalidEvaJSON = errors.New("invalid eva.json")
)

type EvaJSONParser struct {
	CurrentSchemaVersion int
	DefaultModulesFolder string
	props                *properties.Properties
}

func NewEvaJSONParser(props *properties.Properties) *EvaJSONParser {
	currentSchemaVersion := backend.EVA_JSON_SCHEMA_CURRENT_VERSION
	defaultModulesFolder := backend.EVA_DEFAULT_EVA_MODULES

	if props != nil {
		currentSchemaVersion = props.GetInt(backend.EVA_JSON_SCHEMA_CURRENT_VERSION_KEY, backend.EVA_JSON_SCHEMA_CURRENT_VERSION)

		// NOTE: Ensure EVA_DEFAULT_EVA_MODULES_KEY exists and is the correct key for modules folder.
		defaultModulesFolder = props.GetString(backend.EVA_DEFAULT_EVA_MODULES_KEY, backend.EVA_DEFAULT_EVA_MODULES)
	}

	defaultModulesFolder = strings.TrimSpace(defaultModulesFolder)
	if defaultModulesFolder == "" {
		defaultModulesFolder = backend.EVA_DEFAULT_EVA_MODULES
	}

	return &EvaJSONParser{
		props:                props,
		CurrentSchemaVersion: currentSchemaVersion,
		DefaultModulesFolder: defaultModulesFolder,
	}
}

func (p *EvaJSONParser) ParseFile(path string) (eva.EvaProject, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return eva.EvaProject{}, fmt.Errorf("read %q: %w", path, err)
	}
	return p.ParseBytes(b)
}

func (p *EvaJSONParser) ParseBytes(b []byte) (eva.EvaProject, error) {
	if isBlank(b) {
		return p.defaultProject(), nil
	}

	prj, err := p.unmarshalProject(b)
	if err != nil {
		return eva.EvaProject{}, err
	}

	if err := p.normalizeTopLevel(&prj); err != nil {
		return eva.EvaProject{}, err
	}

	if err := p.validateSchemaVersion(prj.SchemaVersion); err != nil {
		return eva.EvaProject{}, err
	}
	if err := p.validateModulesFolder(prj.ModulesFolder); err != nil {
		return eva.EvaProject{}, err
	}

	if err := p.normalizeAndValidateModules(&prj); err != nil {
		return eva.EvaProject{}, err
	}

	return prj, nil
}

/* ----------------------------- helpers ----------------------------- */

func (p *EvaJSONParser) defaultProject() eva.EvaProject {
	return eva.EvaProject{
		SchemaVersion: p.CurrentSchemaVersion,
		ModulesFolder: p.DefaultModulesFolder,
		Modules:       map[string]eva.EvaModuleInfo{},
	}
}

func (p *EvaJSONParser) unmarshalProject(b []byte) (eva.EvaProject, error) {
	var prj eva.EvaProject
	if err := json.Unmarshal(b, &prj); err != nil {
		return eva.EvaProject{}, fmt.Errorf("%w: json parse error: %v", ErrInvalidEvaJSON, err)
	}
	return prj, nil
}

func (p *EvaJSONParser) normalizeTopLevel(prj *eva.EvaProject) error {
	if prj.SchemaVersion == 0 {
		prj.SchemaVersion = 1
	}
	if strings.TrimSpace(prj.ModulesFolder) == "" {
		prj.ModulesFolder = p.DefaultModulesFolder
	}
	prj.ModulesFolder = filepath.ToSlash(filepath.Clean(prj.ModulesFolder))

	if prj.Modules == nil {
		prj.Modules = map[string]eva.EvaModuleInfo{}
	}
	return nil
}

func (p *EvaJSONParser) validateSchemaVersion(v int) error {
	if v < 1 {
		return fmt.Errorf("%w: schemaVersion must be >= 1", ErrInvalidEvaJSON)
	}
	if v > p.CurrentSchemaVersion {
		return fmt.Errorf("%w: schemaVersion %d is newer than supported (%d)", ErrInvalidEvaJSON, v, p.CurrentSchemaVersion)
	}
	return nil
}

func (p *EvaJSONParser) validateModulesFolder(folder string) error {
	folder = strings.TrimSpace(folder)
	if folder == "." || folder == "/" || folder == `\` || folder == "" {
		return fmt.Errorf("%w: modulesFolder is invalid", ErrInvalidEvaJSON)
	}

	// Strong absolute detection across OSes:
	// - filepath.IsAbs handles native absolute paths (C:\..., \\server\share, /tmp, etc.)
	// - prefix checks catch "/something" and "\something" reliably even on Windows edge cases
	if filepath.IsAbs(folder) || strings.HasPrefix(folder, "/") || strings.HasPrefix(folder, `\`) {
		return fmt.Errorf("%w: modulesFolder must be a relative path", ErrInvalidEvaJSON)
	}

	// Path traversal guard
	if strings.Contains(folder, "..") {
		return fmt.Errorf("%w: modulesFolder must not contain '..'", ErrInvalidEvaJSON)
	}

	return nil
}

func (p *EvaJSONParser) normalizeAndValidateModules(prj *eva.EvaProject) error {
	modulesSeen := make(map[string]string)

	for rawKey, rawInfo := range prj.Modules {
		module, version, err := ParseModuleReleaseVersion(rawKey)
		if err != nil {
			return err
		}
		v, ok := modulesSeen[module]
		if ok {
			return fmt.Errorf("module %s has been redeclared with version %s@%s earlier..only one version per module is allowed to be included in the project", module, module, v)
		} else {
			modulesSeen[module] = version
		}

		key := strings.TrimSpace(rawKey)
		if key == "" {
			return fmt.Errorf("%w: modules contains an empty key", ErrInvalidEvaJSON)
		}

		parsedName, parsedVer, err := p.parseModuleKey(key)
		if err != nil {
			return err
		}

		info, err := p.normalizeAndValidateModuleInfo(prj, key, rawInfo, parsedName, parsedVer)
		if err != nil {
			return err
		}

		// Ensure map key is trimmed canonical key
		if key != rawKey {
			delete(prj.Modules, rawKey)
		}
		prj.Modules[key] = info
	}
	return nil
}

func (p *EvaJSONParser) parseModuleKey(key string) (name string, ver string, err error) {
	name, ver, err = ParseModuleReleaseVersion(key)
	if err != nil {
		return "", "", fmt.Errorf("%w: invalid module key %q (%v)", ErrInvalidEvaJSON, key, err)
	}
	name = strings.TrimSpace(name)
	ver = strings.TrimSpace(ver)

	if name == "" {
		return "", "", fmt.Errorf("%w: invalid module key %q (empty module name)", ErrInvalidEvaJSON, key)
	}

	// ver may be: "", "latest", or a concrete "vX.Y.Z" (ParseModuleReleaseVersion normalizes/validates semver)
	return name, ver, nil
}

func (p *EvaJSONParser) normalizeAndValidateModuleInfo(
	prj *eva.EvaProject,
	key string,
	info eva.EvaModuleInfo,
	keyName string,
	keyVer string,
) (eva.EvaModuleInfo, error) {
	// moduleName is always required
	if strings.TrimSpace(info.ModuleName) == "" {
		return eva.EvaModuleInfo{}, fmt.Errorf("%w: module %q missing moduleName", ErrInvalidEvaJSON, key)
	}

	info.ModuleName = strings.TrimSpace(info.ModuleName)

	// Allow floating versions:
	// - keyVer ""     (key is just "module")
	// - keyVer latest (key is "module@latest")
	// - keyVer vX.Y.Z (pinned)
	isFloating := isFloatingVersion(keyVer)

	// Normalize version field:
	// - If missing, infer from key (including "" and "latest").
	// - If present, it must match keyVer exactly.
	if strings.TrimSpace(info.Version) == "" {
		info.Version = keyVer
	} else {
		info.Version = strings.TrimSpace(info.Version)
	}

	// Enforce consistency:
	if info.ModuleName != keyName || info.Version != keyVer {
		return eva.EvaModuleInfo{}, fmt.Errorf(
			"%w: module %q inconsistent fields (key=%s@%s, fields=%s@%s)",
			ErrInvalidEvaJSON, key, keyName, keyVer, info.ModuleName, info.Version,
		)
	}

	// installationFolder:
	// - For pinned versions, force canonical "<modulesFolder>/<name>/<version>".
	// - For floating versions, allow empty; if present, validate it's safe and relative.
	if isFloating {
		if err := p.validateOptionalRelativePath(info.InstallationFolder, "installationFolder", key); err != nil {
			return eva.EvaModuleInfo{}, err
		}
	} else {
		desired := filepath.ToSlash(filepath.Join(prj.ModulesFolder, info.ModuleName, info.Version))
		info.InstallationFolder = desired
	}

	return info, nil
}

func (p *EvaJSONParser) validateOptionalRelativePath(pathValue string, fieldName string, moduleKey string) error {
	if strings.TrimSpace(pathValue) == "" {
		return nil
	}

	// Clean using OS rules first (do not ToSlash before IsAbs)
	cleanOS := filepath.Clean(strings.TrimSpace(pathValue))

	// Cross-platform absolute detection
	if filepath.IsAbs(cleanOS) || strings.HasPrefix(cleanOS, "/") || strings.HasPrefix(cleanOS, `\`) {
		return fmt.Errorf("%w: module %q %s must be relative", ErrInvalidEvaJSON, moduleKey, fieldName)
	}

	if strings.Contains(cleanOS, "..") {
		return fmt.Errorf("%w: module %q %s must not contain '..'", ErrInvalidEvaJSON, moduleKey, fieldName)
	}

	return nil
}
func isBlank(b []byte) bool {
	return len(strings.TrimSpace(string(b))) == 0
}

func isFloatingVersion(v string) bool {
	return v == "" || v == "latest"
}
