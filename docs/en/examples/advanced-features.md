# Advanced Features Examples

This page demonstrates advanced features and use cases of the Go NPM SDK.

## Dependency Management

### Dependency Resolution and Conflict Detection

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    manager := npm.NewDependencyManager()
    ctx := context.Background()
    
    // Define dependencies to resolve
    dependencies := map[string]string{
        "react":       "^18.0.0",
        "react-dom":   "^18.0.0",
        "lodash":      "^4.17.21",
        "axios":       "^1.0.0",
        "typescript":  "^4.5.0",
    }
    
    fmt.Println("Resolving dependencies...")
    
    // Resolve dependencies
    tree, err := manager.ResolveDependencies(ctx, dependencies)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Dependency resolution completed!")
    
    // Check for conflicts
    conflicts := manager.CheckConflicts(tree)
    if len(conflicts) > 0 {
        fmt.Println("\n⚠️  Dependency conflicts found:")
        for _, conflict := range conflicts {
            fmt.Printf("  - %s: %s vs %s\n", 
                conflict.Package, 
                conflict.Version1, 
                conflict.Version2)
        }
    } else {
        fmt.Println("✅ No dependency conflicts found!")
    }
    
    // Check for circular dependencies
    circular := manager.DetectCircularDependencies(tree)
    if len(circular) > 0 {
        fmt.Println("\n🔄 Circular dependencies found:")
        for _, cycle := range circular {
            fmt.Printf("  - %s\n", strings.Join(cycle, " → "))
        }
    } else {
        fmt.Println("✅ No circular dependencies found!")
    }
}
```

## Batch Operations

### Concurrent Package Operations

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
    "github.com/scagogogo/go-npm-sdk/pkg/utils"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Packages to install concurrently
    packages := []string{
        "express", "lodash", "axios", "moment", 
        "uuid", "cors", "helmet", "dotenv",
    }
    
    // Create batch executor
    batchExecutor := utils.NewBatchExecutor(3) // Max 3 concurrent operations
    
    // Prepare batch commands
    var commands []utils.ExecuteOptions
    for _, pkg := range packages {
        commands = append(commands, utils.ExecuteOptions{
            Command:       "npm",
            Args:          []string{"install", pkg},
            WorkingDir:    "/tmp/batch-project",
            CaptureOutput: true,
            Timeout:       2 * time.Minute,
        })
    }
    
    // Execute batch installation
    fmt.Printf("Installing %d packages concurrently...\n", len(packages))
    start := time.Now()
    
    result, err := batchExecutor.ExecuteBatch(ctx, utils.BatchOptions{
        Commands:       commands,
        StopOnError:    false,
        MaxConcurrency: 3,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    duration := time.Since(start)
    
    // Report results
    fmt.Printf("\nBatch installation completed in %v\n", duration)
    fmt.Printf("Success: %v, Failed: %d/%d\n", 
        result.Success, 
        result.FailedCount, 
        len(result.Results))
    
    // Show detailed results
    for i, res := range result.Results {
        status := "✅"
        if !res.Success {
            status = "❌"
        }
        fmt.Printf("%s %s (took %v)\n", 
            status, 
            packages[i], 
            res.Duration)
    }
}
```

### Parallel Project Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "sync"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Multiple projects to set up
    projects := []struct {
        name     string
        packages []string
        scripts  map[string]string
    }{
        {
            name:     "frontend-app",
            packages: []string{"react", "react-dom", "typescript"},
            scripts: map[string]string{
                "start": "react-scripts start",
                "build": "react-scripts build",
                "test":  "react-scripts test",
            },
        },
        {
            name:     "backend-api",
            packages: []string{"express", "cors", "helmet", "dotenv"},
            scripts: map[string]string{
                "start": "node src/index.js",
                "dev":   "nodemon src/index.js",
                "test":  "jest",
            },
        },
        {
            name:     "shared-utils",
            packages: []string{"lodash", "moment", "uuid"},
            scripts: map[string]string{
                "build": "tsc",
                "test":  "jest",
            },
        },
    }
    
    var wg sync.WaitGroup
    
    // Setup projects concurrently
    for _, project := range projects {
        wg.Add(1)
        go func(proj struct {
            name     string
            packages []string
            scripts  map[string]string
        }) {
            defer wg.Done()
            
            projectDir := fmt.Sprintf("/tmp/projects/%s", proj.name)
            os.MkdirAll(projectDir, 0755)
            
            fmt.Printf("Setting up %s...\n", proj.name)
            
            // Initialize project
            err := client.Init(ctx, npm.InitOptions{
                Name:       proj.name,
                Version:    "1.0.0",
                WorkingDir: projectDir,
            })
            if err != nil {
                log.Printf("Failed to init %s: %v", proj.name, err)
                return
            }
            
            // Install packages
            for _, pkg := range proj.packages {
                err := client.InstallPackage(ctx, pkg, npm.InstallOptions{
                    WorkingDir: projectDir,
                })
                if err != nil {
                    log.Printf("Failed to install %s in %s: %v", pkg, proj.name, err)
                }
            }
            
            // Add scripts to package.json
            pkgJSON := npm.NewPackageJSON(projectDir + "/package.json")
            err = pkgJSON.Load()
            if err != nil {
                log.Printf("Failed to load package.json for %s: %v", proj.name, err)
                return
            }
            
            for scriptName, scriptCmd := range proj.scripts {
                pkgJSON.AddScript(scriptName, scriptCmd)
            }
            
            err = pkgJSON.Save()
            if err != nil {
                log.Printf("Failed to save package.json for %s: %v", proj.name, err)
                return
            }
            
            fmt.Printf("✅ %s setup completed!\n", proj.name)
        }(project)
    }
    
    wg.Wait()
    fmt.Println("All projects setup completed!")
}
```

## Custom Registry and Authentication

### Working with Private Registries

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
    "github.com/scagogogo/go-npm-sdk/pkg/utils"
)

func main() {
    // Setup custom executor with authentication
    executor := utils.NewExecutor()
    
    // Set environment variables for authentication
    executor.SetDefaultEnv(map[string]string{
        "NPM_TOKEN":           os.Getenv("NPM_TOKEN"),
        "NPM_CONFIG_REGISTRY": "https://npm.company.com/",
        "NPM_CONFIG_ALWAYS_AUTH": "true",
    })
    
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Install from private registry
    privatePackages := []string{
        "@company/shared-components",
        "@company/api-client",
        "@company/utils",
    }
    
    fmt.Println("Installing packages from private registry...")
    
    for _, pkg := range privatePackages {
        fmt.Printf("Installing %s...\n", pkg)
        
        err := client.InstallPackage(ctx, pkg, npm.InstallOptions{
            Registry:   "https://npm.company.com/",
            WorkingDir: "/tmp/private-project",
        })
        
        if err != nil {
            log.Printf("Failed to install %s: %v", pkg, err)
        } else {
            fmt.Printf("✅ %s installed successfully\n", pkg)
        }
    }
    
    // Publish to private registry
    fmt.Println("\nPublishing to private registry...")
    
    err = client.Publish(ctx, npm.PublishOptions{
        Registry:   "https://npm.company.com/",
        Access:     "restricted",
        Tag:        "latest",
        WorkingDir: "/tmp/private-project",
    })
    
    if err != nil {
        log.Printf("Failed to publish: %v", err)
    } else {
        fmt.Println("✅ Package published successfully")
    }
}
```

## Monitoring and Logging

### Advanced Logging and Monitoring

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
    "github.com/scagogogo/go-npm-sdk/pkg/utils"
)

type OperationLogger struct {
    operations []Operation
}

type Operation struct {
    Type      string
    Package   string
    StartTime time.Time
    EndTime   time.Time
    Success   bool
    Error     error
}

func (ol *OperationLogger) LogOperation(opType, pkg string, start, end time.Time, success bool, err error) {
    ol.operations = append(ol.operations, Operation{
        Type:      opType,
        Package:   pkg,
        StartTime: start,
        EndTime:   end,
        Success:   success,
        Error:     err,
    })
}

func (ol *OperationLogger) PrintSummary() {
    fmt.Println("\n=== Operation Summary ===")
    
    var successful, failed int
    var totalDuration time.Duration
    
    for _, op := range ol.operations {
        duration := op.EndTime.Sub(op.StartTime)
        totalDuration += duration
        
        status := "✅"
        if !op.Success {
            status = "❌"
            failed++
        } else {
            successful++
        }
        
        fmt.Printf("%s %s %s (took %v)\n", 
            status, op.Type, op.Package, duration)
        
        if op.Error != nil {
            fmt.Printf("   Error: %v\n", op.Error)
        }
    }
    
    fmt.Printf("\nTotal operations: %d\n", len(ol.operations))
    fmt.Printf("Successful: %d\n", successful)
    fmt.Printf("Failed: %d\n", failed)
    fmt.Printf("Total time: %v\n", totalDuration)
    fmt.Printf("Average time per operation: %v\n", 
        totalDuration/time.Duration(len(ol.operations)))
}

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    logger := &OperationLogger{}
    
    // Monitored package installation
    packages := []string{"express", "lodash", "axios", "moment", "uuid"}
    
    for _, pkg := range packages {
        start := time.Now()
        
        fmt.Printf("Installing %s...\n", pkg)
        err := client.InstallPackage(ctx, pkg, npm.InstallOptions{
            WorkingDir: "/tmp/monitored-project",
        })
        
        end := time.Now()
        success := err == nil
        
        logger.LogOperation("install", pkg, start, end, success, err)
        
        if err != nil {
            fmt.Printf("❌ Failed to install %s: %v\n", pkg, err)
        } else {
            fmt.Printf("✅ %s installed in %v\n", pkg, end.Sub(start))
        }
    }
    
    // Print detailed summary
    logger.PrintSummary()
}
```

## Error Recovery and Retry Logic

### Robust Installation with Retry

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func installWithRetry(client npm.Client, ctx context.Context, pkg string, options npm.InstallOptions, maxRetries int) error {
    var lastErr error
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("Attempt %d/%d: Installing %s...\n", attempt, maxRetries, pkg)
        
        err := client.InstallPackage(ctx, pkg, options)
        if err == nil {
            fmt.Printf("✅ %s installed successfully on attempt %d\n", pkg, attempt)
            return nil
        }
        
        lastErr = err
        
        // Check if it's a retryable error
        if npm.IsNetworkError(err) || npm.IsNpmNotFound(err) {
            if attempt < maxRetries {
                backoff := time.Duration(attempt) * time.Second
                fmt.Printf("❌ Attempt %d failed: %v. Retrying in %v...\n", 
                    attempt, err, backoff)
                time.Sleep(backoff)
                continue
            }
        } else {
            // Non-retryable error
            fmt.Printf("❌ Non-retryable error: %v\n", err)
            return err
        }
    }
    
    return fmt.Errorf("failed after %d attempts: %v", maxRetries, lastErr)
}

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Packages to install with retry logic
    packages := []string{
        "express",
        "lodash", 
        "axios",
        "nonexistent-package-12345", // This will fail
    }
    
    options := npm.InstallOptions{
        WorkingDir: "/tmp/retry-project",
        SaveDev:    false,
    }
    
    for _, pkg := range packages {
        err := installWithRetry(client, ctx, pkg, options, 3)
        if err != nil {
            log.Printf("Final failure for %s: %v", pkg, err)
        }
        fmt.Println()
    }
}
```

## Performance Optimization

### Caching and Optimization

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

type PackageCache struct {
    cache map[string]*npm.PackageInfo
    mutex sync.RWMutex
}

func NewPackageCache() *PackageCache {
    return &PackageCache{
        cache: make(map[string]*npm.PackageInfo),
    }
}

func (pc *PackageCache) GetPackageInfo(client npm.Client, ctx context.Context, pkg string) (*npm.PackageInfo, error) {
    // Check cache first
    pc.mutex.RLock()
    if info, exists := pc.cache[pkg]; exists {
        pc.mutex.RUnlock()
        fmt.Printf("📦 Cache hit for %s\n", pkg)
        return info, nil
    }
    pc.mutex.RUnlock()
    
    // Fetch from registry
    fmt.Printf("🌐 Fetching %s from registry...\n", pkg)
    info, err := client.GetPackageInfo(ctx, pkg)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    pc.mutex.Lock()
    pc.cache[pkg] = info
    pc.mutex.Unlock()
    
    return info, nil
}

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    cache := NewPackageCache()
    
    packages := []string{
        "express", "lodash", "axios", "moment", "uuid",
        "express", "lodash", // Repeated to test cache
    }
    
    start := time.Now()
    
    for _, pkg := range packages {
        info, err := cache.GetPackageInfo(client, ctx, pkg)
        if err != nil {
            log.Printf("Failed to get info for %s: %v", pkg, err)
            continue
        }
        
        fmt.Printf("📋 %s@%s - %s\n", 
            info.Name, 
            info.Version, 
            info.Description)
    }
    
    duration := time.Since(start)
    fmt.Printf("\nTotal time: %v\n", duration)
    fmt.Printf("Cache size: %d packages\n", len(cache.cache))
}
```

## Integration Testing

### Comprehensive Testing Suite

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "testing"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

type TestSuite struct {
    client    npm.Client
    testDir   string
    ctx       context.Context
    results   []TestResult
}

type TestResult struct {
    Name    string
    Success bool
    Error   error
}

func NewTestSuite() (*TestSuite, error) {
    client, err := npm.NewClient()
    if err != nil {
        return nil, err
    }
    
    testDir := "/tmp/npm-sdk-test"
    os.RemoveAll(testDir)
    os.MkdirAll(testDir, 0755)
    
    return &TestSuite{
        client:  client,
        testDir: testDir,
        ctx:     context.Background(),
    }, nil
}

func (ts *TestSuite) RunTest(name string, testFunc func() error) {
    fmt.Printf("Running test: %s...\n", name)
    
    err := testFunc()
    result := TestResult{
        Name:    name,
        Success: err == nil,
        Error:   err,
    }
    
    ts.results = append(ts.results, result)
    
    if err != nil {
        fmt.Printf("❌ %s failed: %v\n", name, err)
    } else {
        fmt.Printf("✅ %s passed\n", name)
    }
}

func (ts *TestSuite) TestNpmAvailability() error {
    if !ts.client.IsAvailable(ts.ctx) {
        return fmt.Errorf("npm is not available")
    }
    return nil
}

func (ts *TestSuite) TestProjectInit() error {
    return ts.client.Init(ts.ctx, npm.InitOptions{
        Name:       "test-project",
        Version:    "1.0.0",
        WorkingDir: ts.testDir,
    })
}

func (ts *TestSuite) TestPackageInstall() error {
    return ts.client.InstallPackage(ts.ctx, "lodash", npm.InstallOptions{
        WorkingDir: ts.testDir,
    })
}

func (ts *TestSuite) TestPackageList() error {
    packages, err := ts.client.ListPackages(ts.ctx, npm.ListOptions{
        WorkingDir: ts.testDir,
    })
    if err != nil {
        return err
    }
    
    if len(packages) == 0 {
        return fmt.Errorf("no packages found")
    }
    
    return nil
}

func (ts *TestSuite) TestPackageUninstall() error {
    return ts.client.UninstallPackage(ts.ctx, "lodash", npm.UninstallOptions{
        WorkingDir: ts.testDir,
    })
}

func (ts *TestSuite) PrintResults() {
    fmt.Println("\n=== Test Results ===")
    
    var passed, failed int
    for _, result := range ts.results {
        status := "✅"
        if !result.Success {
            status = "❌"
            failed++
        } else {
            passed++
        }
        
        fmt.Printf("%s %s\n", status, result.Name)
        if result.Error != nil {
            fmt.Printf("   Error: %v\n", result.Error)
        }
    }
    
    fmt.Printf("\nTotal tests: %d\n", len(ts.results))
    fmt.Printf("Passed: %d\n", passed)
    fmt.Printf("Failed: %d\n", failed)
    
    if failed > 0 {
        fmt.Printf("❌ Test suite failed\n")
    } else {
        fmt.Printf("✅ All tests passed!\n")
    }
}

func main() {
    suite, err := NewTestSuite()
    if err != nil {
        log.Fatal(err)
    }
    
    // Run test suite
    suite.RunTest("NPM Availability", suite.TestNpmAvailability)
    suite.RunTest("Project Initialization", suite.TestProjectInit)
    suite.RunTest("Package Installation", suite.TestPackageInstall)
    suite.RunTest("Package Listing", suite.TestPackageList)
    suite.RunTest("Package Uninstallation", suite.TestPackageUninstall)
    
    // Print final results
    suite.PrintResults()
    
    // Cleanup
    os.RemoveAll(suite.testDir)
}
```

## Best Practices

1. **Error Handling**: Implement comprehensive error handling with retry logic
2. **Performance**: Use caching and concurrent operations where appropriate
3. **Monitoring**: Log operations for debugging and performance analysis
4. **Testing**: Create comprehensive test suites for validation
5. **Security**: Handle authentication and private registries securely
6. **Resource Management**: Clean up temporary files and directories
7. **Scalability**: Use batch operations for multiple packages

## Next Steps

- [Basic Usage Examples](./basic-usage.md)
- [Package Management Examples](./package-management.md)
