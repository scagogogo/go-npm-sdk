# Portable Installation Examples

This page demonstrates how to use portable Node.js/npm installations with the Go NPM SDK.

## Basic Portable Installation

### Setting Up Portable Manager

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    // Create portable manager with custom directory
    portableDir := "/opt/npm-portable"
    
    // Ensure directory exists
    err := os.MkdirAll(portableDir, 0755)
    if err != nil {
        log.Fatal(err)
    }
    
    manager, err := npm.NewPortableManager(portableDir)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install Node.js 18.17.0
    fmt.Println("Installing Node.js 18.17.0...")
    config, err := manager.Install(ctx, "18.17.0")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Node.js installed successfully!\n")
    fmt.Printf("Version: %s\n", config.Version)
    fmt.Printf("Install Path: %s\n", config.InstallPath)
    fmt.Printf("Node Path: %s\n", config.NodePath)
    fmt.Printf("NPM Path: %s\n", config.NpmPath)
}
```

### Installing Multiple Versions

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    manager, err := npm.NewPortableManager("/opt/npm-portable")
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install multiple Node.js versions
    versions := []string{"16.20.0", "18.17.0", "20.5.0"}
    
    for _, version := range versions {
        fmt.Printf("Installing Node.js %s...\n", version)
        
        config, err := manager.Install(ctx, version)
        if err != nil {
            log.Printf("Failed to install %s: %v", version, err)
            continue
        }
        
        fmt.Printf("Node.js %s installed at %s\n", version, config.InstallPath)
    }
    
    // List all installed versions
    fmt.Println("\nInstalled versions:")
    configs, err := manager.List()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, config := range configs {
        fmt.Printf("  %s - %s\n", config.Version, config.InstallPath)
    }
}
```

## Using Portable Installations

### Creating Clients for Specific Versions

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    manager, err := npm.NewPortableManager("/opt/npm-portable")
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Create client for Node.js 18.17.0
    client18, err := manager.CreateClient("18.17.0")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create client for Node.js 20.5.0
    client20, err := manager.CreateClient("20.5.0")
    if err != nil {
        log.Fatal(err)
    }
    
    // Use different versions
    version18, err := client18.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Client 18 npm version: %s\n", version18)
    
    version20, err := client20.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Client 20 npm version: %s\n", version20)
    
    // Install packages with different versions
    fmt.Println("Installing lodash with Node.js 18...")
    err = client18.InstallPackage(ctx, "lodash", npm.InstallOptions{
        WorkingDir: "/tmp/project18",
    })
    if err != nil {
        log.Printf("Failed: %v", err)
    }
    
    fmt.Println("Installing lodash with Node.js 20...")
    err = client20.InstallPackage(ctx, "lodash", npm.InstallOptions{
        WorkingDir: "/tmp/project20",
    })
    if err != nil {
        log.Printf("Failed: %v", err)
    }
}
```

### Version Switching

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    manager, err := npm.NewPortableManager("/opt/npm-portable")
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Function to test a version
    testVersion := func(version string) {
        fmt.Printf("\n=== Testing Node.js %s ===\n", version)
        
        client, err := manager.CreateClient(version)
        if err != nil {
            log.Printf("Failed to create client for %s: %v", version, err)
            return
        }
        
        // Check if npm is available
        if !client.IsAvailable(ctx) {
            log.Printf("npm not available for version %s", version)
            return
        }
        
        // Get npm version
        npmVersion, err := client.Version(ctx)
        if err != nil {
            log.Printf("Failed to get npm version: %v", err)
            return
        }
        
        fmt.Printf("npm version: %s\n", npmVersion)
        
        // Create a test project
        tempDir := fmt.Sprintf("/tmp/test-project-%s", version)
        err = client.Init(ctx, npm.InitOptions{
            Name:       fmt.Sprintf("test-project-%s", version),
            Version:    "1.0.0",
            WorkingDir: tempDir,
        })
        if err != nil {
            log.Printf("Failed to init project: %v", err)
            return
        }
        
        fmt.Printf("Test project created at %s\n", tempDir)
    }
    
    // Test different versions
    versions := []string{"16.20.0", "18.17.0", "20.5.0"}
    for _, version := range versions {
        testVersion(version)
    }
}
```

## Advanced Portable Management

### Custom Installation with Progress

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    manager, err := npm.NewPortableManager("/opt/npm-portable")
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install with progress callback
    version := "18.17.0"
    fmt.Printf("Installing Node.js %s with progress tracking...\n", version)
    
    // Note: This is a conceptual example - actual implementation may vary
    config, err := manager.Install(ctx, version)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Installation completed!\n")
    fmt.Printf("Installed at: %s\n", config.InstallPath)
    fmt.Printf("Install date: %s\n", config.InstallDate)
}
```

### Cleaning Up Old Versions

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    manager, err := npm.NewPortableManager("/opt/npm-portable")
    if err != nil {
        log.Fatal(err)
    }
    
    // List all installed versions
    configs, err := manager.List()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Currently installed versions:")
    for _, config := range configs {
        fmt.Printf("  %s (installed: %s)\n", config.Version, config.InstallDate)
    }
    
    // Remove old versions (keep only latest 2)
    if len(configs) > 2 {
        versionsToRemove := configs[:len(configs)-2]
        
        for _, config := range versionsToRemove {
            fmt.Printf("Removing Node.js %s...\n", config.Version)
            err = manager.Uninstall(config.Version)
            if err != nil {
                log.Printf("Failed to remove %s: %v", config.Version, err)
            } else {
                fmt.Printf("Node.js %s removed successfully\n", config.Version)
            }
        }
    }
    
    // List remaining versions
    fmt.Println("\nRemaining versions:")
    configs, err = manager.List()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, config := range configs {
        fmt.Printf("  %s\n", config.Version)
    }
}
```

## Integration with CI/CD

### Docker-like Environment

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    // Setup portable environment for CI/CD
    ciDir := "/ci/npm-portable"
    
    // Clean up any existing installation
    os.RemoveAll(ciDir)
    os.MkdirAll(ciDir, 0755)
    
    manager, err := npm.NewPortableManager(ciDir)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install specific Node.js version for CI
    nodeVersion := "18.17.0"
    fmt.Printf("Setting up CI environment with Node.js %s...\n", nodeVersion)
    
    config, err := manager.Install(ctx, nodeVersion)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create client for CI operations
    client, err := manager.CreateClient(nodeVersion)
    if err != nil {
        log.Fatal(err)
    }
    
    // Verify installation
    npmVersion, err := client.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("CI environment ready!\n")
    fmt.Printf("Node.js: %s\n", nodeVersion)
    fmt.Printf("npm: %s\n", npmVersion)
    fmt.Printf("Installation path: %s\n", config.InstallPath)
    
    // Run CI tasks
    projectDir := "/ci/project"
    os.MkdirAll(projectDir, 0755)
    
    // Initialize project
    err = client.Init(ctx, npm.InitOptions{
        Name:       "ci-project",
        Version:    "1.0.0",
        WorkingDir: projectDir,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Install dependencies
    dependencies := []string{"express", "lodash", "axios"}
    for _, dep := range dependencies {
        fmt.Printf("Installing %s...\n", dep)
        err = client.InstallPackage(ctx, dep, npm.InstallOptions{
            WorkingDir: projectDir,
        })
        if err != nil {
            log.Printf("Failed to install %s: %v", dep, err)
        }
    }
    
    // Install dev dependencies
    devDeps := []string{"jest", "eslint", "typescript"}
    for _, dep := range devDeps {
        fmt.Printf("Installing %s as dev dependency...\n", dep)
        err = client.InstallPackage(ctx, dep, npm.InstallOptions{
            WorkingDir: projectDir,
            SaveDev:    true,
        })
        if err != nil {
            log.Printf("Failed to install %s: %v", dep, err)
        }
    }
    
    fmt.Println("CI environment setup completed!")
}
```

## Platform-Specific Portable Installations

### Cross-Platform Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "runtime"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
    "github.com/scagogogo/go-npm-sdk/pkg/platform"
)

func main() {
    // Detect current platform
    detector := platform.NewDetector()
    info, err := detector.Detect()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Detected platform: %s %s\n", info.Platform, info.Architecture)
    
    // Setup platform-specific portable directory
    var portableDir string
    switch runtime.GOOS {
    case "windows":
        portableDir = "C:\\npm-portable"
    case "darwin":
        portableDir = "/usr/local/npm-portable"
    default: // linux
        portableDir = "/opt/npm-portable"
    }
    
    manager, err := npm.NewPortableManager(portableDir)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install Node.js version appropriate for platform
    version := "18.17.0"
    fmt.Printf("Installing Node.js %s for %s...\n", version, info.Platform)
    
    config, err := manager.Install(ctx, version)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Installation successful!\n")
    fmt.Printf("Platform: %s\n", info.Platform)
    fmt.Printf("Architecture: %s\n", info.Architecture)
    fmt.Printf("Install path: %s\n", config.InstallPath)
    fmt.Printf("Node path: %s\n", config.NodePath)
    fmt.Printf("NPM path: %s\n", config.NpmPath)
}
```

## Best Practices

1. **Use specific versions**: Always specify exact Node.js versions for reproducible environments
2. **Organize by purpose**: Use different portable directories for different projects or environments
3. **Clean up regularly**: Remove unused versions to save disk space
4. **Version testing**: Test your application with multiple Node.js versions
5. **CI/CD integration**: Use portable installations for consistent CI/CD environments
6. **Platform awareness**: Consider platform-specific paths and behaviors
7. **Backup configurations**: Keep track of which versions work with your projects

## Next Steps

- [Advanced Features Examples](./advanced-features.md)
- [Basic Usage Examples](./basic-usage.md)
