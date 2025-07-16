---
layout: home

hero:
  name: "Go NPM SDK"
  text: "A comprehensive Go SDK for npm operations"
  tagline: "Cross-platform npm management with automatic installation, portable versions, and complete API coverage"
  image:
    src: /logo.svg
    alt: Go NPM SDK
  actions:
    - theme: brand
      text: Get Started
      link: /en/guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/go-npm-sdk

features:
  - icon: 🚀
    title: Automatic npm Installation
    details: Automatically detect and install npm based on your operating system using package managers or official installers.
  
  - icon: 📦
    title: Portable Version Support
    details: Download and manage portable Node.js/npm versions without system-wide installation.
  
  - icon: 🔧
    title: Complete API Coverage
    details: Full wrapper for all common npm commands including install, uninstall, update, publish, and more.
  
  - icon: 🌍
    title: Cross-Platform Support
    details: Works seamlessly on Windows, macOS, and Linux with platform-specific optimizations.
  
  - icon: 📝
    title: Project Management
    details: Read, write, and manage package.json files with comprehensive dependency management.
  
  - icon: ⚡
    title: High Performance
    details: Asynchronous execution with timeout control, streaming output, and batch operations.
  
  - icon: 🛡️
    title: Type Safety
    details: Comprehensive error handling with structured error types and validation.
  
  - icon: 🧪
    title: Well Tested
    details: Extensive test coverage (69.7%) with comprehensive unit and integration tests.
  
  - icon: 📚
    title: Rich Documentation
    details: Complete API documentation with examples and best practices.
---

## Quick Start

Install the SDK:

```bash
go get github.com/scagogogo/go-npm-sdk
```

Basic usage:

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
}
```

## Why Go NPM SDK?

- **Zero Configuration**: Works out of the box with automatic npm detection and installation
- **Production Ready**: Used in production environments with comprehensive error handling
- **Developer Friendly**: Intuitive API design with extensive documentation and examples
- **Actively Maintained**: Regular updates and community support

## Community

- [GitHub Issues](https://github.com/scagogogo/go-npm-sdk/issues) - Report bugs and request features
- [GitHub Discussions](https://github.com/scagogogo/go-npm-sdk/discussions) - Ask questions and share ideas
- [Contributing Guide](https://github.com/scagogogo/go-npm-sdk/blob/main/CONTRIBUTING.md) - Learn how to contribute

## License

Released under the [MIT License](https://github.com/scagogogo/go-npm-sdk/blob/main/LICENSE).
