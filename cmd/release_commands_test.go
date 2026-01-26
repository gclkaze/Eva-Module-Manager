package cmd

import (
	"testing"
)

// TestModuleSuggestCommand tests the module suggest command
func TestModuleSuggestCommand(t *testing.T) {
	suggestionCmd := NewModuleSuggestionCommand(nil)
	if suggestionCmd == nil {
		t.Fatal("suggestionCmd is nil")
	}

	if suggestionCmd.Use != "suggest [params...]" {
		t.Errorf("Expected command use 'suggest [params...]', got '%s'", suggestionCmd.Use)
	}

	if suggestionCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if suggestionCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestModuleSuggestCommandFlags tests the module suggest command flags
func TestModuleSuggestCommandFlags(t *testing.T) {
	suggestionCmd := NewModuleSuggestionCommand(nil)
	moduleNameFlag := suggestionCmd.Flags().Lookup("module-name")
	if moduleNameFlag == nil {
		t.Error("Expected --module-name flag to exist")
	}

	versionFlag := suggestionCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Error("Expected --version flag to exist")
	}
}

// TestReleaseAcceptCommand tests the release accept command
func TestReleaseAcceptCommand(t *testing.T) {
	releaseAcceptCmd := NewReleaseAcceptCommand(nil)
	if releaseAcceptCmd == nil {
		t.Fatal("releaseAcceptCmd is nil")
	}

	if releaseAcceptCmd.Use != "accept module@version" {
		t.Errorf("Expected command use 'accept module@version', got '%s'", releaseAcceptCmd.Use)
	}

	if releaseAcceptCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if releaseAcceptCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestReleaseAcceptCommandArgs tests the release accept command argument validation
func TestReleaseAcceptCommandArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		shouldErr bool
	}{
		{
			name:      "one argument is valid",
			args:      []string{"module@1.0.0"},
			shouldErr: false,
		},
		{
			name:      "no arguments should fail",
			args:      []string{},
			shouldErr: true,
		},
		{
			name:      "more than one argument should fail",
			args:      []string{"module@1.0.0", "module@2.0.0"},
			shouldErr: true,
		},
	}

	releaseAcceptCmd := NewReleaseAcceptCommand(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := releaseAcceptCmd.Args(releaseAcceptCmd, tt.args)

			if (err != nil) != tt.shouldErr {
				t.Errorf("Expected error: %v, got: %v", tt.shouldErr, err != nil)
			}
		})
	}
}

// TestReleaseDumpCommand tests the release dump command
func TestReleaseDumpCommand(t *testing.T) {
	releaseDumpCmd := NewReleaseDumpCommand(nil)
	if releaseDumpCmd == nil {
		t.Fatal("releaseDumpCmd is nil")
	}

	if releaseDumpCmd.Use != "dump [params...]" {
		t.Errorf("Expected command use 'dump [params...]', got '%s'", releaseDumpCmd.Use)
	}

	if releaseDumpCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if releaseDumpCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestReleaseDumpCommandFlags tests the release dump command flags
func TestReleaseDumpCommandFlags(t *testing.T) {
	releaseDumpCmd := NewReleaseDumpCommand(nil)

	viewFlag := releaseDumpCmd.Flags().Lookup("view")
	if viewFlag == nil {
		t.Error("Expected --view flag to exist")
	}

	statusFlag := releaseDumpCmd.Flags().Lookup("status")
	if statusFlag == nil {
		t.Error("Expected --status flag to exist")
	}

	versionsFlag := releaseDumpCmd.Flags().Lookup("versions")
	if versionsFlag == nil {
		t.Error("Expected --versions flag to exist")
	}

	tagsFlag := releaseDumpCmd.Flags().Lookup("tags")
	if tagsFlag == nil {
		t.Error("Expected --tags flag to exist")
	}

	moduleFlag := releaseDumpCmd.Flags().Lookup("module")
	if moduleFlag == nil {
		t.Error("Expected --module flag to exist")
	}

	repoFlag := releaseDumpCmd.Flags().Lookup("repo")
	if repoFlag == nil {
		t.Error("Expected --repo flag to exist")
	}

	descriptionFlag := releaseDumpCmd.Flags().Lookup("description")
	if descriptionFlag == nil {
		t.Error("Expected --description flag to exist")
	}

	creatorFlag := releaseDumpCmd.Flags().Lookup("creator")
	if creatorFlag == nil {
		t.Error("Expected --creator flag to exist")
	}

	creatorEmailFlag := releaseDumpCmd.Flags().Lookup("creator-email")
	if creatorEmailFlag == nil {
		t.Error("Expected --creator-email flag to exist")
	}

	createdAfterFlag := releaseDumpCmd.Flags().Lookup("created-after")
	if createdAfterFlag == nil {
		t.Error("Expected --created-after flag to exist")
	}

	releasedAfterFlag := releaseDumpCmd.Flags().Lookup("released-after")
	if releasedAfterFlag == nil {
		t.Error("Expected --released-after flag to exist")
	}
}

// TestReleaseCancelCommand tests the release cancel command if it exists
func TestReleaseCancelCommand(t *testing.T) {
	releaseCancelCmd := NewReleaseCancelCommand(nil)

	if releaseCancelCmd == nil {
		t.Fatal("releaseCancelCmd is nil")
	}

	if releaseCancelCmd.Use == "" {
		t.Error("Expected command use, got empty string")
	}

	if releaseCancelCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if releaseCancelCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestReleaseRejectCommand tests the release reject command if it exists
func TestReleaseRejectCommand(t *testing.T) {
	releaseRejectCmd := NewReleaseRejectCommand(nil)

	if releaseRejectCmd == nil {
		t.Fatal("releaseRejectCmd is nil")
	}

	if releaseRejectCmd.Use == "" {
		t.Error("Expected command use, got empty string")
	}

	if releaseRejectCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if releaseRejectCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestReleaseLowerCommand tests the release lower command if it exists
func TestReleaseLowerCommand(t *testing.T) {
	releaseLowerCmd := NewReleaseLowerCommand(nil)

	if releaseLowerCmd == nil {
		t.Fatal("releaseLowerCmd is nil")
	}

	if releaseLowerCmd.Use == "" {
		t.Error("Expected command use, got empty string")
	}

	if releaseLowerCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if releaseLowerCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestSwitchUserCommand tests the switch user command if it exists
func TestSwitchUserCommand(t *testing.T) {
	switchCmd := NewSwitchUserCommand(nil)
	if switchCmd == nil {
		t.Fatal("switchCmd is nil")
	}

	if switchCmd.Use == "" {
		t.Error("Expected command use, got empty string")
	}

	if switchCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if switchCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}

// TestSwitchUserCommandFlags tests the switch user command flags
func TestSwitchUserCommandFlags(t *testing.T) {
	switchCmd := NewSwitchUserCommand(nil)

	emailFlag := switchCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Error("Expected --email flag to exist")
	}

	// Verify the shorthand is correctly set
	if emailFlag != nil && emailFlag.Shorthand != "e" {
		t.Errorf("Expected shorthand 'e', got '%s'", emailFlag.Shorthand)
	}
}

// TestModuleGetUserCommand tests the module get user command if it exists
func TestModuleGetUserCommand(t *testing.T) {
	userModuleGetCmd := NewModuleUserGetCommand(nil)
	if userModuleGetCmd == nil {
		t.Fatal("userModuleGetCmd is nil")
	}

	if userModuleGetCmd.Use == "" {
		t.Error("Expected command use, got empty string")
	}

	if userModuleGetCmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if userModuleGetCmd.RunE == nil {
		t.Error("Expected RunE function to be defined")
	}
}
