# Package Management Examples

This page demonstrates advanced package management features of the Go NPM SDK.

## Installing Packages with Options

### Development Dependencies

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
    
    // Install development dependencies
    devPackages := []string{"jest", "typescript", "@types/node", "eslint"}
    
    for _, pkg := range devPackages {
        fmt.Printf("Installing %s as dev dependency...\n", pkg)
        err = client.InstallPackage(ctx, pkg, npm.InstallOptions{
            SaveDev:   true,
            SaveExact: true,
        })
        if err != nil {
            log.Printf("Failed to install %s: %v", pkg, err)
        } else {
            fmt.Printf("%s installed successfully!\n", pkg)
        }
    }
}
```

### Production Dependencies

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
    
    // Install production dependencies
    prodPackages := map[string]string{
        "express":    "^4.18.0",
        "lodash":     "^4.17.21",
        "axios":      "^1.0.0",
        "dotenv":     "^16.0.0",
    }
    
    for pkg, version := range prodPackages {
        fmt.Printf("Installing %s@%s...\n", pkg, version)
        err = client.InstallPackage(ctx, pkg+"@"+version, npm.InstallOptions{
            SaveDev:   false,
            SaveExact: false,
        })
        if err != nil {
            log.Printf("Failed to install %s: %v", pkg, err)
        } else {
            fmt.Printf("%s@%s installed successfully!\n", pkg, version)
        }
    }
}
```

## Package Information and Search

### Getting Package Information

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
    
    // Get package information
    packageName := "express"
    info, err := client.GetPackageInfo(ctx, packageName)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Package: %s@%s\n", info.Name, info.Version)
    fmt.Printf("Description: %s\n", info.Description)
    fmt.Printf("Homepage: %s\n", info.Homepage)
    fmt.Printf("License: %s\n", info.License)
    
    if info.Author != nil {
        fmt.Printf("Author: %s <%s>\n", info.Author.Name, info.Author.Email)
    }
    
    if info.Repository != nil {
        fmt.Printf("Repository: %s\n", info.Repository.URL)
    }
    
    fmt.Printf("Keywords: %v\n", info.Keywords)
    
    // Show latest versions
    fmt.Println("\nLatest versions:")
    for tag, version := range info.DistTags {
        fmt.Printf("  %s: %s\n", tag, version)
    }
}
```

### Searching Packages

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
    
    // Search for packages
    query := "react testing"
    results, err := client.Search(ctx, query)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Search results for '%s':\n\n", query)
    
    for i, result := range results[:10] { // Show top 10 results
        fmt.Printf("%d. %s@%s\n", i+1, result.Package.Name, result.Package.Version)
        fmt.Printf("   Description: %s\n", result.Package.Description)
        fmt.Printf("   Score: %.2f (Quality: %.2f, Popularity: %.2f, Maintenance: %.2f)\n",
            result.Score.Final,
            result.Score.Detail.Quality,
            result.Score.Detail.Popularity,
            result.Score.Detail.Maintenance)
        
        if result.Package.Author != nil {
            fmt.Printf("   Author: %s\n", result.Package.Author.Name)
        }
        
        fmt.Println()
    }
}
```

## Package Updates and Maintenance

### Updating Packages

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
    
    // List current packages
    packages, err := client.ListPackages(ctx, npm.ListOptions{
        Global: false,
        Depth:  0,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Current packages:")
    for _, pkg := range packages {
        fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
    }
    
    // Update specific packages
    packagesToUpdate := []string{"lodash", "axios", "express"}
    
    fmt.Println("\nUpdating packages...")
    for _, pkg := range packagesToUpdate {
        fmt.Printf("Updating %s...\n", pkg)
        err = client.UpdatePackage(ctx, pkg)
        if err != nil {
            log.Printf("Failed to update %s: %v", pkg, err)
        } else {
            fmt.Printf("%s updated successfully!\n", pkg)
        }
    }
}
```

### Uninstalling Packages

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
    
    // Uninstall development dependencies
    devPackagesToRemove := []string{"@types/jest", "ts-node"}
    
    for _, pkg := range devPackagesToRemove {
        fmt.Printf("Uninstalling %s...\n", pkg)
        err = client.UninstallPackage(ctx, pkg, npm.UninstallOptions{
            SaveDev: true,
        })
        if err != nil {
            log.Printf("Failed to uninstall %s: %v", pkg, err)
        } else {
            fmt.Printf("%s uninstalled successfully!\n", pkg)
        }
    }
    
    // Uninstall production dependencies
    prodPackagesToRemove := []string{"unused-package"}
    
    for _, pkg := range prodPackagesToRemove {
        fmt.Printf("Uninstalling %s...\n", pkg)
        err = client.UninstallPackage(ctx, pkg, npm.UninstallOptions{
            SaveDev: false,
        })
        if err != nil {
            log.Printf("Failed to uninstall %s: %v", pkg, err)
        } else {
            fmt.Printf("%s uninstalled successfully!\n", pkg)
        }
    }
}
```

## Working with package.json

### Reading and Modifying package.json

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    // Load existing package.json
    pkg := npm.NewPackageJSON("./package.json")
    
    err := pkg.Load()
    if err != nil {
        log.Fatal(err)
    }
    
    // Display current information
    fmt.Printf("Current package: %s@%s\n", pkg.GetName(), pkg.GetVersion())
    fmt.Printf("Description: %s\n", pkg.GetDescription())
    fmt.Printf("Author: %s\n", pkg.GetAuthor())
    
    // Modify package information
    pkg.SetDescription("Updated description for my awesome package")
    pkg.SetAuthor("New Author <new.author@example.com>")
    
    // Add new dependencies
    pkg.AddDependency("moment", "^2.29.0")
    pkg.AddDevDependency("nodemon", "^2.0.0")
    
    // Add scripts
    pkg.AddScript("dev", "nodemon src/index.js")
    pkg.AddScript("build", "webpack --mode production")
    pkg.AddScript("test:watch", "jest --watch")
    
    // Add keywords
    pkg.AddKeyword("nodejs")
    pkg.AddKeyword("javascript")
    pkg.AddKeyword("api")
    
    // Set repository information
    repo := &npm.Repository{
        Type: "git",
        URL:  "https://github.com/username/repo.git",
    }
    pkg.SetRepository(repo)
    
    // Set bugs information
    bugs := &npm.Bugs{
        URL:   "https://github.com/username/repo/issues",
        Email: "bugs@example.com",
    }
    pkg.SetBugs(bugs)
    
    // Set homepage
    pkg.SetHomepage("https://github.com/username/repo#readme")
    
    // Validate before saving
    if err := pkg.Validate(); err != nil {
        log.Fatal(err)
    }
    
    // Save changes
    err = pkg.Save()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("package.json updated successfully!")
}
```

## Global Package Management

### Managing Global Packages

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
    
    // Install global packages
    globalPackages := []string{
        "typescript",
        "nodemon",
        "@angular/cli",
        "create-react-app",
        "eslint",
    }
    
    fmt.Println("Installing global packages...")
    for _, pkg := range globalPackages {
        fmt.Printf("Installing %s globally...\n", pkg)
        err = client.InstallPackage(ctx, pkg, npm.InstallOptions{
            Global: true,
        })
        if err != nil {
            log.Printf("Failed to install %s globally: %v", pkg, err)
        } else {
            fmt.Printf("%s installed globally!\n", pkg)
        }
    }
    
    // List global packages
    fmt.Println("\nListing global packages...")
    globalPkgs, err := client.ListPackages(ctx, npm.ListOptions{
        Global: true,
        Depth:  0,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Global packages:")
    for _, pkg := range globalPkgs {
        fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
    }
}
```

## Registry Configuration

### Using Custom Registry

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
    
    // Install from custom registry
    customRegistry := "https://registry.npmjs.org/"
    
    fmt.Println("Installing package from custom registry...")
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        Registry: customRegistry,
        SaveDev:  false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Package installed from custom registry!")
    
    // Install private package
    privateRegistry := "https://npm.company.com/"
    
    fmt.Println("Installing private package...")
    err = client.InstallPackage(ctx, "@company/private-package", npm.InstallOptions{
        Registry: privateRegistry,
        SaveDev:  false,
    })
    if err != nil {
        log.Printf("Failed to install private package: %v", err)
    } else {
        fmt.Println("Private package installed!")
    }
}
```

## Best Practices

1. **Use exact versions for critical dependencies**: Use `SaveExact: true` for important packages
2. **Separate dev and prod dependencies**: Use appropriate `SaveDev` settings
3. **Validate package.json**: Always validate before saving changes
4. **Handle errors gracefully**: Check for specific error types
5. **Use appropriate registries**: Configure custom registries for private packages
6. **Keep dependencies updated**: Regularly update packages for security
7. **Clean up unused packages**: Remove packages that are no longer needed

## Next Steps

- [Portable Installation Examples](./portable-installation.md)
- [Advanced Features Examples](./advanced-features.md)
