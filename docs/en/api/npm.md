# NPM Package API

The `pkg/npm` package provides comprehensive npm functionality including client operations, installation management, package.json handling, and dependency resolution.

## Core Components

### Client

The main npm client implementation.

#### NewClient

```go
func NewClient() (Client, error)
```

Creates a new npm client with automatic npm detection.

#### NewClientWithPath

```go
func NewClientWithPath(npmPath string) (Client, error)
```

Creates a new npm client with a specific npm executable path.

### Installer

Manages npm installation across different platforms and methods.

#### NewInstaller

```go
func NewInstaller() (*Installer, error)
```

Creates a new npm installer with platform detection.

**Example:**
```go
installer, err := npm.NewInstaller()
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
options := npm.NpmInstallOptions{
    Method:  npm.PackageManager, // or npm.OfficialInstaller, npm.Portable
    Version: "latest",
    Force:   false,
}

result, err := installer.Install(ctx, options)
if err != nil {
    log.Fatalf("Installation failed: %v", err)
}

fmt.Printf("npm installed successfully: %s\n", result.Version)
```

#### Installation Methods

The installer supports multiple installation methods:

- `PackageManager`: Use system package manager (apt, yum, brew, etc.)
- `OfficialInstaller`: Download and run official Node.js installer
- `Portable`: Download portable Node.js/npm version
- `Manual`: Manual installation guidance

### Detector

Detects npm availability and version information.

#### NewDetector

```go
func NewDetector() *Detector
```

Creates a new npm detector.

#### IsAvailable

```go
func (d *Detector) IsAvailable(ctx context.Context) bool
```

Checks if npm is available on the system.

#### Detect

```go
func (d *Detector) Detect(ctx context.Context) (*NpmInfo, error)
```

Detects npm installation details.

**Example:**
```go
detector := npm.NewDetector()
ctx := context.Background()

if detector.IsAvailable(ctx) {
    info, err := detector.Detect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("npm found at: %s\n", info.Path)
    fmt.Printf("npm version: %s\n", info.Version)
} else {
    fmt.Println("npm not found")
}
```

### Portable Manager

Manages portable npm installations and versions.

#### NewPortableManager

```go
func NewPortableManager(baseDir string) (*PortableManager, error)
```

Creates a new portable manager with the specified base directory.

#### Install

```go
func (pm *PortableManager) Install(ctx context.Context, version string) (*PortableConfig, error)
```

Installs a portable npm version.

#### List

```go
func (pm *PortableManager) List() ([]*PortableConfig, error)
```

Lists all installed portable versions.

#### Uninstall

```go
func (pm *PortableManager) Uninstall(version string) error
```

Uninstalls a portable npm version.

#### CreateClient

```go
func (pm *PortableManager) CreateClient(version string) (Client, error)
```

Creates a client for a specific portable version.

**Example:**
```go
manager, err := npm.NewPortableManager("/opt/npm-portable")
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

// Install Node.js 18.17.0 with npm
config, err := manager.Install(ctx, "18.17.0")
if err != nil {
    log.Fatal(err)
}

// Create client for this version
client, err := manager.CreateClient("18.17.0")
if err != nil {
    log.Fatal(err)
}

// Use the client
version, err := client.Version(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Portable npm version: %s\n", version)
```

### Package Manager

Handles package.json file operations.

#### NewPackageJSON

```go
func NewPackageJSON(filePath string) *PackageJSON
```

Creates a new package.json manager for the specified file.

#### Load

```go
func (p *PackageJSON) Load() error
```

Loads package.json from disk.

#### Save

```go
func (p *PackageJSON) Save() error
```

Saves package.json to disk.

#### Basic Operations

```go
// Getters
func (p *PackageJSON) GetName() string
func (p *PackageJSON) GetVersion() string
func (p *PackageJSON) GetDescription() string
func (p *PackageJSON) GetAuthor() string
func (p *PackageJSON) GetLicense() string

// Setters
func (p *PackageJSON) SetName(name string)
func (p *PackageJSON) SetVersion(version string)
func (p *PackageJSON) SetDescription(description string)
func (p *PackageJSON) SetAuthor(author string)
func (p *PackageJSON) SetLicense(license string)
```

#### Dependency Management

```go
// Dependencies
func (p *PackageJSON) AddDependency(name, version string)
func (p *PackageJSON) RemoveDependency(name string)
func (p *PackageJSON) GetDependencies() map[string]string

// Dev Dependencies
func (p *PackageJSON) AddDevDependency(name, version string)
func (p *PackageJSON) RemoveDevDependency(name string)
func (p *PackageJSON) GetDevDependencies() map[string]string

// Peer Dependencies
func (p *PackageJSON) AddPeerDependency(name, version string)
func (p *PackageJSON) RemovePeerDependency(name string)
func (p *PackageJSON) GetPeerDependencies() map[string]string

// Optional Dependencies
func (p *PackageJSON) AddOptionalDependency(name, version string)
func (p *PackageJSON) RemoveOptionalDependency(name string)
func (p *PackageJSON) GetOptionalDependencies() map[string]string
```

#### Script Management

```go
func (p *PackageJSON) AddScript(name, command string)
func (p *PackageJSON) RemoveScript(name string)
func (p *PackageJSON) GetScripts() map[string]string
func (p *PackageJSON) HasScript(name string) bool
```

#### Repository and Metadata

```go
func (p *PackageJSON) SetRepository(repo *Repository)
func (p *PackageJSON) GetRepository() *Repository
func (p *PackageJSON) SetRepositoryURL(url string)

func (p *PackageJSON) SetBugs(bugs *Bugs)
func (p *PackageJSON) GetBugs() *Bugs
func (p *PackageJSON) SetBugsURL(url string)

func (p *PackageJSON) SetHomepage(homepage string)
func (p *PackageJSON) GetHomepage() string

func (p *PackageJSON) SetKeywords(keywords []string)
func (p *PackageJSON) GetKeywords() []string
func (p *PackageJSON) AddKeyword(keyword string)
```

#### Validation

```go
func (p *PackageJSON) Validate() error
```

Validates the package.json structure and content.

**Example:**
```go
pkg := npm.NewPackageJSON("./package.json")

// Load existing package.json
if err := pkg.Load(); err != nil {
    log.Fatal(err)
}

// Modify package information
pkg.SetName("my-awesome-package")
pkg.SetVersion("1.0.0")
pkg.SetDescription("An awesome package")

// Add dependencies
pkg.AddDependency("lodash", "^4.17.21")
pkg.AddDevDependency("typescript", "^4.5.0")

// Add scripts
pkg.AddScript("build", "tsc")
pkg.AddScript("test", "jest")

// Set repository
repo := &npm.Repository{
    Type: "git",
    URL:  "https://github.com/user/repo.git",
}
pkg.SetRepository(repo)

// Validate and save
if err := pkg.Validate(); err != nil {
    log.Fatal(err)
}

if err := pkg.Save(); err != nil {
    log.Fatal(err)
}
```

### Dependency Manager

Handles dependency resolution and conflict detection.

#### NewDependencyManager

```go
func NewDependencyManager() *DependencyManager
```

Creates a new dependency manager.

#### ResolveDependencies

```go
func (dm *DependencyManager) ResolveDependencies(ctx context.Context, deps map[string]string) (*DependencyTree, error)
```

Resolves dependencies and builds a dependency tree.

#### CheckConflicts

```go
func (dm *DependencyManager) CheckConflicts(tree *DependencyTree) []Conflict
```

Checks for dependency conflicts in the tree.

#### DetectCircularDependencies

```go
func (dm *DependencyManager) DetectCircularDependencies(tree *DependencyTree) [][]string
```

Detects circular dependencies in the tree.

**Example:**
```go
manager := npm.NewDependencyManager()
ctx := context.Background()

dependencies := map[string]string{
    "lodash": "^4.17.21",
    "axios":  "^0.24.0",
    "react":  "^17.0.0",
}

// Resolve dependencies
tree, err := manager.ResolveDependencies(ctx, dependencies)
if err != nil {
    log.Fatal(err)
}

// Check for conflicts
conflicts := manager.CheckConflicts(tree)
if len(conflicts) > 0 {
    fmt.Println("Dependency conflicts found:")
    for _, conflict := range conflicts {
        fmt.Printf("- %s: %s vs %s\n", conflict.Package, conflict.Version1, conflict.Version2)
    }
}

// Check for circular dependencies
circular := manager.DetectCircularDependencies(tree)
if len(circular) > 0 {
    fmt.Println("Circular dependencies found:")
    for _, cycle := range circular {
        fmt.Printf("- %s\n", strings.Join(cycle, " -> "))
    }
}
```

## Constants

### Installation Methods

```go
const (
    PackageManager    = "package_manager"
    OfficialInstaller = "official_installer"
    Portable          = "portable"
    Manual            = "manual"
)
```

### Error Constants

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

```go
func IsNpmNotFound(err error) bool
func IsPackageNotFound(err error) bool
func IsNetworkError(err error) bool
func IsPlatformError(err error) bool
```

These functions help identify specific error conditions for better error handling.

## Best Practices

1. **Use appropriate managers**: Choose the right component for your use case
2. **Handle errors gracefully**: Use error checking functions for specific error handling
3. **Validate inputs**: Always validate package names and versions
4. **Use contexts**: Pass contexts for timeout and cancellation control
5. **Manage portable versions**: Use portable manager for isolated npm environments
6. **Backup package.json**: Always backup before modifying package.json files
