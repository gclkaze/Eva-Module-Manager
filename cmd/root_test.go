package cmd

import (
	"bytes"
	"emm/internal/app"
	"emm/internal/models"
	"emm/internal/services"
	"testing"

	"github.com/spf13/cobra"
)

// MockPrinter is a mock implementation of output.Printer for testing
type MockPrinter struct {
	verboseFlag bool
	infoLogs    []string
	errorLogs   []string
	warnLogs    []string
}

func NewMockPrinter() *MockPrinter {
	return &MockPrinter{
		infoLogs:  []string{},
		errorLogs: []string{},
		warnLogs:  []string{},
	}
}

func (m *MockPrinter) Info(msg string) {
	m.infoLogs = append(m.infoLogs, msg)
}

func (m *MockPrinter) VerboseInfo(msg string) {
	if m.verboseFlag {
		m.infoLogs = append(m.infoLogs, msg)
	}
}

func (m *MockPrinter) Error(err error) {
	if err != nil {
		m.errorLogs = append(m.errorLogs, err.Error())
	}
}

func (m *MockPrinter) Warn(msg string) {
	m.warnLogs = append(m.warnLogs, msg)
}

func (m *MockPrinter) VerboseWarn(msg string) {
	if m.verboseFlag {
		m.warnLogs = append(m.warnLogs, msg)
	}
}

func (m *MockPrinter) Success(msg string) {
	m.infoLogs = append(m.infoLogs, msg)
}

func (m *MockPrinter) PrintModules(mods []models.Module) {
	// Mock implementation
}

func (m *MockPrinter) PrintModuleInfo(moduleInfo models.ModuleEnrichedInformation) {
	// Mock implementation
}

func (m *MockPrinter) PrintReleaseInfo(moduleRepr string, r models.Release) {
	// Mock implementation
}

func (m *MockPrinter) PrintSummary(summary *models.InstallationSummary) {
	// Mock implementation
}

func (m *MockPrinter) PrintUninstallSummary(summary *models.PurgeSummary) {
	// Mock implementation
}

func (m *MockPrinter) PrintDetailedModuleReleaseInfo(mods []models.ModuleEnrichedDTO) {
	// Mock implementation
}

func (m *MockPrinter) PrintReleaseRows(mods []models.ModuleEnrichedDTO) {
	// Mock implementation
}

func (m *MockPrinter) GetVerboseFlagPointer() *bool {
	return &m.verboseFlag
}

// TestExecuteRoot tests the root command execution
func TestExecuteRoot(t *testing.T) {
	// Create a test command
	cmd := &cobra.Command{
		Use:   "emm",
		Short: "EMM is a CLI tool",
		Long:  "EMM is a CLI tool that can search artifacts using tags",
	}

	if cmd == nil {
		t.Fatal("Failed to create root command")
	}

	if cmd.Use != "emm" {
		t.Errorf("Expected command use 'emm', got '%s'", cmd.Use)
	}
}

// Helper function to create a mock EMMApp for testing
func createMockApp() *app.EMMApp {
	cwd := "/tmp"
	moduleSearchService := services.NewModuleSearchService()
	authService := services.NewAuthService()
	moduleService := services.NewModuleService(authService)
	releaseService := services.NewModuleReleaseService(authService)
	bookKeepingService := &services.ProjectBookkeepingService{}
	installService := services.NewInstallService(cwd, bookKeepingService, releaseService)

	printer := NewMockPrinter()

	app := app.NewEMMApp(
		cwd,
		moduleSearchService,
		authService,
		moduleService,
		releaseService,
		bookKeepingService,
		installService,
		printer,
	)

	return app
}

// Helper function to execute a command and capture output
func executeCommand(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}
