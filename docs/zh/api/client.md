# 客户端接口

`Client`接口是Go NPM SDK中所有npm操作的主要入口点。它提供了管理npm安装、包和项目的全面方法集。

## 创建客户端

### NewClient

使用默认配置创建新的npm客户端。

```go
func NewClient() (Client, error)
```

**示例:**
```go
client, err := npm.NewClient()
if err != nil {
    log.Fatal(err)
}
```

### NewClientWithPath

使用特定的npm可执行文件路径创建新的npm客户端。

```go
func NewClientWithPath(npmPath string) (Client, error)
```

**参数:**
- `npmPath` (string): npm可执行文件的路径

**示例:**
```go
client, err := npm.NewClientWithPath("/usr/local/bin/npm")
if err != nil {
    log.Fatal(err)
}
```

## 基本操作

### IsAvailable

检查npm是否可用并可以执行。

```go
IsAvailable(ctx context.Context) bool
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文

**返回:**
- `bool`: 如果npm可用返回true，否则返回false

**示例:**
```go
ctx := context.Background()
if client.IsAvailable(ctx) {
    fmt.Println("npm可用")
} else {
    fmt.Println("npm不可用")
}
```

### Install

如果系统上没有npm，则自动安装npm。

```go
Install(ctx context.Context) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文

**返回:**
- `error`: 如果安装失败返回错误

**示例:**
```go
ctx := context.Background()
if !client.IsAvailable(ctx) {
    err := client.Install(ctx)
    if err != nil {
        log.Fatalf("安装npm失败: %v", err)
    }
    fmt.Println("npm安装成功")
}
```

### Version

获取当前npm版本。

```go
Version(ctx context.Context) (string, error)
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文

**返回:**
- `string`: npm版本字符串
- `error`: 如果版本获取失败返回错误

**示例:**
```go
ctx := context.Background()
version, err := client.Version(ctx)
if err != nil {
    log.Fatalf("获取npm版本失败: %v", err)
}
fmt.Printf("npm版本: %s\n", version)
```

## 项目管理

### Init

使用package.json初始化新的npm项目。

```go
Init(ctx context.Context, options InitOptions) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `options` (InitOptions): 初始化选项

**返回:**
- `error`: 如果初始化失败返回错误

**示例:**
```go
ctx := context.Background()
options := npm.InitOptions{
    Name:        "my-project",
    Version:     "1.0.0",
    Description: "我的项目",
    Author:      "张三",
    License:     "MIT",
    WorkingDir:  "/path/to/project",
}

err := client.Init(ctx, options)
if err != nil {
    log.Fatalf("初始化项目失败: %v", err)
}
```

## 包操作

### InstallPackage

安装特定的npm包。

```go
InstallPackage(ctx context.Context, pkg string, options InstallOptions) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `pkg` (string): 要安装的包名
- `options` (InstallOptions): 安装选项

**返回:**
- `error`: 如果安装失败返回错误

**示例:**
```go
ctx := context.Background()
options := npm.InstallOptions{
    SaveDev:    true,
    SaveExact:  true,
    WorkingDir: "/path/to/project",
}

err := client.InstallPackage(ctx, "typescript", options)
if err != nil {
    log.Fatalf("安装包失败: %v", err)
}
```

### UninstallPackage

卸载特定的npm包。

```go
UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `pkg` (string): 要卸载的包名
- `options` (UninstallOptions): 卸载选项

**返回:**
- `error`: 如果卸载失败返回错误

**示例:**
```go
ctx := context.Background()
options := npm.UninstallOptions{
    SaveDev:    true,
    WorkingDir: "/path/to/project",
}

err := client.UninstallPackage(ctx, "typescript", options)
if err != nil {
    log.Fatalf("卸载包失败: %v", err)
}
```

### UpdatePackage

将特定的npm包更新到最新版本。

```go
UpdatePackage(ctx context.Context, pkg string) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `pkg` (string): 要更新的包名

**返回:**
- `error`: 如果更新失败返回错误

**示例:**
```go
ctx := context.Background()
err := client.UpdatePackage(ctx, "lodash")
if err != nil {
    log.Fatalf("更新包失败: %v", err)
}
```

### ListPackages

列出当前项目中已安装的包。

```go
ListPackages(ctx context.Context, options ListOptions) ([]Package, error)
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `options` (ListOptions): 列表选项

**返回:**
- `[]Package`: 已安装包的列表
- `error`: 如果列表获取失败返回错误

**示例:**
```go
ctx := context.Background()
options := npm.ListOptions{
    Global:     false,
    Depth:      1,
    WorkingDir: "/path/to/project",
    JSON:       true,
}

packages, err := client.ListPackages(ctx, options)
if err != nil {
    log.Fatalf("列出包失败: %v", err)
}

for _, pkg := range packages {
    fmt.Printf("包: %s@%s\n", pkg.Name, pkg.Version)
}
```

## 脚本执行

### RunScript

执行package.json中定义的npm脚本。

```go
RunScript(ctx context.Context, script string, args ...string) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `script` (string): 要执行的脚本名称
- `args` (...string): 传递给脚本的额外参数

**返回:**
- `error`: 如果脚本执行失败返回错误

**示例:**
```go
ctx := context.Background()

// 运行构建脚本
err := client.RunScript(ctx, "build")
if err != nil {
    log.Fatalf("运行构建脚本失败: %v", err)
}

// 运行带参数的测试脚本
err = client.RunScript(ctx, "test", "--verbose", "--coverage")
if err != nil {
    log.Fatalf("运行测试脚本失败: %v", err)
}
```

## 发布

### Publish

将当前包发布到npm注册表。

```go
Publish(ctx context.Context, options PublishOptions) error
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `options` (PublishOptions): 发布选项

**返回:**
- `error`: 如果发布失败返回错误

**示例:**
```go
ctx := context.Background()
options := npm.PublishOptions{
    Tag:        "latest",
    Access:     "public",
    Registry:   "https://registry.npmjs.org/",
    WorkingDir: "/path/to/project",
    DryRun:     false,
}

err := client.Publish(ctx, options)
if err != nil {
    log.Fatalf("发布包失败: %v", err)
}
```

## 信息检索

### GetPackageInfo

从npm注册表检索包的详细信息。

```go
GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `pkg` (string): 要获取信息的包名

**返回:**
- `*PackageInfo`: 详细的包信息
- `error`: 如果信息检索失败返回错误

**示例:**
```go
ctx := context.Background()
info, err := client.GetPackageInfo(ctx, "lodash")
if err != nil {
    log.Fatalf("获取包信息失败: %v", err)
}

fmt.Printf("包: %s@%s\n", info.Name, info.Version)
fmt.Printf("描述: %s\n", info.Description)
fmt.Printf("主页: %s\n", info.Homepage)
```

### Search

在npm注册表中搜索包。

```go
Search(ctx context.Context, query string) ([]SearchResult, error)
```

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `query` (string): 搜索查询

**返回:**
- `[]SearchResult`: 搜索结果列表
- `error`: 如果搜索失败返回错误

**示例:**
```go
ctx := context.Background()
results, err := client.Search(ctx, "react hooks")
if err != nil {
    log.Fatalf("搜索包失败: %v", err)
}

for _, result := range results {
    fmt.Printf("包: %s@%s\n", result.Package.Name, result.Package.Version)
    fmt.Printf("描述: %s\n", result.Package.Description)
    fmt.Printf("评分: %.2f\n", result.Score.Final)
    fmt.Println("---")
}
```

## 错误处理

客户端方法返回结构化错误，可以检查特定条件：

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

## 最佳实践

1. **始终使用context**: 传递适当的context以进行超时和取消控制
2. **检查可用性**: 在执行操作之前使用`IsAvailable()`
3. **处理错误**: 检查特定错误类型以便更好地处理错误
4. **使用选项**: 使用选项结构体配置操作以获得灵活性
5. **工作目录**: 在选项中设置适当的工作目录以进行项目特定操作
