# Go NPM SDK 使用指南

## 项目概述

Go NPM SDK 是一个用于在Go语言中操作npm的完整SDK，提供了npm常用操作的Go API封装。

## 主要特性

- 🚀 **自动npm安装**: 根据操作系统自动检测和安装npm
- 📦 **便携版支持**: 支持下载便携版Node.js/npm
- 🔧 **完整API封装**: 封装npm的所有常用命令
- 🌍 **跨平台支持**: 支持Windows、macOS、Linux
- 📝 **项目管理**: 提供package.json读写和依赖管理功能
- ⚡ **高性能**: 异步执行，支持超时控制

## 项目结构

```
go-npm-sdk/
├── cmd/                    # 命令行工具示例
├── pkg/
│   ├── npm/               # 核心npm操作
│   │   ├── client.go      # npm客户端实现
│   │   ├── installer.go   # npm安装管理
│   │   ├── detector.go    # npm检测功能
│   │   ├── portable.go    # 便携版管理
│   │   ├── package.go     # package.json管理
│   │   ├── dependency.go  # 依赖管理
│   │   ├── types.go       # 数据类型定义
│   │   └── errors.go      # 错误类型定义
│   ├── platform/          # 平台相关
│   │   ├── detector.go    # 操作系统检测
│   │   └── downloader.go  # 下载器
│   └── utils/             # 工具函数
│       └── executor.go    # 命令执行器
├── examples/              # 使用示例
│   ├── basic/            # 基本使用示例
│   └── portable/         # 便携版使用示例
├── tests/                 # 测试文件
├── go.mod
├── go.sum
├── README.md
├── USAGE.md
└── LICENSE
```

## 快速开始

### 1. 安装

```bash
go get github.com/scagogogo/go-npm-sdk
```

### 2. 基本使用

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
        // 自动安装npm
        if err := client.Install(ctx); err != nil {
            log.Fatal(err)
        }
    }
    
    // 获取npm版本
    version, err := client.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("npm版本: %s\n", version)
}
```

## 核心功能

### 1. npm客户端操作

```go
// 创建客户端
client, err := npm.NewClient()

// 检查npm是否可用
available := client.IsAvailable(ctx)

// 获取npm版本
version, err := client.Version(ctx)

// 自动安装npm
err = client.Install(ctx)
```

### 2. 项目初始化

```go
options := npm.InitOptions{
    Name:        "my-project",
    Version:     "1.0.0",
    Description: "My awesome project",
    Author:      "Your Name",
    License:     "MIT",
    WorkingDir:  "/path/to/project",
    Force:       true,
}

err := client.Init(ctx, options)
```

### 3. 包管理

```go
// 安装包
installOptions := npm.InstallOptions{
    WorkingDir: "/path/to/project",
    SaveDev:    false,
}
err := client.InstallPackage(ctx, "lodash", installOptions)

// 卸载包
uninstallOptions := npm.UninstallOptions{
    WorkingDir: "/path/to/project",
}
err := client.UninstallPackage(ctx, "lodash", uninstallOptions)

// 更新包
err := client.UpdatePackage(ctx, "lodash")

// 列出已安装的包
listOptions := npm.ListOptions{
    WorkingDir: "/path/to/project",
    Depth:      0,
}
packages, err := client.ListPackages(ctx, listOptions)
```

### 4. 脚本执行

```go
// 运行npm脚本
err := client.RunScript(ctx, "test")

// 运行带参数的脚本
err := client.RunScript(ctx, "build", "--production")
```

### 5. 包信息查询

```go
// 获取包信息
packageInfo, err := client.GetPackageInfo(ctx, "lodash")

// 搜索包
results, err := client.Search(ctx, "react")
```

### 6. 发布包

```go
publishOptions := npm.PublishOptions{
    Tag:        "beta",
    Access:     "public",
    WorkingDir: "/path/to/project",
    DryRun:     false,
}
err := client.Publish(ctx, publishOptions)
```

## 高级功能

### 1. 便携版管理

```go
// 创建便携版管理器
portableManager, err := npm.NewPortableManager("/path/to/portable")

// 安装便携版Node.js
progress := func(message string) {
    fmt.Println(message)
}
config, err := portableManager.Install(ctx, "18.17.0", progress)

// 使用便携版创建客户端
client, err := portableManager.CreateClient("18.17.0")

// 列出已安装的版本
configs, err := portableManager.List()

// 设置为默认版本
err = portableManager.SetAsDefault("18.17.0")
```

### 2. package.json管理

```go
// 创建package.json管理器
packageJSON := npm.NewPackageJSON("/path/to/package.json")

// 加载现有文件
err := packageJSON.Load()

// 修改基本信息
packageJSON.SetName("my-package")
packageJSON.SetVersion("1.0.0")
packageJSON.SetDescription("My package")

// 管理依赖
packageJSON.AddDependency("lodash", "^4.17.21")
packageJSON.AddDevDependency("jest", "^27.0.0")

// 管理脚本
packageJSON.AddScript("test", "jest")
packageJSON.AddScript("build", "webpack")

// 保存文件
err = packageJSON.Save()
```

### 3. 依赖管理

```go
// 创建依赖管理器
depManager, err := npm.NewDependencyManager(client, "/path/to/project")

// 添加依赖
operation, err := depManager.Add(ctx, "lodash", "^4.17.21", npm.Production)

// 移除依赖
operation, err := depManager.Remove(ctx, "lodash")

// 更新依赖
operation, err := depManager.Update(ctx, "lodash")

// 列出所有依赖
dependencies, err := depManager.List(ctx)

// 检查过期依赖
outdated, err := depManager.CheckOutdated(ctx)

// 安装所有依赖
err = depManager.Install(ctx)
```

## 平台支持

### 支持的操作系统

- **Windows**: 通过Chocolatey、winget或官方安装程序
- **macOS**: 通过Homebrew、MacPorts或官方安装程序  
- **Linux**: 通过包管理器（apt、yum、pacman等）或官方安装程序

### 支持的架构

- x86_64 (amd64)
- ARM64
- x86 (386)
- ARM

## 错误处理

SDK提供了详细的错误类型：

```go
// 检查特定错误类型
if npm.IsNpmNotFound(err) {
    // npm未找到
}

if npm.IsPackageNotFound(err) {
    // 包未找到
}

if npm.IsNetworkError(err) {
    // 网络错误
}

// 获取详细错误信息
if npmErr, ok := err.(*npm.NpmError); ok {
    fmt.Printf("操作: %s, 退出码: %d\n", npmErr.Op, npmErr.ExitCode)
    fmt.Printf("输出: %s\n", npmErr.Stdout)
    fmt.Printf("错误: %s\n", npmErr.Stderr)
}
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/npm
go test ./pkg/platform

# 运行测试并显示覆盖率
go test -cover ./...
```

## 示例

查看 `examples/` 目录中的完整示例：

- `examples/basic/main.go` - 基本功能演示
- `examples/portable/main.go` - 便携版功能演示

运行示例：

```bash
go run examples/basic/main.go
go run examples/portable/main.go
```

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件
