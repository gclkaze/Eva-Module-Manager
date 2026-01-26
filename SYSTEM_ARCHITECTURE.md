# Eva Module Manager - System Architecture & User Roles

## System Overview

**Eva Module Manager (EMM)** is a command-line interface (CLI) application that manages module packages and their releases in a centralized repository server. It enables users to search for modules, upload new modules, manage releases, and install modules into local projects.

The system operates as a client-side CLI tool that communicates with a backend module repository server via REST API.

---

## User Types & Roles

### 1. **Public Users (Unauthenticated)**

Public users are individuals who have not registered or logged in. They have read-only access to the repository for discovery and download without requiring authentication.

#### Characteristics:
- **No Authentication**: No registered account or login session
- **No Token**: Operate without authentication
- **Read-Only Access**: Can only query and download publicly available content
- **No Persistence**: No stored credentials or session state
- **Unrestricted Discovery**: Full access to search and discovery features

#### Public User Operations (No Login Required):

**Module Discovery:**
- `GET /api/modules/:id` - Find module by ID
- `GET /api/modules/search` - Search modules by tags, name, components
- `GET /api/modules/info` - View module information

**Release Information:**
- `GET /api/releases/latest/:module` - Get latest release of a module
- `GET /api/releases/:id` - View all releases of a module
- `GET /api/releases/:id/release/:releaseId` - View specific release information
- `GET /api/releases/:id/search` - Search releases by keywords
- `GET /api/download/release/:release` - Download public releases

**Project Verification:**
- `verify [--path eva.json]` - Validate eva.json structure (local operation)

#### Token Handling:
```go
// Public access: No token required or token set to empty string
// GET /api/download/release/:release  - No authentication needed
tk := ""  // Public access - no token
err = application.DownloadPublicRelease(ctx, releaseId)
```

#### Example: Public Module Search & Download (No Login Required)
```bash
# No login required - public operations
$ emm search --tags "parser,json"
$ emm info mymodule
$ emm release download mymodule@latest -s ./downloads
```

---

### 2. **Registered User (Basic User)**

Basic registered users have permissions to upload and manage their own modules and releases.

#### Characteristics:
- **Authentication Required**: Must have registered an account and logged in
- **Token-Based Access**: Holds a valid authentication token
- **Persistent Storage**: Credentials securely stored in system keyring
- **Own Module Management**: Can only manage modules they created
- **Release Suggestions**: Can suggest releases for their modules

#### Permissions:
- ✅ `CreateMyModule` - Upload and create new modules
- ✅ `SuggestMyModule` - Suggest new release versions
- ✅ `DeleteMyModule` - Delete own modules
- ✅ `DeleteMyRelease` - Delete own release suggestions
- ❌ All admin/supervision operations

#### User Operations:

**Authentication:**
- `register` - Create new account
- `login` - Authenticate with email/password
- `logout` - Disconnect from current session
- `whoami` - Display currently logged-in user

**Module Management (Own Only):**
- `module upload` - Upload a new module
- `module mylist` - View their own published modules
- `module delete` - Delete their own module

**Release Management (Own Only):**
- `module suggest` - Suggest a new release version for their module
- `release delete` - Delete their own suggested release
- `release cancel` - Cancel their own pending release

**Download:**
- `release download` - Download public releases (authenticated)
- `/api/auth/release/:release` - Download releases as authenticated user

#### Example: User Module Management
```bash
$ emm login --email user@example.com
# Password: [prompted]

$ emm module upload \
  --title "My Parser Module" \
  --repr "my-parser" \
  --tags "parser,json" \
  ./my-parser-v1.0.0.eva

$ emm module suggest --module my-parser --version 1.0.0

$ emm module mylist
```

---

### 3. **Maintainer**

Maintainers have elevated permissions to manage modules and releases across the system, but cannot perform user administration tasks.

#### Characteristics:
- **Authentication Required**: Must be assigned maintainer role
- **Full Module Management**: Can manage any module (not just their own)
- **Release Workflow**: Can accept/reject/cancel any release
- **Supervision Access**: Can filter and manage releases
- **No User Admin**: Cannot ban/unban users

#### Permissions:
- ✅ `CreateMyModule` - Create new modules
- ✅ `SuggestMyModule` - Suggest releases
- ✅ `DeleteModules` - Delete any module
- ✅ `DeleteMyModule` - Delete own modules
- ✅ `DeleteReleases` - Delete any release
- ✅ `DeleteMyRelease` - Delete own releases
- ✅ `UpdateReleases` - Update release information
- ✅ `ChangeReleaseStatuses` - Change release status
- ✅ `RejectReleases` - Reject suggested releases
- ✅ `AcceptReleases` - Accept suggested releases
- ✅ `CancelReleases` - Cancel pending releases
- ❌ `BanUsers` / `UnbanUsers` - Cannot manage users

#### Maintainer Operations:

**Extended Supervision:**
- `GET /api/supervise/release/:module/:version` - Find any release
- `GET /api/supervise/releases?filter=...` - Filter releases
- `POST /api/supervise/reject/:releaseId` - Reject release
- `POST /api/supervise/accept/:releaseId` - Accept release
- `POST /api/supervise/cancel/:releaseId` - Cancel release
- `POST /api/supervise/pending/:releaseId` - Change to pending status

#### Example: Maintainer Release Management
```bash
$ emm login --email maintainer@example.com

# Accept a user's suggested release
$ emm supervision release accept 12345

# Reject a problematic release
$ emm supervision release reject 12346

# Filter releases by status
$ emm supervision releases --filter "status=suggested"
```

---

### 4. **Administrator**

Admins have full permissions across the system, including user management and access to restricted releases.

#### Characteristics:
- **Authentication Required**: Must be assigned administrator role
- **Full System Access**: Can perform any operation
- **User Management**: Can ban/unban users
- **Restricted Downloads**: Can download any release regardless of status
- **Complete Supervision**: Full access to all supervision endpoints

#### Permissions:
- ✅ `CreateMyModule` - Create new modules
- ✅ `SuggestMyModule` - Suggest releases
- ✅ `DeleteModules` - Delete any module
- ✅ `DeleteMyModule` - Delete own modules
- ✅ `DeleteReleases` - Delete any release
- ✅ `DeleteMyRelease` - Delete own releases
- ✅ `UpdateReleases` - Update release information
- ✅ `ChangeReleaseStatuses` - Change release status
- ✅ `RejectReleases` - Reject suggested releases
- ✅ `AcceptReleases` - Accept suggested releases
- ✅ `CancelReleases` - Cancel pending releases
- ✅ `BanUsers` - Ban user accounts
- ✅ `UnbanUsers` - Restore banned user accounts

#### Administrator Operations:

**All Maintainer Operations + User Management:**
- `user list` - View all system users
- `POST /api/supervise/ban/:userId` - Ban a user account
- `POST /api/supervise/unban/:userId` - Restore banned user
- `GET /api/supervise/user/:email` - Get user information
- `GET /api/supervise/users` - List all users
- `GET /api/supervise/download/:releaseId` - Download any release (including restricted)

#### Example: Admin User Management
```bash
$ emm login --email admin@example.com

# Ban a problematic user
$ emm supervision user ban --email spammer@example.com

# View all users
$ emm supervision users list

# Download a restricted release
$ emm supervision download 99999
```

---

## Authorization Matrix by Role

The system implements four tiers of authorization, each with specific permissions:

| Operation | Public | User | Maintainer | Admin |
|-----------|--------|------|-----------|-------|
| **DISCOVERY & DOWNLOAD** | | | | |
| Search modules | ✅ | ✅ | ✅ | ✅ |
| View module info | ✅ | ✅ | ✅ | ✅ |
| Download public releases | ✅ | ✅ | ✅ | ✅ |
| Download restricted releases | ❌ | ❌ | ❌ | ✅ |
| **USER MODULE MANAGEMENT** | | | | |
| Create module | ❌ | ✅ | ✅ | ✅ |
| Upload module | ❌ | ✅ | ✅ | ✅ |
| Delete own module | ❌ | ✅ | ✅ | ✅ |
| Delete any module | ❌ | ❌ | ✅ | ✅ |
| **RELEASE MANAGEMENT** | | | | |
| Suggest release (own module) | ❌ | ✅ | ✅ | ✅ |
| Delete own release | ❌ | ✅ | ✅ | ✅ |
| Delete any release | ❌ | ❌ | ✅ | ✅ |
| Accept any release | ❌ | ❌ | ✅ | ✅ |
| Reject any release | ❌ | ❌ | ✅ | ✅ |
| Cancel any release | ❌ | ❌ | ✅ | ✅ |
| Change release status | ❌ | ❌ | ✅ | ✅ |
| Update release info | ❌ | ❌ | ✅ | ✅ |
| **USER ADMINISTRATION** | | | | |
| View all users | ❌ | ❌ | ❌ | ✅ |
| Ban user | ❌ | ❌ | ❌ | ✅ |
| Unban user | ❌ | ❌ | ❌ | ✅ |
| Get user info | ❌ | ❌ | ❌ | ✅ |
| **AUTHENTICATION** | | | | |
| Register account | ✅ | ✅ | ✅ | ✅ |
| Login | ✅ | ✅ | ✅ | ✅ |
| Logout | N/A | ✅ | ✅ | ✅ |

---

## Architecture Components

### 1. **Authentication Service** (`AuthService`)
- Manages user registration and login
- Maintains current user session
- Stores credentials securely via `KeyringService`
- Provides token management
- Handles user switching and logout

### 2. **Module Service** (`ModuleService`)
- Upload and management of module packages
- Module metadata updates
- User module listing
- Requires authentication token

### 3. **Release Service** (`ModuleReleaseService`)
- Manages module releases and versions
- Release status transitions (suggest → accept/reject)
- Release filtering and discovery
- Download management
- Supports both authenticated and public access

### 4. **Search Service** (`ModuleSearchService`)
- Full-text search across modules
- Filtering by tags, names, descriptions
- Public access (no authentication required)

### 5. **Installation Service** (`InstallService`)
- Local project file management (eva.json)
- Module installation/uninstallation
- Dependency resolution
- Both authenticated and public operations supported

### 6. **Supervision Service** (`SupervisionService`)
- Administrative user management
- User banning/unbanning
- Admin-only operations
- Requires admin privilege token

### 7. **Backend Service** (`Backend`)
- REST API communication with server
- HTTP request/response handling
- Server configuration management
- Base URL and endpoint management

---

## Data Flow Comparison

### Public User Flow (Search & Download)
```
Public User
    ↓
search/info command (NO auth needed)
    ↓
SearchService (reads public data)
    ↓
Backend API (empty token or public endpoint)
    ↓
Repository Server (public data)
    ↓
Display results
```

### Authenticated User Flow (Upload & Manage)
```
Authenticated User
    ↓
login (register credentials)
    ↓
AuthService (stores token securely)
    ↓
module upload/update command (with token)
    ↓
ModuleService (authenticated requests)
    ↓
Backend API (includes auth token)
    ↓
Repository Server (validates token, performs operation)
    ↓
Success/Error response
```

---

## Detailed Permission Mapping

### **Permission Definitions**

#### User Permissions (Basic Registered User)
```go
models.User: {
    CreateMyModule,      // Can upload new modules
    SuggestMyModule,     // Can suggest releases for own modules
    DeleteMyModule,      // Can delete their own modules
    DeleteMyRelease,     // Can delete their own release suggestions
}
```

#### Maintainer Permissions
```go
models.Maintainer: {
    CreateMyModule,         // Can create new modules
    SuggestMyModule,        // Can suggest releases
    DeleteModules,          // Can delete ANY module
    DeleteMyModule,         // Can delete own modules
    DeleteReleases,         // Can delete ANY release
    DeleteMyRelease,        // Can delete own releases
    UpdateReleases,         // Can update release metadata
    ChangeReleaseStatuses,  // Can change release status
    RejectReleases,         // Can reject suggested releases
    AcceptReleases,         // Can accept suggested releases
    CancelReleases,         // Can cancel pending releases
    // No user ban/unban permissions
}
```

#### Administrator Permissions
```go
models.Admin: {
    CreateMyModule,         // Can create new modules
    SuggestMyModule,        // Can suggest releases
    DeleteModules,          // Can delete ANY module
    DeleteMyModule,         // Can delete own modules
    DeleteReleases,         // Can delete ANY release
    DeleteMyRelease,        // Can delete own releases
    UpdateReleases,         // Can update release metadata
    ChangeReleaseStatuses,  // Can change release status
    RejectReleases,         // Can reject suggested releases
    AcceptReleases,         // Can accept suggested releases
    CancelReleases,         // Can cancel pending releases
    BanUsers,               // Can ban user accounts
    UnbanUsers,             // Can unban user accounts
    // Full system access
}
```

### **Access Control by Operation**

#### **Public Access (No Token Required)**
- `GET /api/modules/:id` - Find module by ID
- `GET /api/modules/search` - Search modules by components/tags
- `GET /api/releases/:id` - List all releases of a module
- `GET /api/releases/:id/release/:releaseId` - Get specific release details
- `GET /api/releases/:id/search` - Search releases by keywords
- `GET /api/releases/latest/:module` - Get latest release
- `GET /api/download/release/:release` - Download public releases

#### **User Access (Token + CreateMyModule + SuggestMyModule + DeleteMyModule + DeleteMyRelease)**
- `POST /api/modules/upload` - Upload new module
- `POST /api/modules/update` - Update own module
- `POST /api/modules/delete` - Delete own module
- `POST /api/releases/:id/delete/:releaseId` - Delete own release suggestion
- `POST /api/releases/:id/cancel/:releaseId` - Cancel own release suggestion
- `GET /api/modules/mylist` - View own modules
- `GET /api/auth/release/:release` - Download as authenticated user

#### **Maintainer Access (All User permissions + UpdateReleases + ChangeReleaseStatuses + RejectReleases + AcceptReleases + CancelReleases + DeleteModules + DeleteReleases)**
- `GET /api/supervise/release/:module/:version` - Find any release
- `GET /api/supervise/releases` - Filter all releases
- `POST /api/supervise/reject/:releaseId` - Reject any release
- `POST /api/supervise/accept/:releaseId` - Accept any release
- `POST /api/supervise/cancel/:releaseId` - Cancel any release
- `POST /api/supervise/pending/:releaseId` - Change to pending status
- All user-level operations on any module/release

#### **Admin Access (All Maintainer permissions + BanUsers + UnbanUsers)**
- `POST /api/supervise/ban/:userId` - Ban user account
- `POST /api/supervise/unban/:userId` - Unban user account
- `GET /api/supervise/user/:email` - Get user information
- `GET /api/supervise/users` - List all system users
- `GET /api/supervise/download/:releaseId` - Download any release (including restricted)
- All maintainer-level operations

---

## Command Decision Tree

```
┌─ User starts EMM CLI
│
├─ Search / View Info / Download (Public)? ──→ No auth needed ✅
│
├─ Install / Verify eva.json (Local only)? ──→ No auth needed ✅
│
├─ Upload / Update / Manage Module? ──→ Must login first
│   │
│   ├─ Logged in? ──→ Use token ✅
│   └─ Not logged in? ──→ `emm login` required ❌
│
├─ Admin operations? ──→ Admin token required
│   │
│   └─ Is admin? ──→ Proceed ✅
│
└─ Logout? ──→ Clear session
```

---

## Security Model

### Token Storage
- Authentication tokens stored securely in system **keyring** (OS password manager)
- Tokens retrieved transparently when needed
- Credentials never stored in plain text

### Session Management
- **Current User**: One active user per EMM installation
- **Multiple Accounts**: Can register multiple accounts and switch between them
- **Logout**: Removes active session but keeps credentials in keyring for re-login

### Authorization Flow
```
1. User runs command
2. Command checks if token required
3. If needed: GetCurrentUserToken() → Retrieve from keyring
4. If no token: Return error "User needs to be logged in first"
5. Pass token with API request
6. Server validates token and grants access
```

---

## Usage Examples by User Role

### Example 1: Public User - Search and Download

```bash
# No login required - public operations
$ emm search --tags "parser,json"
# Returns: list of modules with parser/json tags

$ emm info rss-parser
# Returns: module information and available releases

$ emm release download rss-parser@1.0.0 -s ./modules
# Downloads public release (anyone can download)
```

### Example 2: Basic User - Create & Manage Own Modules

```bash
# Must login first
$ emm login --email developer@example.com
# Password: [prompted]

# Now can upload own modules
$ emm module upload \
  --title "My Parser Module" \
  --repr "my-parser" \
  --tags "parser,json" \
  ./my-parser-v1.0.0.eva

# Suggest release for own module
$ emm module suggest \
  --module-id my-parser \
  --version 1.0.0

# View own modules
$ emm module mylist

# Delete own module
$ emm module delete --module-id my-parser

# Delete own release suggestion
$ emm release delete --release-id 12345
```

### Example 3: Maintainer - Manage All Releases

```bash
# Login as maintainer
$ emm login --email maintainer@example.com

# View all releases (with filtering capability)
$ emm release dump

# Accept a user's suggested release (promote to production)
$ emm release accept --release-id 12345
# Module now available to all users

# Reject a problematic release
$ emm release reject --release-id 12346
# Module author gets notification

# Cancel a pending release
$ emm release cancel --release-id 12347

# Download and inspect a specific release
$ emm release download rss-parser@1.0.0 -s ./inspect
```

### Example 4: Administrator - Full System Management

```bash
# Login as admin
$ emm login --email admin@example.com

# Perform all maintainer release operations
$ emm release accept --release-id 12345      # Accept releases
$ emm release reject --release-id 12346      # Reject releases
$ emm release cancel --release-id 12347      # Cancel releases

# Additionally, manage users
$ emm user list
# Output: All registered users in system

# Ban problematic user
$ emm user ban --user-id 999
# User cannot upload/suggest releases anymore

# Unban user when issue is resolved
$ emm user unban --user-id 999

# Download any release (including restricted ones)
$ emm release download restricted-module@1.0.0 -s ./modules
# Works even if release has restricted status
```

### Example 5: User Registration

```bash
# Register a new account (before login)
$ emm register \
  --email developer@example.com \
  --handle "dev-user" \
  --first-name "John" \
  --last-name "Doe"
# Password: [prompted securely]

# Account created and can now login
$ emm login --email developer@example.com
# Password: [prompted]

# Confirm logged-in user
$ emm whoami
# Output: john.doe@example.com (dev-user)
```

### Example 6: Role Transition

```bash
# User starts as basic User
$ emm login --email developer@example.com
# Can: upload own modules, suggest own releases

# User is promoted to Maintainer (by admin)
# Same login, now additional capabilities:
$ emm release dump
# Can now: accept/reject ANY release, manage any module

# User is promoted to Administrator
# Same login, now full capabilities:
$ emm user list
$ emm user ban --user-id 123
# Can now: ban/unban users, download restricted releases
```

---

## Summary Table: Role Capabilities

| Aspect | Public | User | Maintainer | Admin |
|--------|--------|------|-----------|-------|
| **Login Required** | No | Yes | Yes | Yes |
| **Token Storage** | N/A | Keyring | Keyring | Keyring |
| **Search Modules** | ✅ Full | ✅ Full | ✅ Full | ✅ Full |
| **Download Public Releases** | ✅ | ✅ | ✅ | ✅ |
| **Download Restricted Releases** | ❌ | ❌ | ❌ | ✅ |
| **Upload Modules** | ❌ | ✅ Own | ✅ Any | ✅ Any |
| **Manage Own Modules** | ❌ | ✅ | ✅ | ✅ |
| **Manage Any Module** | ❌ | ❌ | ✅ | ✅ |
| **Suggest Releases** | ❌ | ✅ Own | ✅ Any | ✅ Any |
| **Accept Releases** | ❌ | ❌ | ✅ | ✅ |
| **Reject Releases** | ❌ | ❌ | ✅ | ✅ |
| **Cancel Releases** | ❌ | ❌ | ✅ | ✅ |
| **Manage Users** | ❌ | ❌ | ❌ | ✅ |
| **Ban/Unban Users** | ❌ | ❌ | ❌ | ✅ |
| **Filter Releases** | ❌ | ❌ | ✅ | ✅ |
| **Session Switching** | N/A | ✅ | ✅ | ✅ |
| **Data Persistence** | None | Full | Full | Full |

---

## Key Insights

1. **Public by Default**: Search and download are public operations requiring no authentication
2. **Token-Optional**: Some endpoints work with or without tokens, providing different data scopes
3. **Secure Credentials**: Tokens stored in OS keyring, not in files
4. **Session Switching**: Power users can manage multiple registered accounts
5. **Admin Segregation**: Administrative operations require explicit admin privileges
6. **Offline Capable**: Local operations (verify, install/uninstall) work without server connection
