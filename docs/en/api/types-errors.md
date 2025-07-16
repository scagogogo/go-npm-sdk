# Types & Errors API

This document covers all data types, structures, and error handling mechanisms used throughout the Go NPM SDK.

## Core Types

### Client Interface

The main interface for npm operations:

```go
type Client interface {
    IsAvailable(ctx context.Context) bool
    Install(ctx context.Context) error
    Version(ctx context.Context) (string, error)
    Init(ctx context.Context, options InitOptions) error
    InstallPackage(ctx context.Context, pkg string, options InstallOptions) error
    UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error
    UpdatePackage(ctx context.Context, pkg string) error
    ListPackages(ctx context.Context, options ListOptions) ([]Package, error)
    RunScript(ctx context.Context, script string, args ...string) error
    Publish(ctx context.Context, options PublishOptions) error
    GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)
    Search(ctx context.Context, query string) ([]SearchResult, error)
}
```

## Option Types

### InitOptions

Configuration for project initialization:

```go
type InitOptions struct {
    Name        string `json:"name,omitempty"`
    Version     string `json:"version,omitempty"`
    Description string `json:"description,omitempty"`
    Author      string `json:"author,omitempty"`
    License     string `json:"license,omitempty"`
    Private     bool   `json:"private,omitempty"`
    WorkingDir  string `json:"-"`
    Force       bool   `json:"-"`
}
```

### InstallOptions

Configuration for package installation:

```go
type InstallOptions struct {
    SaveDev       bool              `json:"save_dev,omitempty"`
    SaveOptional  bool              `json:"save_optional,omitempty"`
    SaveExact     bool              `json:"save_exact,omitempty"`
    Global        bool              `json:"global,omitempty"`
    Production    bool              `json:"production,omitempty"`
    WorkingDir    string            `json:"working_dir,omitempty"`
    Registry      string            `json:"registry,omitempty"`
    Force         bool              `json:"force,omitempty"`
    IgnoreScripts bool              `json:"ignore_scripts,omitempty"`
}
```

### UninstallOptions

Configuration for package uninstallation:

```go
type UninstallOptions struct {
    SaveDev    bool   `json:"save_dev,omitempty"`
    Global     bool   `json:"global,omitempty"`
    WorkingDir string `json:"working_dir,omitempty"`
}
```

### ListOptions

Configuration for package listing:

```go
type ListOptions struct {
    Global     bool   `json:"global,omitempty"`
    Depth      int    `json:"depth,omitempty"`
    Production bool   `json:"production,omitempty"`
    WorkingDir string `json:"working_dir,omitempty"`
    JSON       bool   `json:"json,omitempty"`
}
```

### PublishOptions

Configuration for package publishing:

```go
type PublishOptions struct {
    Tag        string `json:"tag,omitempty"`
    Access     string `json:"access,omitempty"`
    Registry   string `json:"registry,omitempty"`
    WorkingDir string `json:"working_dir,omitempty"`
    DryRun     bool   `json:"dry_run,omitempty"`
}
```

## Data Types

### Package

Represents an npm package:

```go
type Package struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description,omitempty"`
    Dependencies map[string]string `json:"dependencies,omitempty"`
    DevDeps      map[string]string `json:"devDependencies,omitempty"`
    OptionalDeps map[string]string `json:"optionalDependencies,omitempty"`
    PeerDeps     map[string]string `json:"peerDependencies,omitempty"`
    Scripts      map[string]string `json:"scripts,omitempty"`
    Keywords     []string          `json:"keywords,omitempty"`
    Author       string            `json:"author,omitempty"`
    License      string            `json:"license,omitempty"`
    Homepage     string            `json:"homepage,omitempty"`
    Repository   *Repository       `json:"repository,omitempty"`
    Bugs         *Bugs             `json:"bugs,omitempty"`
    Main         string            `json:"main,omitempty"`
    Private      bool              `json:"private,omitempty"`
}
```

### Repository

Repository information:

```go
type Repository struct {
    Type string `json:"type"`
    URL  string `json:"url"`
}
```

### Bugs

Bug reporting information:

```go
type Bugs struct {
    URL   string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}
```

### PackageInfo

Detailed package information from registry:

```go
type PackageInfo struct {
    Name         string                 `json:"name"`
    Version      string                 `json:"version"`
    Description  string                 `json:"description,omitempty"`
    Keywords     []string               `json:"keywords,omitempty"`
    Homepage     string                 `json:"homepage,omitempty"`
    Repository   *Repository            `json:"repository,omitempty"`
    Author       *Person                `json:"author,omitempty"`
    License      string                 `json:"license,omitempty"`
    Dependencies map[string]string      `json:"dependencies,omitempty"`
    DevDeps      map[string]string      `json:"devDependencies,omitempty"`
    Versions     map[string]interface{} `json:"versions,omitempty"`
    Time         map[string]time.Time   `json:"time,omitempty"`
    DistTags     map[string]string      `json:"dist-tags,omitempty"`
}
```

### Person

Person information (author, maintainer, etc.):

```go
type Person struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    URL   string `json:"url,omitempty"`
}
```

### SearchResult

Search result from npm registry:

```go
type SearchResult struct {
    Package     SearchPackage `json:"package"`
    Score       SearchScore   `json:"score"`
    SearchScore float64       `json:"searchScore"`
}
```

### SearchPackage

Package information in search results:

```go
type SearchPackage struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Keywords    []string          `json:"keywords,omitempty"`
    Date        time.Time         `json:"date"`
    Links       map[string]string `json:"links,omitempty"`
    Author      *Person           `json:"author,omitempty"`
    Publisher   *Person           `json:"publisher,omitempty"`
    Maintainers []*Person         `json:"maintainers,omitempty"`
}
```

### SearchScore

Scoring information for search results:

```go
type SearchScore struct {
    Final  float64     `json:"final"`
    Detail ScoreDetail `json:"detail"`
}

type ScoreDetail struct {
    Quality     float64 `json:"quality"`
    Popularity  float64 `json:"popularity"`
    Maintenance float64 `json:"maintenance"`
}
```

### CommandResult

Result of command execution:

```go
type CommandResult struct {
    Success  bool          `json:"success"`
    ExitCode int           `json:"exit_code"`
    Stdout   string        `json:"stdout,omitempty"`
    Stderr   string        `json:"stderr,omitempty"`
    Duration time.Duration `json:"duration"`
    Error    error         `json:"-"`
}
```

## Installation Types

### NpmInstallOptions

Options for npm installation:

```go
type NpmInstallOptions struct {
    Method      string                `json:"method"`
    Version     string                `json:"version,omitempty"`
    InstallPath string                `json:"install_path,omitempty"`
    Force       bool                  `json:"force,omitempty"`
    Global      bool                  `json:"global,omitempty"`
    Progress    func(message string)  `json:"-"`
}
```

### InstallResult

Result of npm installation:

```go
type InstallResult struct {
    Success  bool          `json:"success"`
    Method   string        `json:"method"`
    Version  string        `json:"version"`
    Path     string        `json:"path"`
    Duration time.Duration `json:"duration"`
    Error    string        `json:"error,omitempty"`
}
```

### NpmInfo

Information about npm installation:

```go
type NpmInfo struct {
    Version   string `json:"version"`
    Path      string `json:"path"`
    Available bool   `json:"available"`
}
```

### PortableConfig

Configuration for portable npm installation:

```go
type PortableConfig struct {
    Version     string `json:"version"`
    InstallPath string `json:"install_path"`
    NodePath    string `json:"node_path"`
    NpmPath     string `json:"npm_path"`
    InstallDate string `json:"install_date"`
}
```

## Error Types

### NpmError

Base error type for npm operations:

```go
type NpmError struct {
    Op       string `json:"op"`
    Package  string `json:"package,omitempty"`
    ExitCode int    `json:"exit_code"`
    Stdout   string `json:"stdout,omitempty"`
    Stderr   string `json:"stderr,omitempty"`
    Err      error  `json:"-"`
}

func (e *NpmError) Error() string
func (e *NpmError) Unwrap() error
```

### ValidationError

Error for validation failures:

```go
type ValidationError struct {
    Field  string `json:"field"`
    Value  string `json:"value"`
    Reason string `json:"reason"`
}

func (e *ValidationError) Error() string
```

### InstallError

Error for installation failures:

```go
type InstallError struct {
    Package string `json:"package"`
    Reason  string `json:"reason"`
    Err     error  `json:"-"`
}

func (e *InstallError) Error() string
func (e *InstallError) Unwrap() error
```

### UninstallError

Error for uninstallation failures:

```go
type UninstallError struct {
    Package string `json:"package"`
    Reason  string `json:"reason"`
    Err     error  `json:"-"`
}

func (e *UninstallError) Error() string
func (e *UninstallError) Unwrap() error
```

### PlatformError

Error for platform-related issues:

```go
type PlatformError struct {
    Platform string `json:"platform"`
    Reason   string `json:"reason"`
    Err      error  `json:"-"`
}

func (e *PlatformError) Error() string
func (e *PlatformError) Unwrap() error
```

### DownloadError

Error for download failures:

```go
type DownloadError struct {
    URL    string `json:"url"`
    Reason string `json:"reason"`
    Err    error  `json:"-"`
}

func (e *DownloadError) Error() string
func (e *DownloadError) Unwrap() error
```

## Error Constants

Predefined error constants:

```go
var (
    ErrNpmNotFound         = errors.New("npm not found")
    ErrPackageNotFound     = errors.New("package not found")
    ErrInvalidVersion      = errors.New("invalid version")
    ErrInvalidPackageName  = errors.New("invalid package name")
    ErrPermissionDenied    = errors.New("permission denied")
    ErrUnsupportedPlatform = errors.New("unsupported platform")
)
```

## Error Checking Functions

Helper functions for error type checking:

```go
func IsNpmNotFound(err error) bool
func IsPackageNotFound(err error) bool
func IsDownloadError(err error) bool
func IsPlatformError(err error) bool
```

**Example usage:**
```go
err := client.InstallPackage(ctx, "nonexistent-package", npm.InstallOptions{})
if err != nil {
    if npm.IsPackageNotFound(err) {
        fmt.Println("Package not found in registry")
    } else if npm.IsNpmNotFound(err) {
        fmt.Println("npm is not installed")
    } else if npmErr, ok := err.(*npm.NpmError); ok {
        fmt.Printf("npm command failed: %s\n", npmErr.Stderr)
    } else {
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

## Constructor Functions

Helper functions for creating error types:

```go
func NewNpmError(op, pkg string, exitCode int, stdout, stderr string, err error) *NpmError
func NewValidationError(field, value, reason string) *ValidationError
func NewInstallError(pkg, reason string, err error) *InstallError
func NewUninstallError(pkg, reason string, err error) *UninstallError
func NewPlatformError(platform, reason string, err error) *PlatformError
func NewDownloadError(url, reason string, err error) *DownloadError
```

## Best Practices

1. **Error Handling**: Always check for specific error types using the provided helper functions
2. **Validation**: Use ValidationError for input validation failures
3. **Context**: Include relevant context information in error messages
4. **Wrapping**: Use error wrapping to preserve the original error chain
5. **Structured Data**: Use structured error types instead of plain strings for better error handling
6. **JSON Serialization**: Most types support JSON serialization for API responses
7. **Nil Checks**: Always check for nil pointers when working with optional fields
