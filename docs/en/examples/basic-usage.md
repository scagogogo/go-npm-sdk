# Basic Usage Examples

This page provides basic usage examples for the Go NPM SDK.

## Getting Started

### Simple npm Client

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
}
```

### Installing Packages

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install production dependency
    fmt.Println("Installing lodash...")
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        SaveDev:   false,
        SaveExact: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("lodash installed successfully!")
    
    // Install development dependency
    fmt.Println("Installing jest...")
    err = client.InstallPackage(ctx, "jest", npm.InstallOptions{
        SaveDev: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("jest installed successfully!")
    
    // Install global package
    fmt.Println("Installing typescript globally...")
    err = client.InstallPackage(ctx, "typescript", npm.InstallOptions{
        Global: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("typescript installed globally!")
}
```

### Project Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Initialize a new project
    options := npm.InitOptions{
        Name:        "my-awesome-project",
        Version:     "1.0.0",
        Description: "An awesome Node.js project",
        Author:      "Your Name <your.email@example.com>",
        License:     "MIT",
        Private:     false,
        WorkingDir:  "./my-project",
    }
    
    fmt.Println("Initializing project...")
    err = client.Init(ctx, options)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Project initialized successfully!")
}
```

### Running Scripts

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Run build script
    fmt.Println("Running build script...")
    err = client.RunScript(ctx, "build")
    if err != nil {
        log.Printf("Build failed: %v", err)
    } else {
        fmt.Println("Build completed successfully!")
    }
    
    // Run test script with arguments
    fmt.Println("Running tests...")
    err = client.RunScript(ctx, "test", "--verbose", "--coverage")
    if err != nil {
        log.Printf("Tests failed: %v", err)
    } else {
        fmt.Println("Tests passed!")
    }
}
```

### Listing Packages

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // List local packages
    fmt.Println("Local packages:")
    packages, err := client.ListPackages(ctx, npm.ListOptions{
        Global: false,
        Depth:  1,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, pkg := range packages {
        fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
    }
    
    // List global packages
    fmt.Println("\nGlobal packages:")
    globalPackages, err := client.ListPackages(ctx, npm.ListOptions{
        Global: true,
        Depth:  0,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, pkg := range globalPackages {
        fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
    }
}
```

### Error Handling

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Try to install a non-existent package
    err = client.InstallPackage(ctx, "this-package-does-not-exist-12345", npm.InstallOptions{})
    if err != nil {
        // Check for specific error types
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
}
```

### Using Context for Timeout

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Install package with timeout
    fmt.Println("Installing package with 30-second timeout...")
    err = client.InstallPackage(ctx, "express", npm.InstallOptions{})
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Installation timed out")
        } else {
            log.Fatal(err)
        }
    } else {
        fmt.Println("Package installed successfully!")
    }
}
```

### Working with Different Directories

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
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Create a temporary directory
    tempDir, err := os.MkdirTemp("", "npm-example")
    if err != nil {
        log.Fatal(err)
    }
    defer os.RemoveAll(tempDir)
    
    fmt.Printf("Working in directory: %s\n", tempDir)
    
    // Initialize project in specific directory
    err = client.Init(ctx, npm.InitOptions{
        Name:       "temp-project",
        Version:    "1.0.0",
        WorkingDir: tempDir,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Install package in specific directory
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        WorkingDir: tempDir,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Project created and package installed in temporary directory!")
}
```

## Best Practices

1. **Always use context**: Pass appropriate contexts for timeout and cancellation control
2. **Check availability first**: Use `IsAvailable()` before performing operations
3. **Handle errors properly**: Use error type checking for specific error handling
4. **Set working directories**: Specify working directories for project-specific operations
5. **Use timeouts**: Set reasonable timeouts for long-running operations
6. **Clean up resources**: Ensure proper cleanup of temporary files and directories

## Next Steps

- [Package Management Examples](./package-management.md)
- [Portable Installation Examples](./portable-installation.md)
- [Advanced Features Examples](./advanced-features.md)
