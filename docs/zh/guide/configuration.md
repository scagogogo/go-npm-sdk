# 配置

了解如何根据您的特定需求配置Go NPM SDK。

## 基本配置

### 客户端配置

```go
// 使用默认配置创建客户端
client, err := npm.NewClient()

// 使用特定npm路径创建客户端
client, err := npm.NewClientWithPath("/usr/local/bin/npm")
```

### 工作目录

为npm操作设置工作目录：

```go
options := npm.InstallOptions{
    WorkingDir: "/path/to/project",
}

err := client.InstallPackage(ctx, "express", options)
```

### 注册表配置

配置自定义npm注册表：

```go
options := npm.InstallOptions{
    Registry: "https://registry.npmjs.org/",
}

err := client.InstallPackage(ctx, "package", options)
```

## 高级配置

### 超时配置

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := client.InstallPackage(ctx, "large-package", npm.InstallOptions{})
```

### 环境变量

```go
import "github.com/scagogogo/go-npm-sdk/pkg/utils"

executor := utils.NewExecutor()
executor.SetDefaultEnv(map[string]string{
    "NODE_ENV": "production",
    "NPM_CONFIG_REGISTRY": "https://registry.npmjs.org/",
})
```

## 平台特定配置

### Windows配置

```go
// Windows特定的npm路径
client, err := npm.NewClientWithPath("C:\\Program Files\\nodejs\\npm.cmd")
```

### macOS配置

```go
// 使用Homebrew的macOS
client, err := npm.NewClientWithPath("/opt/homebrew/bin/npm")
```

### Linux配置

```go
// Linux系统npm
client, err := npm.NewClientWithPath("/usr/bin/npm")
```

## 下一步

- [平台支持](./platform-support.md) - 平台特定信息
- [API参考](/zh/api/overview.md) - 完整API文档
