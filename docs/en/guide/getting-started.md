# Getting Started

Welcome to Go NPM SDK! This guide will help you get up and running with the SDK quickly.

## Installation

Install the SDK using Go modules:

```bash
go get github.com/scagogogo/go-npm-sdk
```

## Quick Start

Here's a simple example to get you started:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    // Create npm client
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Check if npm is available
    if !client.IsAvailable(ctx) {
        fmt.Println("npm not found, installing...")
        if err := client.Install(ctx); err != nil {
            log.Fatal(err)
        }
        fmt.Println("npm installed successfully!")
    }
    
    // Get npm version
    version, err := client.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("npm version: %s\n", version)
    
    // Install a package
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        SaveDev: false,
        SaveExact: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("lodash installed successfully!")
}
```

## Core Concepts

### Client Interface

The `Client` interface is the main entry point for all npm operations. It provides methods for:

- **Package Management**: Install, uninstall, update packages
- **Project Management**: Initialize projects, manage package.json
- **Script Execution**: Run npm scripts
- **Information Retrieval**: Get package info, search packages
- **Publishing**: Publish packages to registry

### Context Usage

All operations accept a `context.Context` parameter for:

- **Timeout Control**: Set operation timeouts
- **Cancellation**: Cancel long-running operations
- **Request Tracing**: Add request metadata

```go
// Set a timeout for the operation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.InstallPackage(ctx, "express", npm.InstallOptions{})
```

### Error Handling

The SDK provides structured error types for better error handling:

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

## Basic Operations

### Check npm Availability

```go
client, _ := npm.NewClient()
ctx := context.Background()

if client.IsAvailable(ctx) {
    fmt.Println("npm is available")
} else {
    fmt.Println("npm is not available")
}
```

### Install npm Automatically

```go
if !client.IsAvailable(ctx) {
    err := client.Install(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("npm installed successfully")
}
```

### Get npm Version

```go
version, err := client.Version(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("npm version: %s\n", version)
```

### Initialize a Project

```go
options := npm.InitOptions{
    Name:        "my-project",
    Version:     "1.0.0",
    Description: "My awesome project",
    Author:      "Your Name",
    License:     "MIT",
    WorkingDir:  "/path/to/project",
}

err := client.Init(ctx, options)
if err != nil {
    log.Fatal(err)
}
```

### Install Packages

```go
// Install production dependency
err := client.InstallPackage(ctx, "express", npm.InstallOptions{
    SaveDev:   false,
    SaveExact: true,
})

// Install development dependency
err = client.InstallPackage(ctx, "jest", npm.InstallOptions{
    SaveDev: true,
})

// Install global package
err = client.InstallPackage(ctx, "typescript", npm.InstallOptions{
    Global: true,
})
```

### Run Scripts

```go
// Run build script
err := client.RunScript(ctx, "build")

// Run test script with arguments
err = client.RunScript(ctx, "test", "--verbose", "--coverage")
```

### List Packages

```go
packages, err := client.ListPackages(ctx, npm.ListOptions{
    Global: false,
    Depth:  1,
})
if err != nil {
    log.Fatal(err)
}

for _, pkg := range packages {
    fmt.Printf("%s@%s\n", pkg.Name, pkg.Version)
}
```

### Search Packages

```go
results, err := client.Search(ctx, "react hooks")
if err != nil {
    log.Fatal(err)
}

for _, result := range results {
    fmt.Printf("%s@%s - %s\n", 
        result.Package.Name, 
        result.Package.Version, 
        result.Package.Description)
}
```

## Advanced Features

### Portable npm Management

Use portable npm versions without system-wide installation:

```go
import "github.com/scagogogo/go-npm-sdk/pkg/npm"

manager, err := npm.NewPortableManager("/opt/npm-portable")
if err != nil {
    log.Fatal(err)
}

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

// Use the client normally
version, _ := client.Version(ctx)
fmt.Printf("Portable npm version: %s\n", version)
```

### Package.json Management

Direct package.json file manipulation:

```go
import "github.com/scagogogo/go-npm-sdk/pkg/npm"

pkg := npm.NewPackageJSON("./package.json")

// Load existing package.json
err := pkg.Load()
if err != nil {
    log.Fatal(err)
}

// Modify package information
pkg.SetName("my-package")
pkg.SetVersion("2.0.0")
pkg.AddDependency("lodash", "^4.17.21")
pkg.AddScript("build", "webpack")

// Save changes
err = pkg.Save()
if err != nil {
    log.Fatal(err)
}
```

### Platform Detection

Detect the current platform for platform-specific operations:

```go
import "github.com/scagogogo/go-npm-sdk/pkg/platform"

detector := platform.NewDetector()
info, err := detector.Detect()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Platform: %s\n", info.Platform)
fmt.Printf("Architecture: %s\n", info.Architecture)

if info.IsLinux() {
    fmt.Printf("Linux distribution: %s\n", info.Distribution)
}
```

## Configuration

### Working Directory

Set the working directory for npm operations:

```go
options := npm.InstallOptions{
    WorkingDir: "/path/to/project",
}

err := client.InstallPackage(ctx, "express", options)
```

### Registry Configuration

Use custom npm registry:

```go
options := npm.InstallOptions{
    Registry: "https://registry.npmjs.org/",
}

err := client.InstallPackage(ctx, "private-package", options)
```

### Timeout Configuration

Set operation timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := client.InstallPackage(ctx, "large-package", npm.InstallOptions{})
```

## Best Practices

1. **Always use contexts**: Pass appropriate contexts for timeout and cancellation control
2. **Check availability first**: Use `IsAvailable()` before performing operations
3. **Handle errors properly**: Use error type checking for specific error handling
4. **Set working directories**: Specify working directories for project-specific operations
5. **Use structured options**: Configure operations using option structs
6. **Validate inputs**: Always validate package names and versions
7. **Clean up resources**: Ensure proper cleanup of temporary files and processes

## Next Steps

- [Installation Guide](./installation.md) - Detailed installation instructions
- [Configuration](./configuration.md) - Advanced configuration options
- [Platform Support](./platform-support.md) - Platform-specific information
- [API Reference](/en/api/overview.md) - Complete API documentation
- [Examples](/en/examples/basic-usage.md) - More examples and use cases
