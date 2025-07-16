# 便携版安装示例

本页面演示如何使用Go NPM SDK进行便携版Node.js/npm安装。

## 基本便携版安装

### 设置便携版管理器

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
    // 使用自定义目录创建便携版管理器
    portableDir := "/opt/npm-portable"
    
    // 确保目录存在
    err := os.MkdirAll(portableDir, 0755)
    if err != nil {
        log.Fatal(err)
    }
    
    manager, err := npm.NewPortableManager(portableDir)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 安装Node.js 18.17.0
    fmt.Println("正在安装Node.js 18.17.0...")
    config, err := manager.Install(ctx, "18.17.0")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Node.js安装成功！\n")
    fmt.Printf("版本: %s\n", config.Version)
    fmt.Printf("安装路径: %s\n", config.InstallPath)
    fmt.Printf("Node路径: %s\n", config.NodePath)
    fmt.Printf("NPM路径: %s\n", config.NpmPath)
}
```

### 安装多个版本

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
    
    // 安装多个Node.js版本
    versions := []string{"16.20.0", "18.17.0", "20.5.0"}
    
    for _, version := range versions {
        fmt.Printf("正在安装Node.js %s...\n", version)
        
        config, err := manager.Install(ctx, version)
        if err != nil {
            log.Printf("安装%s失败: %v", version, err)
            continue
        }
        
        fmt.Printf("Node.js %s安装在%s\n", version, config.InstallPath)
    }
    
    // 列出所有已安装版本
    fmt.Println("\n已安装版本:")
    configs, err := manager.List()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, config := range configs {
        fmt.Printf("  %s - %s\n", config.Version, config.InstallPath)
    }
}
```

## 使用便携版安装

### 为特定版本创建客户端

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
    
    // 为Node.js 18.17.0创建客户端
    client18, err := manager.CreateClient("18.17.0")
    if err != nil {
        log.Fatal(err)
    }
    
    // 为Node.js 20.5.0创建客户端
    client20, err := manager.CreateClient("20.5.0")
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用不同版本
    version18, err := client18.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("客户端18的npm版本: %s\n", version18)
    
    version20, err := client20.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("客户端20的npm版本: %s\n", version20)
    
    // 使用不同版本安装包
    fmt.Println("使用Node.js 18安装lodash...")
    err = client18.InstallPackage(ctx, "lodash", npm.InstallOptions{
        WorkingDir: "/tmp/project18",
    })
    if err != nil {
        log.Printf("失败: %v", err)
    }
    
    fmt.Println("使用Node.js 20安装lodash...")
    err = client20.InstallPackage(ctx, "lodash", npm.InstallOptions{
        WorkingDir: "/tmp/project20",
    })
    if err != nil {
        log.Printf("失败: %v", err)
    }
}
```

### 版本切换

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
    
    // 测试版本的函数
    testVersion := func(version string) {
        fmt.Printf("\n=== 测试Node.js %s ===\n", version)
        
        client, err := manager.CreateClient(version)
        if err != nil {
            log.Printf("为%s创建客户端失败: %v", version, err)
            return
        }
        
        // 检查npm是否可用
        if !client.IsAvailable(ctx) {
            log.Printf("版本%s的npm不可用", version)
            return
        }
        
        // 获取npm版本
        npmVersion, err := client.Version(ctx)
        if err != nil {
            log.Printf("获取npm版本失败: %v", err)
            return
        }
        
        fmt.Printf("npm版本: %s\n", npmVersion)
        
        // 创建测试项目
        tempDir := fmt.Sprintf("/tmp/test-project-%s", version)
        err = client.Init(ctx, npm.InitOptions{
            Name:       fmt.Sprintf("test-project-%s", version),
            Version:    "1.0.0",
            WorkingDir: tempDir,
        })
        if err != nil {
            log.Printf("初始化项目失败: %v", err)
            return
        }
        
        fmt.Printf("测试项目创建在%s\n", tempDir)
    }
    
    // 测试不同版本
    versions := []string{"16.20.0", "18.17.0", "20.5.0"}
    for _, version := range versions {
        testVersion(version)
    }
}
```

## 高级便携版管理

### 带进度的自定义安装

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
    
    // 带进度回调的安装
    version := "18.17.0"
    fmt.Printf("正在安装Node.js %s，带进度跟踪...\n", version)
    
    // 注意：这是概念示例 - 实际实现可能有所不同
    config, err := manager.Install(ctx, version)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("安装完成！\n")
    fmt.Printf("安装在: %s\n", config.InstallPath)
    fmt.Printf("安装日期: %s\n", config.InstallDate)
}
```

### 清理旧版本

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
    
    // 列出所有已安装版本
    configs, err := manager.List()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("当前已安装版本:")
    for _, config := range configs {
        fmt.Printf("  %s (安装于: %s)\n", config.Version, config.InstallDate)
    }
    
    // 删除旧版本（只保留最新的2个）
    if len(configs) > 2 {
        versionsToRemove := configs[:len(configs)-2]
        
        for _, config := range versionsToRemove {
            fmt.Printf("正在删除Node.js %s...\n", config.Version)
            err = manager.Uninstall(config.Version)
            if err != nil {
                log.Printf("删除%s失败: %v", config.Version, err)
            } else {
                fmt.Printf("Node.js %s删除成功\n", config.Version)
            }
        }
    }
    
    // 列出剩余版本
    fmt.Println("\n剩余版本:")
    configs, err = manager.List()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, config := range configs {
        fmt.Printf("  %s\n", config.Version)
    }
}
```

## 与CI/CD集成

### Docker式环境

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
    // 为CI/CD设置便携版环境
    ciDir := "/ci/npm-portable"
    
    // 清理任何现有安装
    os.RemoveAll(ciDir)
    os.MkdirAll(ciDir, 0755)
    
    manager, err := npm.NewPortableManager(ciDir)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 为CI安装特定Node.js版本
    nodeVersion := "18.17.0"
    fmt.Printf("正在设置CI环境，使用Node.js %s...\n", nodeVersion)
    
    config, err := manager.Install(ctx, nodeVersion)
    if err != nil {
        log.Fatal(err)
    }
    
    // 为CI操作创建客户端
    client, err := manager.CreateClient(nodeVersion)
    if err != nil {
        log.Fatal(err)
    }
    
    // 验证安装
    npmVersion, err := client.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("CI环境就绪！\n")
    fmt.Printf("Node.js: %s\n", nodeVersion)
    fmt.Printf("npm: %s\n", npmVersion)
    fmt.Printf("安装路径: %s\n", config.InstallPath)
    
    // 运行CI任务
    projectDir := "/ci/project"
    os.MkdirAll(projectDir, 0755)
    
    // 初始化项目
    err = client.Init(ctx, npm.InitOptions{
        Name:       "ci-project",
        Version:    "1.0.0",
        WorkingDir: projectDir,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 安装依赖
    dependencies := []string{"express", "lodash", "axios"}
    for _, dep := range dependencies {
        fmt.Printf("正在安装%s...\n", dep)
        err = client.InstallPackage(ctx, dep, npm.InstallOptions{
            WorkingDir: projectDir,
        })
        if err != nil {
            log.Printf("安装%s失败: %v", dep, err)
        }
    }
    
    // 安装开发依赖
    devDeps := []string{"jest", "eslint", "typescript"}
    for _, dep := range devDeps {
        fmt.Printf("正在安装%s作为开发依赖...\n", dep)
        err = client.InstallPackage(ctx, dep, npm.InstallOptions{
            WorkingDir: projectDir,
            SaveDev:    true,
        })
        if err != nil {
            log.Printf("安装%s失败: %v", dep, err)
        }
    }
    
    fmt.Println("CI环境设置完成！")
}
```

## 平台特定便携版安装

### 跨平台设置

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
    // 检测当前平台
    detector := platform.NewDetector()
    info, err := detector.Detect()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("检测到平台: %s %s\n", info.Platform, info.Architecture)
    
    // 设置平台特定的便携版目录
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
    
    // 为平台安装适当的Node.js版本
    version := "18.17.0"
    fmt.Printf("正在为%s安装Node.js %s...\n", info.Platform, version)
    
    config, err := manager.Install(ctx, version)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("安装成功！\n")
    fmt.Printf("平台: %s\n", info.Platform)
    fmt.Printf("架构: %s\n", info.Architecture)
    fmt.Printf("安装路径: %s\n", config.InstallPath)
    fmt.Printf("Node路径: %s\n", config.NodePath)
    fmt.Printf("NPM路径: %s\n", config.NpmPath)
}
```

## 最佳实践

1. **使用特定版本**: 始终为可重现环境指定确切的Node.js版本
2. **按用途组织**: 为不同项目或环境使用不同的便携版目录
3. **定期清理**: 删除未使用的版本以节省磁盘空间
4. **版本测试**: 使用多个Node.js版本测试您的应用程序
5. **CI/CD集成**: 使用便携版安装实现一致的CI/CD环境
6. **平台感知**: 考虑平台特定的路径和行为
7. **备份配置**: 跟踪哪些版本适用于您的项目

## 下一步

- [高级功能示例](./advanced-features.md)
- [基本用法示例](./basic-usage.md)
