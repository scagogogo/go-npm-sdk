# 安装

本指南介绍安装和使用Go NPM SDK的不同方法。

## 前提条件

- Go 1.19或更高版本
- Git（用于克隆仓库）
- 互联网连接（用于下载依赖）

## 安装SDK

### 使用Go模块（推荐）

安装Go NPM SDK最简单的方法是使用Go模块：

```bash
go get github.com/scagogogo/go-npm-sdk
```

这将下载SDK及其所有依赖项。

### 使用Git克隆

您也可以直接克隆仓库：

```bash
git clone https://github.com/scagogogo/go-npm-sdk.git
cd go-npm-sdk
go mod download
```

## 验证安装

创建一个简单的测试文件来验证安装：

```go
// test.go
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
    
    if client.IsAvailable(ctx) {
        version, err := client.Version(ctx)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("npm版本: %s\n", version)
    } else {
        fmt.Println("npm不可用")
    }
}
```

运行测试：

```bash
go run test.go
```

## 下一步

- [快速开始](./getting-started.md) - 学习基础知识
- [配置](./configuration.md) - 配置SDK
- [平台支持](./platform-support.md) - 平台特定信息
