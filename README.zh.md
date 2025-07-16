# Go NPM SDK

[English](README.md) | [简体中文](README.zh.md)

一个用于在Go语言中操作npm的全面SDK，支持跨平台操作。

## 📚 文档

**🌐 [完整文档网站](https://scagogogo.github.io/go-npm-sdk/)**

访问我们的综合文档网站，获取详细指南、API参考和示例。

## 特性

- **自动npm安装**: 根据操作系统自动检测和安装npm
- **便携版支持**: 下载和管理便携版Node.js/npm
- **完整API覆盖**: 完整封装所有常用npm命令
- **跨平台支持**: 支持Windows、macOS和Linux
- **项目管理**: 读取、写入和管理package.json文件
- **高性能**: 异步执行，支持超时控制
- **类型安全**: 全面的错误处理和结构化错误类型

## 安装

```bash
go get github.com/scagogogo/go-npm-sdk
```

## 快速开始

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
    
    // 安装包
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        SaveDev: false,
        SaveExact: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("包安装成功！")
}
```

> 💡 **需要更多示例？** 查看我们的[完整文档](https://scagogogo.github.io/go-npm-sdk/)获取详细指南和高级用法模式。

## 核心功能

### 自动npm安装

SDK可以根据您的操作系统自动检测和安装npm：

```go
client, _ := npm.NewClient()
ctx := context.Background()

if !client.IsAvailable(ctx) {
    // 使用最适合您操作系统的方法自动安装npm
    err := client.Install(ctx)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 便携版npm管理

下载和管理便携版Node.js/npm，无需系统级安装：

```go
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
```

### Package.json管理

读取、写入和管理package.json文件：

```go
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

## API文档

完整的API文档请访问我们的文档网站：

**📚 [完整文档网站](https://scagogogo.github.io/go-npm-sdk/)**

文档包括：
- 完整的API参考
- 使用指南和教程
- 示例和最佳实践
- 平台特定信息

## 示例

查看[examples](./examples/)目录获取更多综合示例：

- [基本用法](./examples/basic_usage.go) - SDK入门
- [包管理](./examples/package_management.go) - 安装和管理包
- [便携版安装](./examples/portable_installation.go) - 使用便携版npm
- [平台检测](./examples/platform_detection.go) - 检测平台信息
- [依赖管理](./examples/dependency_management.go) - 管理依赖

## 支持的平台

- **Windows**: Windows 10/11, Windows Server 2019/2022
- **macOS**: macOS 10.15+ (Intel和Apple Silicon)
- **Linux**: Ubuntu, Debian, CentOS, RHEL, Fedora, SUSE, Arch, Alpine

## 安装方法

SDK支持多种npm安装方法：

1. **包管理器**: 使用系统包管理器（apt、yum、brew等）
2. **官方安装程序**: 下载并运行官方Node.js安装程序
3. **便携版**: 下载便携版Node.js/npm
4. **手动**: 手动安装指导

## 系统要求

- Go 1.19或更高版本
- 互联网连接（用于下载npm/Node.js，如果尚未安装）

## 贡献

我们欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

## 许可证

本项目基于MIT许可证发布 - 详见[LICENSE](LICENSE)文件。

## 支持

- [GitHub Issues](https://github.com/scagogogo/go-npm-sdk/issues) - 报告错误和请求功能
- [GitHub Discussions](https://github.com/scagogogo/go-npm-sdk/discussions) - 提问和分享想法
- **[📖 文档网站](https://scagogogo.github.io/go-npm-sdk/)** - 完整文档和指南
