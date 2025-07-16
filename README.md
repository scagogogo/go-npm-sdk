# Go NPM SDK

[English](README.md) | [简体中文](README.zh.md)

A comprehensive Go SDK for npm operations with cross-platform support.

## 📚 Documentation

**🌐 [Complete Documentation Website](https://scagogogo.github.io/go-npm-sdk/)**

Visit our comprehensive documentation website for detailed guides, API references, and examples.

## Features

- **Automatic npm Installation**: Detect and install npm automatically based on your operating system
- **Portable Version Support**: Download and manage portable Node.js/npm versions
- **Complete API Coverage**: Full wrapper for all common npm commands
- **Cross-Platform Support**: Works on Windows, macOS, and Linux
- **Project Management**: Read, write, and manage package.json files
- **High Performance**: Asynchronous execution with timeout control
- **Type Safety**: Comprehensive error handling with structured error types

## Installation

```bash
go get github.com/scagogogo/go-npm-sdk
```

## Quick Start

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
        // Auto-install npm
        if err := client.Install(ctx); err != nil {
            log.Fatal(err)
        }
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
    fmt.Println("Package installed successfully!")
}
```

> 💡 **Need more examples?** Check out our [comprehensive documentation](https://scagogogo.github.io/go-npm-sdk/) for detailed guides and advanced usage patterns.

## Core Features

### Automatic npm Installation

The SDK can automatically detect and install npm based on your operating system:

```go
client, _ := npm.NewClient()
ctx := context.Background()

if !client.IsAvailable(ctx) {
    // Automatically install npm using the best method for your OS
    err := client.Install(ctx)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Portable npm Management

Download and manage portable Node.js/npm versions without system-wide installation:

```go
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
```

### Package.json Management

Read, write, and manage package.json files:

```go
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

## API Documentation

For complete API documentation, visit our documentation website:

**📚 [Complete Documentation Website](https://scagogogo.github.io/go-npm-sdk/)**

The documentation includes:
- Complete API reference
- Usage guides and tutorials
- Examples and best practices
- Platform-specific information

## Examples

See the [examples](./examples/) directory for more comprehensive examples:

- [Basic Usage](./examples/basic_usage.go) - Getting started with the SDK
- [Package Management](./examples/package_management.go) - Installing and managing packages
- [Portable Installation](./examples/portable_installation.go) - Using portable npm versions
- [Platform Detection](./examples/platform_detection.go) - Detecting platform information
- [Dependency Management](./examples/dependency_management.go) - Managing dependencies

## Supported Platforms

- **Windows**: Windows 10/11, Windows Server 2019/2022
- **macOS**: macOS 10.15+ (Intel and Apple Silicon)
- **Linux**: Ubuntu, Debian, CentOS, RHEL, Fedora, SUSE, Arch, Alpine

## Installation Methods

The SDK supports multiple npm installation methods:

1. **Package Manager**: Use system package managers (apt, yum, brew, etc.)
2. **Official Installer**: Download and run official Node.js installer
3. **Portable**: Download portable Node.js/npm version
4. **Manual**: Manual installation guidance

## Requirements

- Go 1.19 or later
- Internet connection for downloading npm/Node.js (if not already installed)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- [GitHub Issues](https://github.com/scagogogo/go-npm-sdk/issues) - Report bugs and request features
- [GitHub Discussions](https://github.com/scagogogo/go-npm-sdk/discussions) - Ask questions and share ideas
- **[📖 Documentation Website](https://scagogogo.github.io/go-npm-sdk/)** - Complete documentation and guides
