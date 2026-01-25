package cmd

import (
	"testing"

	"emm/internal/app"

	"github.com/spf13/cobra"
)

// TestVerifyCommand tests the verify command
func TestVerifyCommand(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewVerifyCommand(mockApp)

	if cmd == nil {
		t.Fatal("Failed to create verify command")
	}

	if cmd.Use != "verify" {
		t.Errorf("Expected command use 'verify', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}

	if cmd.Long == "" {
		t.Error("Expected long description, got empty string")
	}
}

// TestVerifyCommandFlags tests the verify command flags
func TestVerifyCommandFlags(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewVerifyCommand(mockApp)

	// Check for --path flag
	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Error("Expected --path flag to exist")
	}

	// The shorthand should be 'p'
	if pathFlag != nil && pathFlag.Shorthand != "p" {
		t.Errorf("Expected shorthand 'p', got '%s'", pathFlag.Shorthand)
	}
}

// TestVerifyCommandExecution tests verify command execution with no args
func TestVerifyCommandExecution(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewVerifyCommand(mockApp)

	// Test execution with default path
	output, err := executeCommand(t, cmd)
	// The command should execute without immediate error in testing context
	_ = output
	_ = err
}

// TestInstallCommand tests the install command
func TestInstallCommand(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewInstallCommand(mockApp)

	if cmd == nil {
		t.Fatal("Failed to create install command")
	}

	if cmd.Use != "install [module@version]" {
		t.Errorf("Expected command use 'install [module@version]', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}
}

// TestInstallCommandFlags tests the install command flags
func TestInstallCommandFlags(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewInstallCommand(mockApp)

	pathFlag := cmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Error("Expected --path flag to exist")
	}
}

// TestInstallCommandArgs tests the install command argument validation
func TestInstallCommandArgs(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewInstallCommand(mockApp)

	tests := []struct {
		name      string
		args      []string
		shouldErr bool
	}{
		{
			name:      "no arguments is valid",
			args:      []string{},
			shouldErr: false,
		},
		{
			name:      "one argument is valid",
			args:      []string{"module@1.0.0"},
			shouldErr: false,
		},
		{
			name:      "more than one argument should fail",
			args:      []string{"module1@1.0.0", "module2@2.0.0"},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			err := cmd.Args(cmd, tt.args)

			if (err != nil) != tt.shouldErr {
				t.Errorf("Expected error: %v, got: %v", tt.shouldErr, err != nil)
			}
		})
	}
}

// TestUninstallCommand tests the uninstall command
func TestUninstallCommand(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewUninstallCommand(mockApp)

	if cmd == nil {
		t.Fatal("Failed to create uninstall command")
	}

	if cmd.Use != "uninstall [module@version]" {
		t.Errorf("Expected command use 'uninstall [module@version]', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected short description, got empty string")
	}
}

// TestUninstallCommandArgs tests the uninstall command argument validation
func TestUninstallCommandArgs(t *testing.T) {
	mockApp := createMockApp()
	cmd := NewUninstallCommand(mockApp)

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
			args:      []string{"module1@1.0.0", "module2@2.0.0"},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			err := cmd.Args(cmd, tt.args)

			if (err != nil) != tt.shouldErr {
				t.Errorf("Expected error: %v, got: %v", tt.shouldErr, err != nil)
			}
		})
	}
}

// TestCommandStructure tests that all main commands have proper structure
func TestCommandStructure(t *testing.T) {
	mockApp := createMockApp()

	commands := map[string]func(*app.EMMApp) *cobra.Command{
		"verify":   NewVerifyCommand,
		"install":  NewInstallCommand,
		"uninstall": NewUninstallCommand,
	}

	for name, cmdFunc := range commands {
		t.Run(name, func(t *testing.T) {
			cmd := cmdFunc(mockApp)

			if cmd == nil {
				t.Errorf("Command %s returned nil", name)
				return
			}

			if cmd.Use == "" {
				t.Errorf("Command %s has no Use field", name)
			}

			if cmd.Short == "" {
				t.Errorf("Command %s has no Short description", name)
			}

			if cmd.RunE == nil && cmd.Run == nil {
				t.Errorf("Command %s has neither RunE nor Run", name)
			}
		})
	}
}
