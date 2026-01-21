package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyExisting_WithDirContainingEvaJSON_OK(t *testing.T) {
	td := t.TempDir()
	writeEvaJSON(t, filepath.Join(td, "eva.json"), validFloatingEvaJSON())

	svc, err := NewProjectBookkeepingService(td)
	if err != nil {
		t.Fatalf("NewProjectBookkeepingService: %v", err)
	}
	svc.SetProperties(nil)

	abs, err := svc.VerifyExisting(td)
	if err != nil {
		t.Fatalf("VerifyExisting: expected no error, got %v", err)
	}

	want := filepath.Join(td, "eva.json")
	if !samePath(abs, want) {
		t.Fatalf("expected %q, got %q", want, abs)
	}
}

func TestVerifyExisting_WithFilePath_OK(t *testing.T) {
	td := t.TempDir()
	p := filepath.Join(td, "eva.json")
	writeEvaJSON(t, p, validPinnedEvaJSON())

	svc, err := NewProjectBookkeepingService(td)
	if err != nil {
		t.Fatalf("NewProjectBookkeepingService: %v", err)
	}
	svc.SetProperties(nil)

	abs, err := svc.VerifyExisting(p)
	if err != nil {
		t.Fatalf("VerifyExisting: expected no error, got %v", err)
	}
	if !samePath(abs, p) {
		t.Fatalf("expected %q, got %q", p, abs)
	}
}

func TestVerifyExisting_WithDirMissingEvaJSON_Err(t *testing.T) {
	td := t.TempDir()

	svc, err := NewProjectBookkeepingService(td)
	if err != nil {
		t.Fatalf("NewProjectBookkeepingService: %v", err)
	}
	svc.SetProperties(nil)

	_, err = svc.VerifyExisting(td)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "eva.json not found") {
		t.Fatalf("expected 'eva.json not found' error, got %v", err)
	}
}

func TestVerifyExisting_WithNonEvaJSONFile_Err(t *testing.T) {
	td := t.TempDir()
	other := filepath.Join(td, "not-eva.json")
	if err := os.WriteFile(other, []byte(`{}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	svc, err := NewProjectBookkeepingService(td)
	if err != nil {
		t.Fatalf("NewProjectBookkeepingService: %v", err)
	}
	svc.SetProperties(nil)

	_, err = svc.VerifyExisting(other)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "expected eva.json") {
		t.Fatalf("expected 'expected eva.json' error, got %v", err)
	}
}

func TestVerifyExisting_WithInvalidEvaJSON_Err(t *testing.T) {
	td := t.TempDir()
	p := filepath.Join(td, "eva.json")
	writeEvaJSON(t, p, `{ not json }`)

	svc, err := NewProjectBookkeepingService(td)
	if err != nil {
		t.Fatalf("NewProjectBookkeepingService: %v", err)
	}
	svc.SetProperties(nil)

	_, err = svc.VerifyExisting(td)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "invalid eva.json") {
		t.Fatalf("expected 'invalid eva.json' error, got %v", err)
	}
}

func TestVerifyExisting_DefaultPath_UsesCWD(t *testing.T) {
	td := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(orig) }()

	if err := os.Chdir(td); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	writeEvaJSON(t, filepath.Join(td, "eva.json"), validFloatingEvaJSON())

	svc, err := NewProjectBookkeepingService(td)
	if err != nil {
		t.Fatalf("NewProjectBookkeepingService: %v", err)
	}
	svc.SetProperties(nil)

	abs, err := svc.VerifyExisting("")
	if err != nil {
		t.Fatalf("VerifyExisting: expected no error, got %v", err)
	}

	want := filepath.Join(td, "eva.json")
	if !samePath(abs, want) {
		t.Fatalf("expected %q, got %q", want, abs)
	}
}

/* helpers */

func writeEvaJSON(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func validFloatingEvaJSON() string {
	return `{
  "schemaVersion": 1,
  "modulesFolder": "eva-modules",
  "modules": {
    "foo": { "moduleName": "foo" },
    "bar@latest": { "moduleName": "bar", "version": "latest", "installationFolder": "some/relative/path" }
  }
}`
}

func validPinnedEvaJSON() string {
	return `{
  "schemaVersion": 1,
  "modulesFolder": "eva-modules",
  "modules": {
    "baz@v1.2.3": { "moduleName": "baz", "version": "v1.2.3", "installationFolder": "eva-modules/baz/v1.2.3" }
  }
}`
}

func samePath(a, b string) bool {
	aa, err1 := filepath.Abs(a)
	bb, err2 := filepath.Abs(b)
	if err1 != nil || err2 != nil {
		return a == b
	}
	return strings.EqualFold(filepath.Clean(aa), filepath.Clean(bb))
}
