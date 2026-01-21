package test

import (
	"emm/internal/backend"
	"emm/pkg/utils"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/magiconair/properties"
)

func TestEvaJSONParser_ParseBytes_BlankDefaults(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	prj, err := p.ParseBytes([]byte("   \n\t"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if prj.SchemaVersion != p.CurrentSchemaVersion {
		t.Fatalf("schemaVersion: expected %d, got %d", p.CurrentSchemaVersion, prj.SchemaVersion)
	}
	if prj.ModulesFolder != p.DefaultModulesFolder {
		t.Fatalf("modulesFolder: expected %q, got %q", p.DefaultModulesFolder, prj.ModulesFolder)
	}
	if prj.Modules == nil || len(prj.Modules) != 0 {
		t.Fatalf("modules: expected empty map, got %#v", prj.Modules)
	}
}

func TestEvaJSONParser_ParseBytes_InvalidJSON(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	_, err := p.ParseBytes([]byte(`{not json}`))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "invalid eva.json") {
		t.Fatalf("expected ErrInvalidEvaJSON, got %v", err)
	}
}

func TestEvaJSONParser_ParseBytes_SchemaVersionTooNew(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := `{
	  "schemaVersion": 999,
	  "modulesFolder": "eva-modules",
	  "modules": {}
	}`

	_, err := p.ParseBytes([]byte(json))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "newer than supported") {
		t.Fatalf("expected schema version error, got %v", err)
	}
}

func TestEvaJSONParser_ParseBytes_ModulesFolderInvalidAbsolute(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := fmt.Sprintf(`{
	  "schemaVersion": 1,
	  "modulesFolder": %q,
	  "modules": {}
	}`, absModulesFolderForTest())

	_, err := p.ParseBytes([]byte(json))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "modulesFolder must be a relative path") {
		t.Fatalf("expected relative path error, got %v", err)
	}
}

func TestEvaJSONParser_ParseBytes_ModulesFolderInvalidDotDot(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "../eva-modules",
	  "modules": {}
	}`

	_, err := p.ParseBytes([]byte(json))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "must not contain '..'") {
		t.Fatalf("expected '..' error, got %v", err)
	}
}

func TestEvaJSONParser_ParseBytes_Floating_NoVersionKey_AllowsEmptyVersionAndNoInstallFolder(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "foo": { "moduleName": "foo" }
	  }
	}`

	prj, err := p.ParseBytes([]byte(json))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mi, ok := prj.Modules["foo"]
	if !ok {
		t.Fatalf("expected key %q to exist", "foo")
	}
	if mi.ModuleName != "foo" {
		t.Fatalf("moduleName: expected %q, got %q", "foo", mi.ModuleName)
	}
	if mi.Version != "" {
		t.Fatalf("version: expected empty, got %q", mi.Version)
	}
	if mi.InstallationFolder != "" {
		t.Fatalf("installationFolder: expected empty for floating, got %q", mi.InstallationFolder)
	}
}

func TestEvaJSONParser_ParseBytes_Floating_Latest_AllowsLatestAndOptionalInstallFolderValidated(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "bar@latest": { "moduleName": "bar", "version": "latest", "installationFolder": "some/relative/path" }
	  }
	}`

	prj, err := p.ParseBytes([]byte(json))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mi, ok := prj.Modules["bar@latest"]
	if !ok {
		t.Fatalf("expected key %q to exist", "bar@latest")
	}
	if mi.ModuleName != "bar" || mi.Version != "latest" {
		t.Fatalf("expected bar@latest, got %s@%s", mi.ModuleName, mi.Version)
	}
	// For floating, parser does not force canonical installation folder.
	if mi.InstallationFolder != "some/relative/path" {
		t.Fatalf("installationFolder: expected preserved %q, got %q", "some/relative/path", mi.InstallationFolder)
	}
}

func TestEvaJSONParser_ParseBytes_Floating_Latest_RejectsAbsoluteInstallFolder(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := fmt.Sprintf(`{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "bar@latest": { "moduleName": "bar", "version": "latest", "installationFolder": %q }
	  }
	}`, absPathForTest())

	_, err := p.ParseBytes([]byte(json))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "installationFolder must be relative") {
		t.Fatalf("expected relative installationFolder error, got %v", err)
	}
}
func TestEvaJSONParser_ParseBytes_Pinned_SelfHealsInstallationFolder(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	// NOTE: ParseModuleReleaseVersion normalizes semver to have "v" prefix.
	// We intentionally use a key without "v" to ensure the parser still accepts it
	// and canonicalizes the installation folder based on the normalized version.
	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "baz@1.2.3": { "moduleName": "baz", "version": "v1.2.3", "installationFolder": "wrong/path" }
	  }
	}`

	prj, err := p.ParseBytes([]byte(json))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mi, ok := prj.Modules["baz@1.2.3"]
	if !ok {
		t.Fatalf("expected key %q to exist", "baz@1.2.3")
	}

	want := "eva-modules/baz/v1.2.3"
	if mi.InstallationFolder != want {
		t.Fatalf("installationFolder: expected %q, got %q", want, mi.InstallationFolder)
	}
}

func TestEvaJSONParser_ParseBytes_Pinned_InferVersionFromKeyWhenMissing(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	// version omitted in object; inferred from key.
	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "baz@1.2.3": { "moduleName": "baz" }
	  }
	}`

	prj, err := p.ParseBytes([]byte(json))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mi := prj.Modules["baz@1.2.3"]
	if mi.Version != "v1.2.3" {
		t.Fatalf("version: expected %q, got %q", "v1.2.3", mi.Version)
	}
	if mi.InstallationFolder != "eva-modules/baz/v1.2.3" {
		t.Fatalf("installationFolder: expected %q, got %q", "eva-modules/baz/v1.2.3", mi.InstallationFolder)
	}
}

func TestEvaJSONParser_ParseBytes_InconsistentModuleNameRejected(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "foo": { "moduleName": "bar" }
	  }
	}`

	_, err := p.ParseBytes([]byte(json))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "inconsistent fields") {
		t.Fatalf("expected inconsistent fields error, got %v", err)
	}
}

func TestEvaJSONParser_ParseBytes_InconsistentVersionRejected(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	json := `{
	  "schemaVersion": 1,
	  "modulesFolder": "eva-modules",
	  "modules": {
	    "foo@latest": { "moduleName": "foo", "version": "v1.2.3" }
	  }
	}`

	_, err := p.ParseBytes([]byte(json))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "inconsistent fields") {
		t.Fatalf("expected inconsistent fields error, got %v", err)
	}
}

func TestEvaJSONParser_ParseBytes_NormalizesTopLevelDefaultsWhenMissing(t *testing.T) {
	p := utils.NewEvaJSONParser(nil)

	// schemaVersion and modulesFolder missing
	json := `{
	  "modules": {}
	}`

	prj, err := p.ParseBytes([]byte(json))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if prj.SchemaVersion != 1 {
		t.Fatalf("schemaVersion: expected default 1, got %d", prj.SchemaVersion)
	}
	if prj.ModulesFolder == "" {
		t.Fatalf("modulesFolder: expected non-empty default")
	}
}

func TestEvaJSONParser_ParseBytes_UsesConfiguredDefaultsFromProps(t *testing.T) {
	props := propertiesForTest(map[string]string{
		backend.EVA_JSON_SCHEMA_CURRENT_VERSION_KEY: "7",
		backend.EVA_DEFAULT_EVA_MODULES_KEY:         "my-mods",
	})
	p := utils.NewEvaJSONParser(props)

	prj, err := p.ParseBytes([]byte(" \n"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if prj.SchemaVersion != 7 {
		t.Fatalf("schemaVersion: expected %d, got %d", 7, prj.SchemaVersion)
	}
	if prj.ModulesFolder != "my-mods" {
		t.Fatalf("modulesFolder: expected %q, got %q", "my-mods", prj.ModulesFolder)
	}
}

/* ----------------------------- test helpers ----------------------------- */

// propertiesForTest avoids pulling in magiconair/properties directly in every test.
// Your project already uses it, so this is fine.
func propertiesForTest(kv map[string]string) *properties.Properties {
	// Import locally to keep the top of file tidy.
	// (Go requires imports at top; so we include it there if you prefer.)
	p := properties.NewProperties()
	for k, v := range kv {
		_, _, _ = p.Set(k, v)
	}
	return p
}
func absPathForTest() string {
	if runtime.GOOS == "windows" {
		// Use something that IsAbs will always treat as absolute on Windows.
		return `C:\abs\path`
	}
	return "/abs/path"
}

func absModulesFolderForTest() string {
	// Same as above; separated just in case you want different values later.
	return absPathForTest()
}
