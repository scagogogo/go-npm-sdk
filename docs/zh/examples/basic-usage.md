# 基本用法示例

本页面提供Go NPM SDK的基本用法示例。

## 快速开始

### 简单npm客户端

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    // 创建npm客户端
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 检查npm是否可用
    if !client.IsAvailable(ctx) {
        fmt.Println("未找到npm，正在安装...")
        if err := client.Install(ctx); err != nil {
            log.Fatal(err)
        }
        fmt.Println("npm安装成功！")
    }
    
    // 获取npm版本
    version, err := client.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("npm版本: %s\n", version)
}
```

### 安装包

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
    
    // 安装生产依赖
    fmt.Println("正在安装lodash...")
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        SaveDev:   false,
        SaveExact: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("lodash安装成功！")
    
    // 安装开发依赖
    fmt.Println("正在安装jest...")
    err = client.InstallPackage(ctx, "jest", npm.InstallOptions{
        SaveDev: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("jest安装成功！")
    
    // 安装全局包
    fmt.Println("正在全局安装typescript...")
    err = client.InstallPackage(ctx, "typescript", npm.InstallOptions{
        Global: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("typescript全局安装成功！")
}
```

### 项目初始化

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
    
    // 初始化新项目
    options := npm.InitOptions{
        Name:        "my-awesome-project",
        Version:     "1.0.0",
        Description: "一个很棒的Node.js项目",
        Author:      "您的姓名 <your.email@example.com>",
        License:     "MIT",
        Private:     false,
        WorkingDir:  "./my-project",
    }
    
    fmt.Println("正在初始化项目...")
    err = client.Init(ctx, options)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("项目初始化成功！")
}
```

### 运行脚本

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
    
    // 运行构建脚本
    fmt.Println("正在运行构建脚本...")
    err = client.RunScript(ctx, "build")
    if err != nil {
        log.Printf("构建失败: %v", err)
    } else {
        fmt.Println("构建成功完成！")
    }
    
    // 运行带参数的测试脚本
    fmt.Println("正在运行测试...")
    err = client.RunScript(ctx, "test", "--verbose", "--coverage")
    if err != nil {
        log.Printf("测试失败: %v", err)
    } else {
        fmt.Println("测试通过！")
    }
}
```

### 列出包

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
    
    // 列出本地包
    fmt.Println("本地包:")
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
    
    // 列出全局包
    fmt.Println("\n全局包:")
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

### 错误处理

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
    
    // 尝试安装不存在的包
    err = client.InstallPackage(ctx, "this-package-does-not-exist-12345", npm.InstallOptions{})
    if err != nil {
        // 检查特定错误类型
        if npm.IsPackageNotFound(err) {
            fmt.Println("注册表中未找到包")
        } else if npm.IsNpmNotFound(err) {
            fmt.Println("npm未安装")
        } else if npmErr, ok := err.(*npm.NpmError); ok {
            fmt.Printf("npm命令失败: %s\n", npmErr.Stderr)
        } else {
            fmt.Printf("未知错误: %v\n", err)
        }
    }
}
```

### 使用Context进行超时控制

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
    
    // 创建带超时的context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 带超时安装包
    fmt.Println("正在安装包，30秒超时...")
    err = client.InstallPackage(ctx, "express", npm.InstallOptions{})
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("安装超时")
        } else {
            log.Fatal(err)
        }
    } else {
        fmt.Println("包安装成功！")
    }
}
```

### 在不同目录中工作

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
    
    // 创建临时目录
    tempDir, err := os.MkdirTemp("", "npm-example")
    if err != nil {
        log.Fatal(err)
    }
    defer os.RemoveAll(tempDir)
    
    fmt.Printf("在目录中工作: %s\n", tempDir)
    
    // 在特定目录中初始化项目
    err = client.Init(ctx, npm.InitOptions{
        Name:       "temp-project",
        Version:    "1.0.0",
        WorkingDir: tempDir,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 在特定目录中安装包
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        WorkingDir: tempDir,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("项目已创建并在临时目录中安装了包！")
}
```

## 最佳实践

1. **始终使用context**: 传递适当的context以进行超时和取消控制
2. **首先检查可用性**: 在执行操作之前使用`IsAvailable()`
3. **正确处理错误**: 使用错误类型检查进行特定错误处理
4. **设置工作目录**: 为项目特定操作指定工作目录
5. **使用超时**: 为长时间运行的操作设置合理的超时
6. **清理资源**: 确保正确清理临时文件和目录

## 下一步

- [包管理示例](./package-management.md)
- [便携版安装示例](./portable-installation.md)
- [高级功能示例](./advanced-features.md)
