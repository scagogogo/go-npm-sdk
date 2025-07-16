# 快速开始

欢迎使用Go NPM SDK！本指南将帮助您快速上手使用SDK。

## 安装

使用Go模块安装SDK：

```bash
go get github.com/scagogogo/go-npm-sdk
```

## 快速开始

这是一个简单的示例来帮助您开始：

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
    
    // 安装包
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        SaveDev: false,
        SaveExact: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("lodash安装成功！")
}
```

## 核心概念

### 客户端接口

`Client`接口是所有npm操作的主要入口点。它提供以下方法：

- **包管理**: 安装、卸载、更新包
- **项目管理**: 初始化项目、管理package.json
- **脚本执行**: 运行npm脚本
- **信息检索**: 获取包信息、搜索包
- **发布**: 发布包到注册表

### Context使用

所有操作都接受`context.Context`参数用于：

- **超时控制**: 设置操作超时
- **取消**: 取消长时间运行的操作
- **请求跟踪**: 添加请求元数据

```go
// 为操作设置超时
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.InstallPackage(ctx, "express", npm.InstallOptions{})
```

### 错误处理

SDK提供结构化错误类型以便更好地处理错误：

```go
err := client.InstallPackage(ctx, "nonexistent-package", npm.InstallOptions{})
if err != nil {
    if npm.IsPackageNotFound(err) {
        fmt.Println("包未找到")
    } else if npm.IsNpmNotFound(err) {
        fmt.Println("npm未安装")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 基本操作

### 检查npm可用性

```go
client, _ := npm.NewClient()
ctx := context.Background()

if client.IsAvailable(ctx) {
    fmt.Println("npm可用")
} else {
    fmt.Println("npm不可用")
}
```

### 自动安装npm

```go
if !client.IsAvailable(ctx) {
    err := client.Install(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("npm安装成功")
}
```

### 获取npm版本

```go
version, err := client.Version(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("npm版本: %s\n", version)
```

### 初始化项目

```go
options := npm.InitOptions{
    Name:        "my-project",
    Version:     "1.0.0",
    Description: "我的项目",
    Author:      "您的姓名",
    License:     "MIT",
    WorkingDir:  "/path/to/project",
}

err := client.Init(ctx, options)
if err != nil {
    log.Fatal(err)
}
```

### 安装包

```go
// 安装生产依赖
err := client.InstallPackage(ctx, "express", npm.InstallOptions{
    SaveDev:   false,
    SaveExact: true,
})

// 安装开发依赖
err = client.InstallPackage(ctx, "jest", npm.InstallOptions{
    SaveDev: true,
})

// 安装全局包
err = client.InstallPackage(ctx, "typescript", npm.InstallOptions{
    Global: true,
})
```

### 运行脚本

```go
// 运行构建脚本
err := client.RunScript(ctx, "build")

// 运行带参数的测试脚本
err = client.RunScript(ctx, "test", "--verbose", "--coverage")
```

### 列出包

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

### 搜索包

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

## 高级功能

### 便携版npm管理

使用便携版npm而无需系统级安装：

```go
import "github.com/scagogogo/go-npm-sdk/pkg/npm"

manager, err := npm.NewPortableManager("/opt/npm-portable")
if err != nil {
    log.Fatal(err)
}

// 安装Node.js 18.17.0和npm
config, err := manager.Install(ctx, "18.17.0")
if err != nil {
    log.Fatal(err)
}

// 为此版本创建客户端
client, err := manager.CreateClient("18.17.0")
if err != nil {
    log.Fatal(err)
}

// 正常使用客户端
version, _ := client.Version(ctx)
fmt.Printf("便携版npm版本: %s\n", version)
```

### Package.json管理

直接操作package.json文件：

```go
import "github.com/scagogogo/go-npm-sdk/pkg/npm"

pkg := npm.NewPackageJSON("./package.json")

// 加载现有package.json
err := pkg.Load()
if err != nil {
    log.Fatal(err)
}

// 修改包信息
pkg.SetName("my-package")
pkg.SetVersion("2.0.0")
pkg.AddDependency("lodash", "^4.17.21")
pkg.AddScript("build", "webpack")

// 保存更改
err = pkg.Save()
if err != nil {
    log.Fatal(err)
}
```

### 平台检测

检测当前平台以进行平台特定操作：

```go
import "github.com/scagogogo/go-npm-sdk/pkg/platform"

detector := platform.NewDetector()
info, err := detector.Detect()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("平台: %s\n", info.Platform)
fmt.Printf("架构: %s\n", info.Architecture)

if info.IsLinux() {
    fmt.Printf("Linux发行版: %s\n", info.Distribution)
}
```

## 配置

### 工作目录

为npm操作设置工作目录：

```go
options := npm.InstallOptions{
    WorkingDir: "/path/to/project",
}

err := client.InstallPackage(ctx, "express", options)
```

### 注册表配置

使用自定义npm注册表：

```go
options := npm.InstallOptions{
    Registry: "https://registry.npmjs.org/",
}

err := client.InstallPackage(ctx, "private-package", options)
```

### 超时配置

设置操作超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := client.InstallPackage(ctx, "large-package", npm.InstallOptions{})
```

## 最佳实践

1. **始终使用context**: 传递适当的context以进行超时和取消控制
2. **首先检查可用性**: 在执行操作之前使用`IsAvailable()`
3. **正确处理错误**: 使用错误类型检查进行特定错误处理
4. **设置工作目录**: 为项目特定操作指定工作目录
5. **使用结构化选项**: 使用选项结构体配置操作
6. **验证输入**: 始终验证包名和版本
7. **清理资源**: 确保正确清理临时文件和进程

## 下一步

- [安装指南](./installation.md) - 详细安装说明
- [配置](./configuration.md) - 高级配置选项
- [平台支持](./platform-support.md) - 平台特定信息
- [API参考](/zh/api/overview.md) - 完整API文档
- [示例](/zh/examples/basic-usage.md) - 更多示例和用例
