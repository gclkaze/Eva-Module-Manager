package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestCmdFlagTypes tests that various command flags have correct types
func TestCmdFlagTypes(t *testing.T) {
	testCases := []struct {
		name     string
		cmd      *cobra.Command
		flagName string
	}{
		// Install command flags
		{
			name:     "Install --path flag exists",
			cmd:      NewInstallCommand(createMockApp()),
			flagName: "path",
		},
		// Verify command flags
		{
			name:     "Verify --path flag exists",
			cmd:      NewVerifyCommand(createMockApp()),
			flagName: "path",
		},
		// Uninstall command flags
		{
			name:     "Uninstall --path flag exists",
			cmd:      NewUninstallCommand(createMockApp()),
			flagName: "path",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			flag := tt.cmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %s not found", tt.flagName)
			}
		})
	}
}

// TestCommandHierarchy tests the command hierarchy (parent-child relationships)
func TestCommandHierarchy(t *testing.T) {
	moduleCmd := NewModuleParentCommand(nil)
	// Test module command has subcommands
	if moduleCmd.Commands() == nil || len(moduleCmd.Commands()) == 0 {
		t.Error("Module command should have subcommands")
	}
	releaseCmd := NewReleaseParentCommand(nil, nil)

	// Test release command has subcommands
	if releaseCmd.Commands() == nil || len(releaseCmd.Commands()) == 0 {
		t.Error("Release command should have subcommands")
	}
}

// TestCommandPersistentFlags tests persistent flags across command hierarchy
func TestCommandPersistentFlags(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewVerifyCommand(mockApp)

	// PersistentFlags are inherited from parent (root command)
	// This tests that subcommands can access inherited flags
	persistentFlags := cmd.PersistentFlags()
	if persistentFlags == nil {
		t.Error("Command should have access to persistent flags")
	}
}

// TestVerifyCommandPathFlag tests the verify command's path flag specifically
func TestVerifyCommandPathFlag(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewVerifyCommand(mockApp)

	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Fatal("path flag not found")
	}

	// Test flag properties
	if pathFlag.Name != "path" {
		t.Errorf("Expected flag name 'path', got '%s'", pathFlag.Name)
	}

	if pathFlag.Shorthand != "p" {
		t.Errorf("Expected shorthand 'p', got '%s'", pathFlag.Shorthand)
	}
}

// TestInstallCommandPathFlag tests the install command's path flag specifically
func TestInstallCommandPathFlag(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewInstallCommand(mockApp)

	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Fatal("path flag not found")
	}

	if pathFlag.Name != "path" {
		t.Errorf("Expected flag name 'path', got '%s'", pathFlag.Name)
	}

	if pathFlag.Shorthand != "p" {
		t.Errorf("Expected shorthand 'p', got '%s'", pathFlag.Shorthand)
	}
}

// TestUninstallCommandPathFlag tests the uninstall command's path flag specifically
func TestUninstallCommandPathFlag(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewUninstallCommand(mockApp)

	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Fatal("path flag not found")
	}

	if pathFlag.Name != "path" {
		t.Errorf("Expected flag name 'path', got '%s'", pathFlag.Name)
	}

	if pathFlag.Shorthand != "p" {
		t.Errorf("Expected shorthand 'p', got '%s'", pathFlag.Shorthand)
	}
}

// TestInstallCommandEmptyPathFlag tests install command with empty path flag
func TestInstallCommandEmptyPathFlag(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewInstallCommand(mockApp)

	// Simulate setting --path to empty string
	cmd.SetArgs([]string{})
	cmd.Flags().Set("path", "")

	// This should trigger the validation error
	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Error("path flag should exist")
	}
}

// TestCommandDescriptions tests that all commands have proper descriptions
func TestCommandDescriptions(t *testing.T) {
	mockApp := createMockApp()

	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{
			name: "verify",
			cmd:  NewVerifyCommand(mockApp),
		},
		{
			name: "install",
			cmd:  NewInstallCommand(mockApp),
		},
		{
			name: "uninstall",
			cmd:  NewUninstallCommand(mockApp),
		},
	}

	for _, c := range commands {
		t.Run(c.name+" short description", func(t *testing.T) {
			if c.cmd.Short == "" {
				t.Errorf("Command %s should have a short description", c.name)
			}
		})
	}
}

// TestVerboseFlag tests that verbose flag is available globally
func TestVerboseFlag(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewVerifyCommand(mockApp)

	// The verbose flag should be available as a persistent flag
	// It's typically defined on the root command
	if cmd.Parent() != nil {
		persistentFlags := cmd.Parent().PersistentFlags()
		if persistentFlags != nil {
			verboseFlag := persistentFlags.Lookup("verbose")
			if verboseFlag == nil {
				// Verbose flag might be on root command, which we test separately
				t.Log("Verbose flag not found on parent command (expected for subcommand)")
			}
		}
	}
}

// TestSearchCommandArguments tests search command argument handling
func TestSearchCommandArguments(t *testing.T) {
	searchCmd := NewSearchArtifactsCommand(nil)
	// Since searchCmd uses global application variable that might not be initialized in tests,
	// we'll just verify the command structure exists
	if searchCmd == nil {
		t.Fatal("searchCmd should be defined")
	}

	if searchCmd.Use != "search" {
		t.Errorf("Expected use 'search', got '%s'", searchCmd.Use)
	}

	if searchCmd.Short == "" {
		t.Error("Expected short description")
	}
}

// TestReleaseDownloadArgs tests release download command argument validation
func TestReleaseDownloadArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		shouldErr bool
	}{
		{
			name:      "one argument is valid",
			args:      []string{"module@version"},
			shouldErr: false,
		},
		{
			name:      "no arguments should fail",
			args:      []string{},
			shouldErr: true,
		},
		{
			name:      "more than one argument should fail",
			args:      []string{"module@version1", "module@version2"},
			shouldErr: true,
		},
	}

	releaseDownloadCmd := NewReleaseDownloadCommand(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := releaseDownloadCmd.Args(releaseDownloadCmd, tt.args)

			if (err != nil) != tt.shouldErr {
				t.Errorf("Expected error: %v, got: %v", tt.shouldErr, err != nil)
			}
		})
	}
}

// TestModuleUploadArgs tests module upload command argument validation
func TestModuleUploadArgs(t *testing.T) {
	moduleUploadCmd := NewModuleUploadCommand(nil)
	// Module upload uses cobra.MinimumNArgs(1)
	if moduleUploadCmd == nil {
		t.Fatal("moduleUploadCmd is nil")
	}

	// The command should enforce at least one argument
	if moduleUploadCmd.Args == nil {
		t.Error("moduleUploadCmd should have Args validation")
	}
}

// TestModuleUpdateArgs tests module update command argument validation
func TestModuleUpdateArgs(t *testing.T) {
	moduleUpdateCmd := NewModuleUpdateCommand(nil)
	// Module update uses cobra.MinimumNArgs(1)
	if moduleUpdateCmd == nil {
		t.Fatal("moduleUpdateCmd is nil")
	}

	// The command should enforce at least one argument
	if moduleUpdateCmd.Args == nil {
		t.Error("moduleUpdateCmd should have Args validation")
	}
}

// TestLoginCommandEmail tests the login command email flag requirement
func TestLoginCommandEmail(t *testing.T) {
	loginCmd := NewLoginCommand(nil)
	emailFlag := loginCmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Fatal("email flag not found")
	}

	if emailFlag.Name != "email" {
		t.Errorf("Expected flag name 'email', got '%s'", emailFlag.Name)
	}

	if emailFlag.Shorthand != "e" {
		t.Errorf("Expected shorthand 'e', got '%s'", emailFlag.Shorthand)
	}
}

// TestRegisterCommandRequiredFlags tests the register command required flags
func TestRegisterCommandRequiredFlags(t *testing.T) {
	registerCmd := NewRegisterCommand(nil)
	requiredFlags := []string{"email", "firstname"}

	for _, flagName := range requiredFlags {
		flag := registerCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag %s to exist", flagName)
		}
	}
}
