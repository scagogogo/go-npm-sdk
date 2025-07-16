# Go NPM SDK

一个用于在Go语言中操作npm的SDK，提供了npm常用操作的Go API封装。

## 特性

- 🚀 **自动npm安装**: 根据操作系统自动检测和安装npm
- 📦 **便携版支持**: 支持下载便携版Node.js/npm
- 🔧 **完整API封装**: 封装npm的所有常用命令
- 🌍 **跨平台支持**: 支持Windows、macOS、Linux
- 📝 **项目管理**: 提供package.json读写和依赖管理功能
- ⚡ **高性能**: 异步执行，支持超时控制

## 快速开始

### 安装

```bash
go get github.com/scagogogo/go-npm-sdk
```

### 基本使用

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
    
    // 检查npm是否可用
    if !client.IsAvailable(context.Background()) {
        // 自动安装npm
        if err := client.Install(context.Background()); err != nil {
            log.Fatal(err)
        }
    }
    
    // 初始化项目
    if err := client.Init(context.Background(), "my-project"); err != nil {
        log.Fatal(err)
    }
    
    // 安装依赖
    if err := client.InstallPackage(context.Background(), "lodash"); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("项目初始化完成！")
}
```

## API文档

### 核心接口

#### Client

```go
type Client interface {
    // 检查npm是否可用
    IsAvailable(ctx context.Context) bool
    
    // 安装npm
    Install(ctx context.Context) error
    
    // 获取npm版本
    Version(ctx context.Context) (string, error)
    
    // 项目初始化
    Init(ctx context.Context, name string) error
    
    // 安装包
    InstallPackage(ctx context.Context, pkg string) error
    
    // 卸载包
    UninstallPackage(ctx context.Context, pkg string) error
    
    // 更新包
    UpdatePackage(ctx context.Context, pkg string) error
    
    // 列出已安装的包
    ListPackages(ctx context.Context) ([]Package, error)
    
    // 运行脚本
    RunScript(ctx context.Context, script string) error
}
```

### 项目管理

```go
// 读取package.json
pkg, err := npm.ReadPackageJSON("./package.json")

// 添加依赖
pkg.AddDependency("lodash", "^4.17.21")

// 保存package.json
err = pkg.Save("./package.json")
```

## 支持的操作系统

- **Windows**: 通过Chocolatey或官方安装程序安装
- **macOS**: 通过Homebrew或官方安装程序安装  
- **Linux**: 通过包管理器（apt、yum、pacman等）安装

## 开发

### 构建

```bash
go build ./...
```

### 测试

```bash
go test ./...
```

### 运行示例

```bash
go run examples/basic/main.go
```

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件
