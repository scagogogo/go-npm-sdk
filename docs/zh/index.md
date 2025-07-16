---
layout: home

hero:
  name: "Go NPM SDK"
  text: "全面的Go语言npm操作SDK"
  tagline: "跨平台npm管理，支持自动安装、便携版本和完整API覆盖"
  image:
    src: /logo.svg
    alt: Go NPM SDK
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 查看GitHub
      link: https://github.com/scagogogo/go-npm-sdk

features:
  - icon: 🚀
    title: 自动npm安装
    details: 根据操作系统自动检测和安装npm，支持包管理器或官方安装程序。
  
  - icon: 📦
    title: 便携版支持
    details: 下载和管理便携版Node.js/npm，无需系统级安装。
  
  - icon: 🔧
    title: 完整API封装
    details: 完整封装所有常用npm命令，包括安装、卸载、更新、发布等。
  
  - icon: 🌍
    title: 跨平台支持
    details: 在Windows、macOS和Linux上无缝工作，具有平台特定优化。
  
  - icon: 📝
    title: 项目管理
    details: 读取、写入和管理package.json文件，提供全面的依赖管理。
  
  - icon: ⚡
    title: 高性能
    details: 异步执行，支持超时控制、流式输出和批量操作。
  
  - icon: 🛡️
    title: 类型安全
    details: 全面的错误处理，具有结构化错误类型和验证。
  
  - icon: 🧪
    title: 充分测试
    details: 广泛的测试覆盖率（69.7%），包含全面的单元和集成测试。
  
  - icon: 📚
    title: 丰富文档
    details: 完整的API文档，包含示例和最佳实践。
---

## 快速开始

安装SDK：

```bash
go get github.com/scagogogo/go-npm-sdk
```

基本用法：

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

## 为什么选择Go NPM SDK？

- **零配置**: 开箱即用，自动npm检测和安装
- **生产就绪**: 在生产环境中使用，具有全面的错误处理
- **开发者友好**: 直观的API设计，丰富的文档和示例
- **积极维护**: 定期更新和社区支持

## 社区

- [GitHub Issues](https://github.com/scagogogo/go-npm-sdk/issues) - 报告错误和请求功能
- [GitHub Discussions](https://github.com/scagogogo/go-npm-sdk/discussions) - 提问和分享想法
- [贡献指南](https://github.com/scagogogo/go-npm-sdk/blob/main/CONTRIBUTING.md) - 了解如何贡献

## 许可证

基于[MIT许可证](https://github.com/scagogogo/go-npm-sdk/blob/main/LICENSE)发布。
