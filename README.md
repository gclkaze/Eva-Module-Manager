# Eva Module Manager (EMM)

A command-line interface (CLI) tool for managing, discovering, and installing modules from a centralized repository. Built with Go and the Cobra CLI framework.

## Overview

**Eva Module Manager** enables developers to:
- 🔍 **Search and discover** modules in a repository
- 📦 **Upload and manage** their own modules
- 🎯 **Suggest, accept, and reject** module releases
- 🚀 **Install modules** into local projects
- 👥 **Manage user permissions** and access control

The system supports multiple user roles (Public, User, Maintainer, Admin) with distinct permissions and capabilities.

## Features

### For All Users
- **Module Search**: Find modules by tags, names, and components
- **View Module Information**: Browse module details and releases
- **Download Public Releases**: Access publicly available modules
- **Project Verification**: Validate eva.json project file structure

### For Registered Users
- **Upload Modules**: Contribute new modules to the repository
- **Manage Own Modules**: Upload, update, and delete your modules
- **Suggest Releases**: Propose new versions of your modules
- **Secure Credentials**: Tokens stored securely in OS keyring

### For Maintainers
- **Release Management**: Accept, reject, or cancel any module release
- **Release Filtering**: Filter and inspect all releases
- **Module Oversight**: Manage any module in the system

### For Administrators
- **User Management**: Ban/unban user accounts
- **Download Restrictions**: Access any release regardless of status
- **Full System Control**: Complete access to all operations

## Installation

### Prerequisites
- Go 1.21 or later
- Access to a running EMM backend server

### Build from Source

```bash
git clone https://github.com/yourusername/Eva-Module-Manager.git
cd Eva-Module-Manager
go build -o emm main.go
```

### Install Locally

```bash
go install ./...
```

## Quick Start

### 1. Create an Account

```bash
emm register \
  --email user@example.com \
  --handle "your-username" \
  --first-name "John" \
  --last-name "Doe"
```

You'll be prompted for a password.

### 2. Login

```bash
emm login --email user@example.com
```

Your credentials are securely stored in the system keyring.

### 3. Explore Modules

```bash
# Search for modules
emm search --tags "parser,json"

# View module information
emm info mymodule

# Check a specific release
emm release download mymodule@1.0.0 -s ./downloads
```

### 4. Upload Your First Module

```bash
emm module upload \
  --title "My Parser Module" \
  --repr "my-parser" \
  --tags "parser,json" \
  ./my-parser-v1.0.0.eva

# Suggest a release
emm module suggest \
  --module-id my-parser \
  --version 1.0.0

# View your modules
emm module mylist
```

## Command Reference

### Authentication

```bash
emm register [flags]          # Create a new account
emm login [flags]             # Login with email and password
emm logout                    # Logout and clear session
emm whoami                    # Show currently logged-in user
emm switchuser                # Switch between registered accounts
```

### Module Management

```bash
emm search [flags]            # Search modules by tags/components
emm info [module]             # View module information
emm module upload [flags]     # Upload a new module
emm module update [flags]     # Update module metadata
emm module mylist             # List your uploaded modules
emm module delete [flags]     # Delete your module
emm module suggest [flags]    # Suggest a release for your module
```

### Release Management

```bash
emm release download [module@version] [flags]  # Download a release
emm release dump [flags]                       # View/filter releases
emm release accept [flags]                     # Accept a suggested release (Maintainer+)
emm release reject [flags]                     # Reject a suggested release (Maintainer+)
emm release cancel [flags]                     # Cancel a pending release (Maintainer+)
```

### User Management (Admin Only)

```bash
emm user list                 # List all system users
emm user ban [flags]          # Ban a user account
emm user unban [flags]        # Restore a banned user
```

### Project Management

```bash
emm verify [flags]            # Validate eva.json project file
emm install [module@version]  # Install module to project
emm uninstall [module@version] # Remove module from project
```

### Verbose Output

All commands support the `-v` or `--verbose` flag for detailed debug output:

```bash
emm -v search --tags "parser"
```

## User Roles & Permissions

### 🌐 Public User (No Login Required)
- Search modules
- View module information
- Download public releases
- Verify project files

### 👤 Registered User
- All public operations
- Upload new modules
- Manage own modules
- Suggest releases for own modules
- Install modules to projects

### 🔧 Maintainer
- All user operations
- Accept/reject any release
- Cancel pending releases
- Filter and inspect all releases
- Manage any module

### 👨‍💼 Administrator
- All maintainer operations
- Ban/unban users
- Download restricted releases
- Full system access

See [SYSTEM_ARCHITECTURE.md](./SYSTEM_ARCHITECTURE.md) for detailed information about roles, permissions, and data flows.

## Configuration

EMM reads configuration from `internal/config/application.properties`:

```properties
# Server configuration
server.url=http://localhost:8080
server.api.base=/api

# Keyring settings
credentials.storage=system-keyring
```

You can override the server URL by setting environment variables or through command flags.

## Project Structure

```
eva-module-manager/
├── cmd/                  # Cobra commands
│   ├── root.go          # Root command definition
│   ├── module*.go       # Module management commands
│   ├── release*.go      # Release management commands
│   ├── user*.go         # User management commands
│   └── auth*.go         # Authentication commands
├── internal/
│   ├── app/            # EMMApp orchestrator
│   ├── backend/        # HTTP client & API communication
│   ├── config/         # Configuration management
│   ├── models/         # Data structures
│   ├── output/         # Console output & printing
│   └── services/       # Business logic services
├── pkg/
│   └── utils/          # Utilities & helpers
├── tests/              # Integration test files
└── main.go             # Entry point
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Distribution

```bash
go build -ldflags="-s -w" -o emm main.go
```

### Verbose/Debug Mode

Use the `-v` flag with any command for detailed debug output:

```bash
emm -v module upload --title "Test" --repr "test" ./test.eva
```

## Security

- **Credentials**: Securely stored in the system keyring (Windows Credential Manager, macOS Keychain, Linux secret-service)
- **Tokens**: JWT-based authentication for API requests
- **Ownership Validation**: Server enforces ownership checks on user operations
- **Password Input**: Passwords are read securely without echo to terminal

## Troubleshooting

### "User needs to be logged in" error

Solution: Run `emm login --email your@email.com` first

### "Module not found" error

Solution: Verify the module name with `emm search --tags "your-tag"`

### Credentials not found

Solution: Re-login with `emm login --email your@email.com`

### Server connection issues

Check your server URL in the configuration and ensure the backend is running:

```bash
emm -v search --tags "test"  # Verbose output shows connection details
```

## API Documentation

The CLI communicates with a backend server via REST API. Key endpoints include:

- `GET /api/modules` - List/search modules
- `GET /api/modules/:id` - Get module details
- `POST /api/modules/upload` - Upload new module
- `GET /api/releases/:id` - List releases
- `POST /api/releases/:id/accept/:releaseId` - Accept release (Admin+)
- `POST /api/user/ban/:userId` - Ban user (Admin only)

For complete API documentation, see the backend repository.

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass: `go test ./...`
2. Code is formatted: `go fmt ./...`
3. No linting errors: `golangci-lint run ./...`

## License

See LICENSE file for details.

## Support

For issues, feature requests, or questions:
1. Check [SYSTEM_ARCHITECTURE.md](./SYSTEM_ARCHITECTURE.md) for detailed system design
2. Review existing [test documentation](./cmd/TESTING.md) for command usage
3. Open an issue on the repository

## Version

Current version: 1.0.0

---

**Made with ❤️ for the Eva Module Manager community**