package services

import (
	"emm/internal/backend"
	"emm/internal/models/eva"
	"emm/internal/output"
	"emm/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/magiconair/properties"
)

/*const (
	DefaultEvaFileName   = "eva.json"
	CurrentSchemaVersion = 1
	DefaultModulesFolder = "eva-modules"
)*/

type ProjectBookkeepingService struct {
	backend *backend.Backend
	output  output.Printer

	projectRoot string
	evaFilePath string

	props  *properties.Properties
	parser *utils.EvaJSONParser

	mu      sync.RWMutex
	project eva.EvaProject

	defaultEvaFileName string

	currentSchemaVersion int
	defaultModulesFolder string
}

func NewProjectBookkeepingService(projectRoot string) (*ProjectBookkeepingService, error) {
	if projectRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getwd: %w", err)
		}
		projectRoot = cwd
	}

	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("abs(projectRoot): %w", err)
	}

	s := &ProjectBookkeepingService{
		projectRoot: absRoot,
		//	evaFilePath: filepath.Join(absRoot, DefaultEvaFileName),
	}

	return s, nil
}

func (inst *ProjectBookkeepingService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *ProjectBookkeepingService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst *ProjectBookkeepingService) SetProperties(props *properties.Properties) {
	inst.props = props
	inst.parser = utils.NewEvaJSONParser(props)

	if props != nil {
		inst.defaultEvaFileName = props.GetString(backend.EVA_DEFAULT_PROJECT_FILENAME_KEY, "")
		inst.evaFilePath = filepath.Join(inst.projectRoot, inst.defaultEvaFileName)

		inst.currentSchemaVersion = props.GetInt(backend.EVA_JSON_SCHEMA_CURRENT_VERSION_KEY, backend.EVA_JSON_SCHEMA_CURRENT_VERSION)
		inst.defaultModulesFolder = props.GetString(backend.EVA_DEFAULT_EVA_MODULES_KEY, backend.EVA_DEFAULT_EVA_MODULES)
	} else {

		inst.defaultEvaFileName = backend.EVA_DEFAULT_PROJECT_FILENAME
		inst.evaFilePath = filepath.Join(inst.projectRoot, inst.defaultEvaFileName)
		inst.currentSchemaVersion = backend.EVA_JSON_SCHEMA_CURRENT_VERSION
		inst.defaultModulesFolder = backend.EVA_DEFAULT_EVA_MODULES

	}
}

func (s *ProjectBookkeepingService) ProjectRoot() string { return s.projectRoot }
func (s *ProjectBookkeepingService) EvaFilePath() string { return s.evaFilePath }

// --------------------------- Verify / Resolve ---------------------------

// VerifyExisting resolves eva.json location and parses it.
// - If inputPath is empty: checks ./eva.json
// - If inputPath is a dir: checks <dir>/eva.json
// - If inputPath is a file: must be eva.json
// Returns the absolute eva.json path if OK.
func (s *ProjectBookkeepingService) VerifyExisting(inputPath string) (string, error) {
	absEvaPath, err := s.ResolveEvaJSONPath(inputPath, true)
	if err != nil {
		return "", err
	}

	// Validate via parser (supports floating versions).
	if _, err := s.parser.ParseFile(absEvaPath); err != nil {
		return "", err
	}

	return absEvaPath, nil
}

// ResolveEvaJSONPath resolves the target eva.json path.
// inputPath can be "", a directory, or an eva.json file path.
// If requireExists is true, it errors if the target eva.json does not exist.
func (s *ProjectBookkeepingService) ResolveEvaJSONPath(inputPath string, requireExists bool) (string, error) {
	in := strings.TrimSpace(inputPath)

	// Default: ./eva.json
	if in == "" {
		target := filepath.Join(".", s.defaultEvaFileName)
		return s.ensureExistsOrReturnAbs(target, requireExists, "eva.json not found in current directory")
	}

	abs, err := filepath.Abs(in)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", in, err)
	}

	fi, err := os.Stat(abs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("path not found: %s", abs)
		}
		return "", fmt.Errorf("cannot access %s: %w", abs, err)
	}

	// Directory -> <dir>/eva.json
	if fi.IsDir() {
		target := filepath.Join(abs, s.defaultEvaFileName)
		return s.ensureExistsOrReturnAbs(target, requireExists, "eva.json not found in directory")
	}

	// File -> must be eva.json
	if !strings.EqualFold(filepath.Base(abs), s.defaultEvaFileName) {
		return "", fmt.Errorf("expected %s file, got: %s", s.defaultEvaFileName, abs)
	}

	if requireExists {
		if _, err := os.Stat(abs); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return "", fmt.Errorf("%s not found: %s", s.defaultEvaFileName, abs)
			}
			return "", fmt.Errorf("cannot access %s: %w", abs, err)
		}
	}

	return abs, nil
}

func (s *ProjectBookkeepingService) ensureExistsOrReturnAbs(path string, requireExists bool, notFoundMsg string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", path, err)
	}

	if !requireExists {
		return abs, nil
	}

	if _, err := os.Stat(abs); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%s (%s)", notFoundMsg, abs)
		}
		return "", fmt.Errorf("cannot access %s: %w", abs, err)
	}

	return abs, nil
}

// --------------------------- Load / Save ---------------------------

// LoadExisting loads and parses an existing eva.json into memory.
// It does NOT create the file.
func (s *ProjectBookkeepingService) LoadExisting(inputPath string) error {
	absEvaPath, err := s.ResolveEvaJSONPath(inputPath, true)
	if err != nil {
		return err
	}

	prj, err := s.parser.ParseFile(absEvaPath)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.projectRoot = filepath.Dir(absEvaPath)
	s.evaFilePath = absEvaPath
	s.project = prj
	s.mu.Unlock()

	return nil
}

// LoadOrCreate ensures eva.json exists at the resolved location (default ./eva.json), then loads it.
// This is appropriate for commands like `emm add` that can initialize a project file.
func (s *ProjectBookkeepingService) LoadOrCreate(inputPath string) error {
	absEvaPath, err := s.ResolveEvaJSONPath(inputPath, false)
	if err != nil {
		return err
	}

	// If inputPath was a directory, ResolveEvaJSONPath(requireExists=false) returns "<dir>/eva.json" only
	// when you passed a directory. But for "", it returns "./eva.json". We already have absEvaPath.
	if err := s.ensureEvaFileExists(absEvaPath); err != nil {
		return err
	}

	prj, err := s.parser.ParseFile(absEvaPath)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.projectRoot = filepath.Dir(absEvaPath)
	s.evaFilePath = absEvaPath
	s.project = prj
	s.mu.Unlock()

	return nil
}

func (s *ProjectBookkeepingService) Save() error {
	s.mu.RLock()
	p := cloneProject(s.project)
	target := s.evaFilePath
	s.mu.RUnlock()

	tmp := target + ".tmp"
	if err := writeEvaProjectFile(tmp, p); err != nil {
		return err
	}
	if err := os.Rename(tmp, target); err != nil {
		return fmt.Errorf("rename tmp -> eva.json: %w", err)
	}
	return nil
}

func (s *ProjectBookkeepingService) Get() eva.EvaProject {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneProject(s.project)
}

// --------------------------- Mutation helpers ---------------------------

// AddModule adds a module entry; version may be "", "latest", or pinned.
// Keying rules:
// - version ""      => key is "<name>"
// - version "latest"=> key is "<name>@latest"
// - pinned          => key is "<name>@<version>" (expecting canonical, e.g. v1.2.3)
func (s *ProjectBookkeepingService) AddModule(moduleName, version string) error {
	moduleName = strings.TrimSpace(moduleName)
	version = strings.TrimSpace(version)

	if moduleName == "" {
		return errors.New("moduleName is required")
	}

	key := moduleName
	if version != "" {
		key = fmt.Sprintf("%s@%s", moduleName, version)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.project.Modules == nil {
		s.project.Modules = map[string]eva.EvaModuleInfo{}
	}
	s.project.Modules[key] = eva.EvaModuleInfo{
		ModuleName: moduleName,
		Version:    version,
		// For pinned versions, the parser (and/or install flow) will canonicalize installationFolder.
		// For floating, keep it empty.
	}
	return nil
}

func (s *ProjectBookkeepingService) PurgeKey(moduleKey string) error {
	moduleKey = strings.TrimSpace(moduleKey)
	if moduleKey == "" {
		return errors.New("moduleKey is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.project.Modules == nil {
		return nil
	}
	delete(s.project.Modules, moduleKey)
	return nil
}

// ResolveAndPin rewrites a floating entry to a pinned one.
// Example:
// - "foo" or "foo@latest"  -> "foo@v1.2.3"
// It removes any existing entries for that module name (floating or pinned) and inserts the pinned entry.
func (s *ProjectBookkeepingService) ResolveAndPin(moduleName string, resolvedVersion string) error {
	moduleName = strings.TrimSpace(moduleName)
	resolvedVersion = strings.TrimSpace(resolvedVersion)

	if moduleName == "" || resolvedVersion == "" || resolvedVersion == "latest" {
		return errors.New("moduleName and a concrete resolvedVersion are required")
	}

	newKey := fmt.Sprintf("%s@%s", moduleName, resolvedVersion)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.project.Modules == nil {
		s.project.Modules = map[string]eva.EvaModuleInfo{}
	}

	// Remove all keys for this module (foo, foo@latest, foo@vX)
	prefix := moduleName + "@"
	for k := range s.project.Modules {
		if k == moduleName || strings.HasPrefix(k, prefix) {
			delete(s.project.Modules, k)
		}
	}

	// Canonical installation folder
	modulesFolder := s.project.ModulesFolder
	if strings.TrimSpace(modulesFolder) == "" {
		modulesFolder = s.defaultModulesFolder
	}
	installFolder := filepath.ToSlash(filepath.Join(modulesFolder, moduleName, resolvedVersion))

	s.project.Modules[newKey] = eva.EvaModuleInfo{
		ModuleName:         moduleName,
		Version:            resolvedVersion,
		InstallationFolder: installFolder,
	}

	return nil
}

// --------------------------- IO helpers ---------------------------

func (s *ProjectBookkeepingService) ensureEvaFileExists(absEvaPath string) error {
	_, err := os.Stat(absEvaPath)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", absEvaPath, err)
	}

	// Ensure directory exists
	dir := filepath.Dir(absEvaPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	// Write default file
	defaultProject := eva.EvaProject{
		SchemaVersion: s.currentSchemaVersion,
		ModulesFolder: s.defaultModulesFolder,
		Modules:       map[string]eva.EvaModuleInfo{},
	}

	tmp := absEvaPath + ".tmp"
	if err := writeEvaProjectFile(tmp, defaultProject); err != nil {
		return err
	}
	if err := os.Rename(tmp, absEvaPath); err != nil {
		return fmt.Errorf("rename tmp -> eva.json: %w", err)
	}
	return nil
}

func cloneProject(p eva.EvaProject) eva.EvaProject {
	cp := p
	if cp.Modules == nil {
		cp.Modules = map[string]eva.EvaModuleInfo{}
		return cp
	}
	m := make(map[string]eva.EvaModuleInfo, len(cp.Modules))
	for k, v := range cp.Modules {
		m[k] = v
	}
	cp.Modules = m
	return cp
}

func writeEvaProjectFile(path string, p eva.EvaProject) error {
	b, err := jsonMarshalIndent(p)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	return nil
}

// Separate helper to keep imports clean if you already wrap json formatting elsewhere.
func jsonMarshalIndent(p eva.EvaProject) ([]byte, error) {
	// Inline import pattern avoided; keep it explicit.
	// If you already have a shared JSON util, swap this to your util.
	type marshaler interface {
		MarshalIndent(v any, prefix, indent string) ([]byte, error)
	}
	// Use standard library directly:
	b, err := utilsJSONMarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal eva project: %w", err)
	}
	return append(b, '\n'), nil
}

func utilsJSONMarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
