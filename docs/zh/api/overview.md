# API 概览

Go NPM SDK 为Go应用程序中的npm操作提供了一套全面的API。SDK被组织成几个包，每个包都有特定的用途。

## 包结构

```
github.com/scagogogo/go-npm-sdk/
├── pkg/npm/           # 核心npm操作
├── pkg/platform/      # 平台检测和下载
└── pkg/utils/         # 工具函数
```

## 核心包

### pkg/npm

包含npm客户端实现和相关功能的主要包：

- **客户端接口**: 所有操作的主要npm客户端
- **安装器**: 自动npm安装管理
- **检测器**: npm可用性检测
- **便携版管理器**: 便携版npm版本管理
- **包管理器**: package.json文件操作
- **依赖管理器**: 依赖解析和管理
- **类型**: 数据结构和接口
- **错误**: 错误类型和处理

### pkg/platform

平台特定功能：

- **检测器**: 操作系统和架构检测
- **下载器**: 带进度跟踪的文件下载功能

### pkg/utils

工具函数：

- **执行器**: 具有高级功能的命令执行

## 关键接口

### 客户端接口

npm操作的主要接口：

```go
type Client interface {
    // 基本操作
    IsAvailable(ctx context.Context) bool
    Install(ctx context.Context) error
    Version(ctx context.Context) (string, error)
    
    // 项目管理
    Init(ctx context.Context, options InitOptions) error
    
    // 包操作
    InstallPackage(ctx context.Context, pkg string, options InstallOptions) error
    UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error
    UpdatePackage(ctx context.Context, pkg string) error
    ListPackages(ctx context.Context, options ListOptions) ([]Package, error)
    
    // 脚本执行
    RunScript(ctx context.Context, script string, args ...string) error
    
    // 发布
    Publish(ctx context.Context, options PublishOptions) error
    
    // 信息检索
    GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)
    Search(ctx context.Context, query string) ([]SearchResult, error)
}
```

## 常见模式

### Context使用

所有操作都接受`context.Context`参数用于取消和超时控制：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.InstallPackage(ctx, "lodash", npm.InstallOptions{})
```

### 错误处理

SDK提供结构化错误类型以便更好地处理错误：

```go
if err != nil {
    if npm.IsNpmNotFound(err) {
        // 处理npm未找到
    } else if npm.IsPackageNotFound(err) {
        // 处理包未找到
    } else {
        // 处理其他错误
    }
}
```

### 选项模式

大多数操作使用选项结构体进行配置：

```go
options := npm.InstallOptions{
    SaveDev:    true,
    SaveExact:  true,
    WorkingDir: "/path/to/project",
    Registry:   "https://registry.npmjs.org/",
}

err := client.InstallPackage(ctx, "typescript", options)
```

## 快速开始

1. **导入包**:
   ```go
   import "github.com/scagogogo/go-npm-sdk/pkg/npm"
   ```

2. **创建客户端**:
   ```go
   client, err := npm.NewClient()
   if err != nil {
       log.Fatal(err)
   }
   ```

3. **使用客户端**:
   ```go
   ctx := context.Background()
   
   if !client.IsAvailable(ctx) {
       err := client.Install(ctx)
       if err != nil {
           log.Fatal(err)
       }
   }
   
   version, err := client.Version(ctx)
   if err != nil {
       log.Fatal(err)
   }
   fmt.Printf("npm版本: %s\n", version)
   ```

## 下一步

- [客户端接口](./client.md) - 详细的客户端API文档
- [NPM包](./npm.md) - 完整的npm包参考
- [平台包](./platform.md) - 平台检测和下载API
- [工具包](./utils.md) - 工具函数参考
- [类型与错误](./types-errors.md) - 数据类型和错误处理
