# API Overview

Go NPM SDK provides a comprehensive set of APIs for npm operations in Go applications. The SDK is organized into several packages, each serving specific purposes.

## Package Structure

```
github.com/scagogogo/go-npm-sdk/
├── pkg/npm/           # Core npm operations
├── pkg/platform/      # Platform detection and downloads
└── pkg/utils/         # Utility functions
```

## Core Packages

### pkg/npm

The main package containing npm client implementation and related functionality:

- **Client Interface**: Main npm client for all operations
- **Installer**: Automatic npm installation management
- **Detector**: npm availability detection
- **Portable Manager**: Portable npm version management
- **Package Manager**: package.json file operations
- **Dependency Manager**: Dependency resolution and management
- **Types**: Data structures and interfaces
- **Errors**: Error types and handling

### pkg/platform

Platform-specific functionality:

- **Detector**: Operating system and architecture detection
- **Downloader**: File download capabilities with progress tracking

### pkg/utils

Utility functions:

- **Executor**: Command execution with advanced features

## Key Interfaces

### Client Interface

The main interface for npm operations:

```go
type Client interface {
    // Basic operations
    IsAvailable(ctx context.Context) bool
    Install(ctx context.Context) error
    Version(ctx context.Context) (string, error)
    
    // Project management
    Init(ctx context.Context, options InitOptions) error
    
    // Package operations
    InstallPackage(ctx context.Context, pkg string, options InstallOptions) error
    UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error
    UpdatePackage(ctx context.Context, pkg string) error
    ListPackages(ctx context.Context, options ListOptions) ([]Package, error)
    
    // Script execution
    RunScript(ctx context.Context, script string, args ...string) error
    
    // Publishing
    Publish(ctx context.Context, options PublishOptions) error
    
    // Information retrieval
    GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)
    Search(ctx context.Context, query string) ([]SearchResult, error)
}
```

## Common Patterns

### Context Usage

All operations accept a `context.Context` parameter for cancellation and timeout control:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.InstallPackage(ctx, "lodash", npm.InstallOptions{})
```

### Error Handling

The SDK provides structured error types for better error handling:

```go
if err != nil {
    if npm.IsNpmNotFound(err) {
        // Handle npm not found
    } else if npm.IsPackageNotFound(err) {
        // Handle package not found
    } else {
        // Handle other errors
    }
}
```

### Options Pattern

Most operations use options structs for configuration:

```go
options := npm.InstallOptions{
    SaveDev:    true,
    SaveExact:  true,
    WorkingDir: "/path/to/project",
    Registry:   "https://registry.npmjs.org/",
}

err := client.InstallPackage(ctx, "typescript", options)
```

## Getting Started

1. **Import the package**:
   ```go
   import "github.com/scagogogo/go-npm-sdk/pkg/npm"
   ```

2. **Create a client**:
   ```go
   client, err := npm.NewClient()
   if err != nil {
       log.Fatal(err)
   }
   ```

3. **Use the client**:
   ```go
   ctx := context.Background()
   
   if !client.IsAvailable(ctx) {
       err := client.Install(ctx)
       if err != nil {
           log.Fatal(err)
       }
   }
   
   version, err := client.Version(ctx)
   if err != nil {
       log.Fatal(err)
   }
   fmt.Printf("npm version: %s\n", version)
   ```

## Next Steps

- [Client Interface](./client.md) - Detailed client API documentation
- [NPM Package](./npm.md) - Complete npm package reference
- [Platform Package](./platform.md) - Platform detection and download APIs
- [Utils Package](./utils.md) - Utility functions reference
- [Types & Errors](./types-errors.md) - Data types and error handling
