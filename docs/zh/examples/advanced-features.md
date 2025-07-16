# 高级功能示例

本页面演示Go NPM SDK的高级功能和用例。

## 依赖管理

### 依赖解析和冲突检测

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
    
    // 定义要解析的依赖
    dependencies := map[string]string{
        "react":       "^18.0.0",
        "react-dom":   "^18.0.0",
        "lodash":      "^4.17.21",
        "axios":       "^1.0.0",
        "typescript":  "^4.5.0",
    }
    
    fmt.Println("正在解析依赖...")
    
    // 解析依赖
    tree, err := manager.ResolveDependencies(ctx, dependencies)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("依赖解析完成！")
    
    // 检查冲突
    conflicts := manager.CheckConflicts(tree)
    if len(conflicts) > 0 {
        fmt.Println("\n⚠️  发现依赖冲突:")
        for _, conflict := range conflicts {
            fmt.Printf("  - %s: %s vs %s\n", 
                conflict.Package, 
                conflict.Version1, 
                conflict.Version2)
        }
    } else {
        fmt.Println("✅ 未发现依赖冲突！")
    }
    
    // 检查循环依赖
    circular := manager.DetectCircularDependencies(tree)
    if len(circular) > 0 {
        fmt.Println("\n🔄 发现循环依赖:")
        for _, cycle := range circular {
            fmt.Printf("  - %s\n", strings.Join(cycle, " → "))
        }
    } else {
        fmt.Println("✅ 未发现循环依赖！")
    }
}
```

## 批量操作

### 并发包操作

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

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 要并发安装的包
    packages := []string{
        "express", "lodash", "axios", "moment", 
        "uuid", "cors", "helmet", "dotenv",
    }
    
    // 创建批量执行器
    batchExecutor := utils.NewBatchExecutor(3) // 最多3个并发操作
    
    // 准备批量命令
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
    
    // 执行批量安装
    fmt.Printf("正在并发安装%d个包...\n", len(packages))
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
    
    // 报告结果
    fmt.Printf("\n批量安装在%v内完成\n", duration)
    fmt.Printf("成功: %v, 失败: %d/%d\n", 
        result.Success, 
        result.FailedCount, 
        len(result.Results))
    
    // 显示详细结果
    for i, res := range result.Results {
        status := "✅"
        if !res.Success {
            status = "❌"
        }
        fmt.Printf("%s %s (耗时%v)\n", 
            status, 
            packages[i], 
            res.Duration)
    }
}
```

### 并行项目设置

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
    
    // 要设置的多个项目
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
    
    // 并发设置项目
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
            
            fmt.Printf("正在设置%s...\n", proj.name)
            
            // 初始化项目
            err := client.Init(ctx, npm.InitOptions{
                Name:       proj.name,
                Version:    "1.0.0",
                WorkingDir: projectDir,
            })
            if err != nil {
                log.Printf("初始化%s失败: %v", proj.name, err)
                return
            }
            
            // 安装包
            for _, pkg := range proj.packages {
                err := client.InstallPackage(ctx, pkg, npm.InstallOptions{
                    WorkingDir: projectDir,
                })
                if err != nil {
                    log.Printf("在%s中安装%s失败: %v", proj.name, pkg, err)
                }
            }
            
            // 向package.json添加脚本
            pkgJSON := npm.NewPackageJSON(projectDir + "/package.json")
            err = pkgJSON.Load()
            if err != nil {
                log.Printf("为%s加载package.json失败: %v", proj.name, err)
                return
            }
            
            for scriptName, scriptCmd := range proj.scripts {
                pkgJSON.AddScript(scriptName, scriptCmd)
            }
            
            err = pkgJSON.Save()
            if err != nil {
                log.Printf("为%s保存package.json失败: %v", proj.name, err)
                return
            }
            
            fmt.Printf("✅ %s设置完成！\n", proj.name)
        }(project)
    }
    
    wg.Wait()
    fmt.Println("所有项目设置完成！")
}
```

## 自定义注册表和认证

### 使用私有注册表

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
    // 设置带认证的自定义执行器
    executor := utils.NewExecutor()
    
    // 设置认证环境变量
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
    
    // 从私有注册表安装
    privatePackages := []string{
        "@company/shared-components",
        "@company/api-client",
        "@company/utils",
    }
    
    fmt.Println("正在从私有注册表安装包...")
    
    for _, pkg := range privatePackages {
        fmt.Printf("正在安装%s...\n", pkg)
        
        err := client.InstallPackage(ctx, pkg, npm.InstallOptions{
            Registry:   "https://npm.company.com/",
            WorkingDir: "/tmp/private-project",
        })
        
        if err != nil {
            log.Printf("安装%s失败: %v", pkg, err)
        } else {
            fmt.Printf("✅ %s安装成功\n", pkg)
        }
    }
    
    // 发布到私有注册表
    fmt.Println("\n正在发布到私有注册表...")
    
    err = client.Publish(ctx, npm.PublishOptions{
        Registry:   "https://npm.company.com/",
        Access:     "restricted",
        Tag:        "latest",
        WorkingDir: "/tmp/private-project",
    })
    
    if err != nil {
        log.Printf("发布失败: %v", err)
    } else {
        fmt.Println("✅ 包发布成功")
    }
}
```

## 监控和日志

### 高级日志和监控

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
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
    fmt.Println("\n=== 操作摘要 ===")
    
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
        
        fmt.Printf("%s %s %s (耗时%v)\n", 
            status, op.Type, op.Package, duration)
        
        if op.Error != nil {
            fmt.Printf("   错误: %v\n", op.Error)
        }
    }
    
    fmt.Printf("\n总操作数: %d\n", len(ol.operations))
    fmt.Printf("成功: %d\n", successful)
    fmt.Printf("失败: %d\n", failed)
    fmt.Printf("总时间: %v\n", totalDuration)
    fmt.Printf("每个操作平均时间: %v\n", 
        totalDuration/time.Duration(len(ol.operations)))
}

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    logger := &OperationLogger{}
    
    // 监控包安装
    packages := []string{"express", "lodash", "axios", "moment", "uuid"}
    
    for _, pkg := range packages {
        start := time.Now()
        
        fmt.Printf("正在安装%s...\n", pkg)
        err := client.InstallPackage(ctx, pkg, npm.InstallOptions{
            WorkingDir: "/tmp/monitored-project",
        })
        
        end := time.Now()
        success := err == nil
        
        logger.LogOperation("install", pkg, start, end, success, err)
        
        if err != nil {
            fmt.Printf("❌ 安装%s失败: %v\n", pkg, err)
        } else {
            fmt.Printf("✅ %s在%v内安装完成\n", pkg, end.Sub(start))
        }
    }
    
    // 打印详细摘要
    logger.PrintSummary()
}
```

## 错误恢复和重试逻辑

### 带重试的健壮安装

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
        fmt.Printf("尝试%d/%d: 正在安装%s...\n", attempt, maxRetries, pkg)
        
        err := client.InstallPackage(ctx, pkg, options)
        if err == nil {
            fmt.Printf("✅ %s在第%d次尝试时安装成功\n", pkg, attempt)
            return nil
        }
        
        lastErr = err
        
        // 检查是否为可重试错误
        if npm.IsNetworkError(err) || npm.IsNpmNotFound(err) {
            if attempt < maxRetries {
                backoff := time.Duration(attempt) * time.Second
                fmt.Printf("❌ 第%d次尝试失败: %v. %v后重试...\n", 
                    attempt, err, backoff)
                time.Sleep(backoff)
                continue
            }
        } else {
            // 不可重试错误
            fmt.Printf("❌ 不可重试错误: %v\n", err)
            return err
        }
    }
    
    return fmt.Errorf("经过%d次尝试后失败: %v", maxRetries, lastErr)
}

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 使用重试逻辑安装的包
    packages := []string{
        "express",
        "lodash", 
        "axios",
        "nonexistent-package-12345", // 这个会失败
    }
    
    options := npm.InstallOptions{
        WorkingDir: "/tmp/retry-project",
        SaveDev:    false,
    }
    
    for _, pkg := range packages {
        err := installWithRetry(client, ctx, pkg, options, 3)
        if err != nil {
            log.Printf("%s的最终失败: %v", pkg, err)
        }
        fmt.Println()
    }
}
```

## 性能优化

### 缓存和优化

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
    // 首先检查缓存
    pc.mutex.RLock()
    if info, exists := pc.cache[pkg]; exists {
        pc.mutex.RUnlock()
        fmt.Printf("📦 %s缓存命中\n", pkg)
        return info, nil
    }
    pc.mutex.RUnlock()
    
    // 从注册表获取
    fmt.Printf("🌐 正在从注册表获取%s...\n", pkg)
    info, err := client.GetPackageInfo(ctx, pkg)
    if err != nil {
        return nil, err
    }
    
    // 缓存结果
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
        "express", "lodash", // 重复以测试缓存
    }
    
    start := time.Now()
    
    for _, pkg := range packages {
        info, err := cache.GetPackageInfo(client, ctx, pkg)
        if err != nil {
            log.Printf("获取%s信息失败: %v", pkg, err)
            continue
        }
        
        fmt.Printf("📋 %s@%s - %s\n", 
            info.Name, 
            info.Version, 
            info.Description)
    }
    
    duration := time.Since(start)
    fmt.Printf("\n总时间: %v\n", duration)
    fmt.Printf("缓存大小: %d个包\n", len(cache.cache))
}
```

## 集成测试

### 综合测试套件

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
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
    fmt.Printf("运行测试: %s...\n", name)
    
    err := testFunc()
    result := TestResult{
        Name:    name,
        Success: err == nil,
        Error:   err,
    }
    
    ts.results = append(ts.results, result)
    
    if err != nil {
        fmt.Printf("❌ %s失败: %v\n", name, err)
    } else {
        fmt.Printf("✅ %s通过\n", name)
    }
}

func (ts *TestSuite) TestNpmAvailability() error {
    if !ts.client.IsAvailable(ts.ctx) {
        return fmt.Errorf("npm不可用")
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
        return fmt.Errorf("未找到包")
    }
    
    return nil
}

func (ts *TestSuite) TestPackageUninstall() error {
    return ts.client.UninstallPackage(ts.ctx, "lodash", npm.UninstallOptions{
        WorkingDir: ts.testDir,
    })
}

func (ts *TestSuite) PrintResults() {
    fmt.Println("\n=== 测试结果 ===")
    
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
            fmt.Printf("   错误: %v\n", result.Error)
        }
    }
    
    fmt.Printf("\n总测试数: %d\n", len(ts.results))
    fmt.Printf("通过: %d\n", passed)
    fmt.Printf("失败: %d\n", failed)
    
    if failed > 0 {
        fmt.Printf("❌ 测试套件失败\n")
    } else {
        fmt.Printf("✅ 所有测试通过！\n")
    }
}

func main() {
    suite, err := NewTestSuite()
    if err != nil {
        log.Fatal(err)
    }
    
    // 运行测试套件
    suite.RunTest("NPM可用性", suite.TestNpmAvailability)
    suite.RunTest("项目初始化", suite.TestProjectInit)
    suite.RunTest("包安装", suite.TestPackageInstall)
    suite.RunTest("包列表", suite.TestPackageList)
    suite.RunTest("包卸载", suite.TestPackageUninstall)
    
    // 打印最终结果
    suite.PrintResults()
    
    // 清理
    os.RemoveAll(suite.testDir)
}
```

## 最佳实践

1. **错误处理**: 实现带重试逻辑的全面错误处理
2. **性能**: 在适当的地方使用缓存和并发操作
3. **监控**: 记录操作以便调试和性能分析
4. **测试**: 创建全面的测试套件进行验证
5. **安全**: 安全地处理认证和私有注册表
6. **资源管理**: 清理临时文件和目录
7. **可扩展性**: 对多个包使用批量操作

## 下一步

- [基本用法示例](./basic-usage.md)
- [包管理示例](./package-management.md)
