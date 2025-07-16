# Client Interface

The `Client` interface is the main entry point for all npm operations in the Go NPM SDK. It provides a comprehensive set of methods for managing npm installations, packages, and projects.

## Creating a Client

### NewClient

Creates a new npm client with default configuration.

```go
func NewClient() (Client, error)
```

**Example:**
```go
client, err := npm.NewClient()
if err != nil {
    log.Fatal(err)
}
```

### NewClientWithPath

Creates a new npm client with a specific npm executable path.

```go
func NewClientWithPath(npmPath string) (Client, error)
```

**Parameters:**
- `npmPath` (string): Path to the npm executable

**Example:**
```go
client, err := npm.NewClientWithPath("/usr/local/bin/npm")
if err != nil {
    log.Fatal(err)
}
```

## Basic Operations

### IsAvailable

Checks if npm is available and can be executed.

```go
IsAvailable(ctx context.Context) bool
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout

**Returns:**
- `bool`: true if npm is available, false otherwise

**Example:**
```go
ctx := context.Background()
if client.IsAvailable(ctx) {
    fmt.Println("npm is available")
} else {
    fmt.Println("npm is not available")
}
```

### Install

Automatically installs npm if it's not available on the system.

```go
Install(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout

**Returns:**
- `error`: Error if installation fails

**Example:**
```go
ctx := context.Background()
if !client.IsAvailable(ctx) {
    err := client.Install(ctx)
    if err != nil {
        log.Fatalf("Failed to install npm: %v", err)
    }
    fmt.Println("npm installed successfully")
}
```

### Version

Gets the current npm version.

```go
Version(ctx context.Context) (string, error)
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout

**Returns:**
- `string`: npm version string
- `error`: Error if version retrieval fails

**Example:**
```go
ctx := context.Background()
version, err := client.Version(ctx)
if err != nil {
    log.Fatalf("Failed to get npm version: %v", err)
}
fmt.Printf("npm version: %s\n", version)
```

## Project Management

### Init

Initializes a new npm project with package.json.

```go
Init(ctx context.Context, options InitOptions) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `options` (InitOptions): Initialization options

**Returns:**
- `error`: Error if initialization fails

**Example:**
```go
ctx := context.Background()
options := npm.InitOptions{
    Name:        "my-project",
    Version:     "1.0.0",
    Description: "My awesome project",
    Author:      "John Doe",
    License:     "MIT",
    WorkingDir:  "/path/to/project",
}

err := client.Init(ctx, options)
if err != nil {
    log.Fatalf("Failed to initialize project: %v", err)
}
```

## Package Operations

### InstallPackage

Installs a specific npm package.

```go
InstallPackage(ctx context.Context, pkg string, options InstallOptions) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `pkg` (string): Package name to install
- `options` (InstallOptions): Installation options

**Returns:**
- `error`: Error if installation fails

**Example:**
```go
ctx := context.Background()
options := npm.InstallOptions{
    SaveDev:    true,
    SaveExact:  true,
    WorkingDir: "/path/to/project",
}

err := client.InstallPackage(ctx, "typescript", options)
if err != nil {
    log.Fatalf("Failed to install package: %v", err)
}
```

### UninstallPackage

Uninstalls a specific npm package.

```go
UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `pkg` (string): Package name to uninstall
- `options` (UninstallOptions): Uninstallation options

**Returns:**
- `error`: Error if uninstallation fails

**Example:**
```go
ctx := context.Background()
options := npm.UninstallOptions{
    SaveDev:    true,
    WorkingDir: "/path/to/project",
}

err := client.UninstallPackage(ctx, "typescript", options)
if err != nil {
    log.Fatalf("Failed to uninstall package: %v", err)
}
```

### UpdatePackage

Updates a specific npm package to the latest version.

```go
UpdatePackage(ctx context.Context, pkg string) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `pkg` (string): Package name to update

**Returns:**
- `error`: Error if update fails

**Example:**
```go
ctx := context.Background()
err := client.UpdatePackage(ctx, "lodash")
if err != nil {
    log.Fatalf("Failed to update package: %v", err)
}
```

### ListPackages

Lists installed packages in the current project.

```go
ListPackages(ctx context.Context, options ListOptions) ([]Package, error)
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `options` (ListOptions): Listing options

**Returns:**
- `[]Package`: List of installed packages
- `error`: Error if listing fails

**Example:**
```go
ctx := context.Background()
options := npm.ListOptions{
    Global:     false,
    Depth:      1,
    WorkingDir: "/path/to/project",
    JSON:       true,
}

packages, err := client.ListPackages(ctx, options)
if err != nil {
    log.Fatalf("Failed to list packages: %v", err)
}

for _, pkg := range packages {
    fmt.Printf("Package: %s@%s\n", pkg.Name, pkg.Version)
}
```

## Script Execution

### RunScript

Executes an npm script defined in package.json.

```go
RunScript(ctx context.Context, script string, args ...string) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `script` (string): Script name to execute
- `args` (...string): Additional arguments to pass to the script

**Returns:**
- `error`: Error if script execution fails

**Example:**
```go
ctx := context.Background()

// Run build script
err := client.RunScript(ctx, "build")
if err != nil {
    log.Fatalf("Failed to run build script: %v", err)
}

// Run test script with arguments
err = client.RunScript(ctx, "test", "--verbose", "--coverage")
if err != nil {
    log.Fatalf("Failed to run test script: %v", err)
}
```

## Publishing

### Publish

Publishes the current package to npm registry.

```go
Publish(ctx context.Context, options PublishOptions) error
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `options` (PublishOptions): Publishing options

**Returns:**
- `error`: Error if publishing fails

**Example:**
```go
ctx := context.Background()
options := npm.PublishOptions{
    Tag:        "latest",
    Access:     "public",
    Registry:   "https://registry.npmjs.org/",
    WorkingDir: "/path/to/project",
    DryRun:     false,
}

err := client.Publish(ctx, options)
if err != nil {
    log.Fatalf("Failed to publish package: %v", err)
}
```

## Information Retrieval

### GetPackageInfo

Retrieves detailed information about a package from the npm registry.

```go
GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `pkg` (string): Package name to get information for

**Returns:**
- `*PackageInfo`: Detailed package information
- `error`: Error if information retrieval fails

**Example:**
```go
ctx := context.Background()
info, err := client.GetPackageInfo(ctx, "lodash")
if err != nil {
    log.Fatalf("Failed to get package info: %v", err)
}

fmt.Printf("Package: %s@%s\n", info.Name, info.Version)
fmt.Printf("Description: %s\n", info.Description)
fmt.Printf("Homepage: %s\n", info.Homepage)
```

### Search

Searches for packages in the npm registry.

```go
Search(ctx context.Context, query string) ([]SearchResult, error)
```

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `query` (string): Search query

**Returns:**
- `[]SearchResult`: List of search results
- `error`: Error if search fails

**Example:**
```go
ctx := context.Background()
results, err := client.Search(ctx, "react hooks")
if err != nil {
    log.Fatalf("Failed to search packages: %v", err)
}

for _, result := range results {
    fmt.Printf("Package: %s@%s\n", result.Package.Name, result.Package.Version)
    fmt.Printf("Description: %s\n", result.Package.Description)
    fmt.Printf("Score: %.2f\n", result.Score.Final)
    fmt.Println("---")
}
```

## Error Handling

The client methods return structured errors that can be checked for specific conditions:

```go
err := client.InstallPackage(ctx, "nonexistent-package", npm.InstallOptions{})
if err != nil {
    if npm.IsPackageNotFound(err) {
        fmt.Println("Package not found")
    } else if npm.IsNpmNotFound(err) {
        fmt.Println("npm is not installed")
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Best Practices

1. **Always use context**: Pass appropriate context for timeout and cancellation control
2. **Check availability**: Use `IsAvailable()` before performing operations
3. **Handle errors**: Check for specific error types for better error handling
4. **Use options**: Configure operations using option structs for flexibility
5. **Working directory**: Set appropriate working directory in options for project-specific operations
