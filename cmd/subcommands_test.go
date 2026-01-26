package cmd

import (
	"testing"
)

// TestWhoamiCommand tests the whoami command
func TestWhoamiCommand(t *testing.T) {
	// whoamiCmd is a package-level variable defined in whoami.go
	whoamiCmd := NewWhoamiCommand(nil)
	if whoamiCmd == nil {
		t.Fatal("whoamiCmd is nil")
	}

	if whoamiCmd.Use != "whoami" {
		t.Errorf("Expected command use 'whoami', got '%s'", whoamiCmd.Use)
	}

	if whoamiCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if whoamiCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestLoginCommand tests the login command
func TestLoginCommand(t *testing.T) {
	loginCmd := NewLoginCommand(nil)
	if loginCmd == nil {
		t.Fatal("loginCmd is nil")
	}

	if loginCmd.Use != "login" {
		t.Errorf("Expected command use 'login', got '%s'", loginCmd.Use)
	}

	if loginCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if loginCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestLoginCommandFlags tests the login command flags
func TestLoginCommandFlags(t *testing.T) {
	loginCmd := NewLoginCommand(nil)
	emailFlag := loginCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Error("Expected --email flag to exist")
	}

	// Verify the shorthand is correctly set
	if emailFlag != nil && emailFlag.Shorthand != "e" {
		t.Errorf("Expected shorthand 'e', got '%s'", emailFlag.Shorthand)
	}
}

// TestRegisterCommand tests the register command
func TestRegisterCommand(t *testing.T) {
	if registerCmd == nil {
		t.Fatal("registerCmd is nil")
	}

	if registerCmd.Use != "register" {
		t.Errorf("Expected command use 'register', got '%s'", registerCmd.Use)
	}

	if registerCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if registerCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestRegisterCommandFlags tests the register command flags
func TestRegisterCommandFlags(t *testing.T) {
	emailFlag := registerCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Error("Expected --email flag to exist")
	}

	firstnameFlag := registerCmd.Flags().Lookup("firstname")
	if firstnameFlag == nil {
		t.Error("Expected --firstname flag to exist")
	}
}

// TestLogoutCommand tests the logout command
func TestLogoutCommand(t *testing.T) {
	if logoutCmd == nil {
		t.Fatal("logoutCmd is nil")
	}

	if logoutCmd.Use != "logout" {
		t.Errorf("Expected command use 'logout', got '%s'", logoutCmd.Use)
	}

	if logoutCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if logoutCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestLogoutCommandFlags tests the logout command flags
func TestLogoutCommandFlags(t *testing.T) {
	emailFlag := logoutCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Error("Expected --email flag to exist")
	}

	// Verify the shorthand is correctly set
	if emailFlag != nil && emailFlag.Shorthand != "e" {
		t.Errorf("Expected shorthand 'e', got '%s'", emailFlag.Shorthand)
	}
}

// TestSearchCommand tests the search command
func TestSearchCommand(t *testing.T) {
	if searchCmd == nil {
		t.Fatal("searchCmd is nil")
	}

	if searchCmd.Use != "search" {
		t.Errorf("Expected command use 'search', got '%s'", searchCmd.Use)
	}

	if searchCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if searchCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestSearchCommandFlags tests the search command flags
func TestSearchCommandFlags(t *testing.T) {
	tagsFlag := searchCmd.Flags().Lookup("tags")
	if tagsFlag == nil {
		t.Error("Expected --tags flag to exist")
	}

	nameFlag := searchCmd.Flags().Lookup("name")
	if nameFlag == nil {
		t.Error("Expected --name flag to exist")
	}

	descriptionFlag := searchCmd.Flags().Lookup("description")
	if descriptionFlag == nil {
		t.Error("Expected --description flag to exist")
	}
}

// TestModuleCommand tests the module command
func TestModuleCommand(t *testing.T) {
	if moduleCmd == nil {
		t.Fatal("moduleCmd is nil")
	}

	if moduleCmd.Use != "module" {
		t.Errorf("Expected command use 'module', got '%s'", moduleCmd.Use)
	}

	if moduleCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}
}

// TestReleaseCommand tests the release command
func TestReleaseCommand(t *testing.T) {
	if releaseCmd == nil {
		t.Fatal("releaseCmd is nil")
	}

	if releaseCmd.Use != "release" {
		t.Errorf("Expected command use 'release', got '%s'", releaseCmd.Use)
	}

	if releaseCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}
}

// TestInfoCommand tests the info command
func TestInfoCommand(t *testing.T) {
	showCmd := NewShowModuleInfoCommand(nil)
	if showCmd == nil {
		t.Fatal("showCmd is nil")
	}

	if showCmd.Use != "info" {
		t.Errorf("Expected command use 'info', got '%s'", showCmd.Use)
	}

	if showCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if showCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestModuleUploadCommand tests the module upload command
func TestModuleUploadCommand(t *testing.T) {
	if moduleUploadCmd == nil {
		t.Fatal("moduleUploadCmd is nil")
	}

	if moduleUploadCmd.Use != "upload [params...]" {
		t.Errorf("Expected command use 'upload [params...]', got '%s'", moduleUploadCmd.Use)
	}

	if moduleUploadCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if moduleUploadCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestModuleUploadCommandFlags tests the module upload command flags
func TestModuleUploadCommandFlags(t *testing.T) {
	titleFlag := moduleUploadCmd.Flags().Lookup("title")
	if titleFlag == nil {
		t.Error("Expected --title flag to exist")
	}

	reprFlag := moduleUploadCmd.Flags().Lookup("repr")
	if reprFlag == nil {
		t.Error("Expected --repr flag to exist")
	}

	tagsFlag := moduleUploadCmd.Flags().Lookup("tags")
	if tagsFlag == nil {
		t.Error("Expected --tags flag to exist")
	}

	descriptionFlag := moduleUploadCmd.Flags().Lookup("description")
	if descriptionFlag == nil {
		t.Error("Expected --description flag to exist")
	}
}

// TestModuleUpdateCommand tests the module update command
func TestModuleUpdateCommand(t *testing.T) {
	if moduleUpdateCmd == nil {
		t.Fatal("moduleUpdateCmd is nil")
	}

	if moduleUpdateCmd.Use != "update [params...]" {
		t.Errorf("Expected command use 'update [params...]', got '%s'", moduleUpdateCmd.Use)
	}

	if moduleUpdateCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if moduleUpdateCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestModuleUpdateCommandFlags tests the module update command flags
func TestModuleUpdateCommandFlags(t *testing.T) {
	moduleNameFlag := moduleUpdateCmd.Flags().Lookup("module-name")
	if moduleNameFlag == nil {
		t.Error("Expected --module-name flag to exist")
	}

	titleFlag := moduleUpdateCmd.Flags().Lookup("title")
	if titleFlag == nil {
		t.Error("Expected --title flag to exist")
	}

	reprFlag := moduleUpdateCmd.Flags().Lookup("repr")
	if reprFlag == nil {
		t.Error("Expected --repr flag to exist")
	}
}

// TestReleaseDownloadCommand tests the release download command
func TestReleaseDownloadCommand(t *testing.T) {
	if releaseDownloadCmd == nil {
		t.Fatal("releaseDownloadCmd is nil")
	}

	if releaseDownloadCmd.Use != "download module@version" {
		t.Errorf("Expected command use 'download module@version', got '%s'", releaseDownloadCmd.Use)
	}

	if releaseDownloadCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if releaseDownloadCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestReleaseDownloadCommandFlags tests the release download command flags
func TestReleaseDownloadCommandFlags(t *testing.T) {
	savelocationFlag := releaseDownloadCmd.Flags().Lookup("savelocation")
	if savelocationFlag == nil {
		t.Error("Expected --savelocation flag to exist")
	}

	// Verify the shorthand is correctly set
	if savelocationFlag != nil && savelocationFlag.Shorthand != "s" {
		t.Errorf("Expected shorthand 's', got '%s'", savelocationFlag.Shorthand)
	}
}
