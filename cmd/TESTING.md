# Eva Module Manager - Command Tests Documentation

## Overview

This comprehensive test suite for the Eva Module Manager (EMM) CLI application provides thorough coverage of all cobra commands in the `cmd` package. The test suite validates command structure, flags, arguments, and overall command hierarchy.

## Test Files

### 1. **root_test.go**
Contains the foundational testing infrastructure and helper functions used across all test files.

**Key Components:**
- `MockPrinter`: A mock implementation of the `output.Printer` interface for testing without external dependencies
- `createMockApp()`: Factory function that creates a mock EMMApp for testing
- `executeCommand()`: Helper function to execute commands and capture output
- `TestExecuteRoot()`: Validates the root command structure

**Test Coverage:**
- Root command creation and validation
- Mock printer initialization and all required methods
- Mock app initialization with all services

### 2. **commands_test.go**
Tests for the main user-facing commands: verify, install, and uninstall.

**Commands Tested:**
- **verify command**
  - `TestVerifyCommand()`: Validates command structure
  - `TestVerifyCommandFlags()`: Validates --path/-p flag existence
  - `TestVerifyCommandExecution()`: Tests command execution
  - `TestVerifyCommandPathFlag()`: Validates flag properties

- **install command**
  - `TestInstallCommand()`: Validates command structure
  - `TestInstallCommandFlags()`: Validates --path flag
  - `TestInstallCommandArgs()`: Tests argument validation
    - No arguments (valid)
    - One argument (valid)
    - Multiple arguments (invalid)
  - `TestInstallCommandPathFlag()`: Flag property validation
  - `TestInstallCommandEmptyPathFlag()`: Empty path handling

- **uninstall command**
  - `TestUninstallCommand()`: Validates command structure
  - `TestUninstallCommandArgs()`: Tests argument validation
    - One argument required (valid)
    - No arguments (invalid)
    - Multiple arguments (invalid)
  - `TestUninstallCommandPathFlag()`: Flag validation
  - `TestUninstallCommandFlags()`: Flag structure validation

**Common Tests:**
- `TestCommandStructure()`: Validates all main commands have proper Use, Short, and RunE fields
- `TestCommandStructure()`: Runs tests for all main command factories

### 3. **subcommands_test.go**
Tests for user authentication, search, and module management subcommands.

**Commands Tested:**
- **whoami**: Shows current active user
- **login**: User login command
  - `TestLoginCommand()`: Structure validation
  - `TestLoginCommandFlags()`: Email flag with -e shorthand
  - `TestLoginCommandEmail()`: Email flag details

- **register**: User registration command
  - `TestRegisterCommand()`: Structure validation
  - `TestRegisterCommandFlags()`: Email and firstname flags
  - `TestRegisterCommandRequiredFlags()`: Validates required flags

- **logout**: User logout command
  - `TestLogoutCommand()`: Structure validation
  - `TestLogoutCommandFlags()`: Email flag validation

- **search**: Module search command
  - `TestSearchCommand()`: Structure validation
  - `TestSearchCommandFlags()`: Tests --tags, --name, --description flags
  - `TestSearchCommandArguments()`: Validates command structure without executing

- **module**: Module management parent command
  - `TestModuleCommand()`: Validates it exists and has subcommands
  - **upload subcommand**
    - `TestModuleUploadCommand()`: Structure validation
    - `TestModuleUploadCommandFlags()`: --title, --repr, --tags, --description flags
    - `TestModuleUploadArgs()`: Minimum args validation

  - **update subcommand**
    - `TestModuleUpdateCommand()`: Structure validation
    - `TestModuleUpdateCommandFlags()`: --module-name, --title, --repr flags
    - `TestModuleUpdateArgs()`: Minimum args validation

  - **suggest subcommand**
    - `TestModuleSuggestCommand()`: Structure validation
    - `TestModuleSuggestCommandFlags()`: --module-name and --version flags

  - **mylist subcommand** (user modules)
    - `TestModuleGetUserCommand()`: Structure validation

### 4. **release_commands_test.go**
Tests for all release management commands.

**Commands Tested:**
- **release**: Parent command
  - `TestReleaseCommand()`: Validates parent command
  - **accept**: Accept suggested release
    - `TestReleaseAcceptCommand()`: Structure validation
    - `TestReleaseAcceptCommandArgs()`: Exact args=1 validation

  - **download**: Download a module release
    - `TestReleaseDownloadCommand()`: Structure validation
    - `TestReleaseDownloadCommandFlags()`: --savelocation/-s flag
    - `TestReleaseDownloadArgs()`: Exact args=1 validation

  - **dump**: Dump release information with filters
    - `TestReleaseDumpCommand()`: Structure validation
    - `TestReleaseDumpCommandFlags()`: Tests 11 different filter flags
      - --view, --status, --versions, --tags, --module, --repo
      - --description, --creator, --creator-email
      - --created-after, --released-after

  - **cancel**: Cancel a release
    - `TestReleaseCancelCommand()`: Structure validation

  - **reject**: Reject a release
    - `TestReleaseRejectCommand()`: Structure validation

  - **lower**: Lower release priority
    - `TestReleaseLowerCommand()`: Structure validation

- **switchuser**: Switch to different user
  - `TestSwitchUserCommand()`: Structure validation
  - `TestSwitchUserCommandFlags()`: Email flag with -e shorthand

### 5. **integration_test.go**
Tests for command integration, flags, and hierarchy.

**Test Categories:**

**Command Flag Testing:**
- `TestCmdFlagTypes()`: Validates main commands have expected flags
- `TestVerifyCommandPathFlag()`: Detailed flag property testing
- `TestInstallCommandPathFlag()`: Flag shorthand validation
- `TestUninstallCommandPathFlag()`: Flag configuration

**Command Hierarchy:**
- `TestCommandHierarchy()`: Validates parent-child relationships
- `TestCommandPersistentFlags()`: Tests flag inheritance

**Command Descriptions:**
- `TestCommandDescriptions()`: Validates all commands have descriptions
- `TestVerboseFlag()`: Tests verbose flag availability

**Argument Validation:**
- `TestReleaseDownloadArgs()`: ExactArgs validation
- `TestModuleUploadArgs()`: MinimumNArgs validation
- `TestModuleUpdateArgs()`: MinimumNArgs validation
- `TestSearchCommandArguments()`: Argument handling

## Test Infrastructure

### MockPrinter Implementation
All methods required by the `output.Printer` interface are implemented:
- `Info()`, `VerboseInfo()`
- `Error()`, `Warn()`, `VerboseWarn()`
- `Success()`
- `PrintModules()`, `PrintModuleInfo()`
- `PrintReleaseInfo()`, `PrintDetailedModuleReleaseInfo()`, `PrintReleaseRows()`
- `PrintSummary()`, `PrintUninstallSummary()`
- `GetVerboseFlagPointer()`

### Helper Functions
- `createMockApp()`: Creates a complete mock EMMApp with all services
- `executeCommand()`: Executes commands with captured output

## Test Execution

### Run All Tests
```bash
go test ./cmd -v
```

### Run Specific Test
```bash
go test ./cmd -v -run "TestVerifyCommand"
```

### Run with Coverage
```bash
go test ./cmd -v -cover
```

### Run Specific Test File
```bash
go test ./cmd -v -run "TestInstall"  # Runs all TestInstall* tests
```

## Test Statistics

The complete test suite includes:
- **68+ test functions**
- **Coverage of all cobra commands** in the cmd package
- **Validation of command structure, flags, and arguments**
- **Integration tests for command hierarchy**

## Test Design Patterns

### 1. Command Structure Validation
```go
func TestCommandName(t *testing.T) {
    cmd := cmdFactory(app)
    if cmd.Use != "expected-use" {
        t.Errorf("Use mismatch")
    }
    if cmd.Short == "" {
        t.Error("Missing short description")
    }
}
```

### 2. Flag Validation
```go
func TestCommandFlags(t *testing.T) {
    flag := cmd.Flags().Lookup("flagname")
    if flag == nil {
        t.Error("Flag missing")
    }
    if flag.Shorthand != "f" {
        t.Errorf("Shorthand incorrect")
    }
}
```

### 3. Argument Validation
```go
func TestCommandArgs(t *testing.T) {
    err := cmd.Args(cmd, args)
    if (err != nil) != shouldErr {
        t.Errorf("Unexpected error state")
    }
}
```

## Limitations

1. **Global State**: Some tests cannot exercise the full command Args validation due to global `application` variable initialization requirements
2. **Interactive Input**: Tests that require password input are skipped
3. **External Services**: Mock services are used, not actual backend connections

## Future Enhancements

1. Add integration tests with test fixtures
2. Add end-to-end tests with real file operations
3. Add benchmarking for command performance
4. Add fuzzing tests for argument validation
5. Increase coverage to test error paths more thoroughly

## Running Tests in CI/CD

The test suite is designed to run in CI/CD pipelines:

```bash
# Run all tests with coverage reporting
go test ./cmd -v -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run tests with short timeout for CI
go test ./cmd -timeout=30s -v

# Run tests in parallel
go test ./cmd -parallel 4 -v
```

## Contributing

When adding new commands:
1. Add command structure tests in the appropriate test file
2. Test all command flags and their shorthands
3. Test argument validation
4. Add the command to command hierarchy tests if it's a parent command
5. Ensure all tests pass before committing
