# 平台包API

`pkg/platform`包提供平台检测和文件下载功能，支持Windows、macOS和Linux的跨平台操作。

## 检测器

检测器组件识别当前操作系统、架构和Linux发行版。

### NewDetector

```go
func NewDetector() *Detector
```

创建新的平台检测器。

### Detect

```go
func (d *Detector) Detect() (*Info, error)
```

检测当前平台信息。

**返回:**
- `*Info`: 平台信息结构
- `error`: 如果检测失败返回错误

**示例:**
```go
detector := platform.NewDetector()
info, err := detector.Detect()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("平台: %s\n", info.Platform)
fmt.Printf("架构: %s\n", info.Architecture)
fmt.Printf("发行版: %s\n", info.Distribution)
fmt.Printf("版本: %s\n", info.Version)
fmt.Printf("内核: %s\n", info.Kernel)
```

## 平台信息

### Info结构

```go
type Info struct {
    Platform     Platform     `json:"platform"`
    Architecture Architecture `json:"architecture"`
    Distribution Distribution `json:"distribution,omitempty"`
    Version      string       `json:"version,omitempty"`
    Kernel       string       `json:"kernel,omitempty"`
}
```

### 平台方法

```go
func (i *Info) IsWindows() bool
func (i *Info) IsMacOS() bool
func (i *Info) IsLinux() bool
func (i *Info) IsX86() bool
func (i *Info) IsARM() bool
func (i *Info) String() string
```

**示例:**
```go
info, _ := detector.Detect()

if info.IsLinux() {
    fmt.Println("运行在Linux上")
    if info.Distribution == platform.Ubuntu {
        fmt.Println("检测到Ubuntu发行版")
    }
}

if info.IsX86() {
    fmt.Println("x86架构")
} else if info.IsARM() {
    fmt.Println("ARM架构")
}

fmt.Printf("平台字符串: %s\n", info.String())
```

## 平台常量

### 平台类型

```go
const (
    Windows Platform = "windows"
    MacOS   Platform = "darwin"
    Linux   Platform = "linux"
    Unknown Platform = "unknown"
)
```

### 架构类型

```go
const (
    AMD64 Architecture = "amd64"
    I386  Architecture = "386"
    ARM64 Architecture = "arm64"
    ARM   Architecture = "arm"
)
```

### Linux发行版

```go
const (
    Ubuntu        Distribution = "ubuntu"
    Debian        Distribution = "debian"
    CentOS        Distribution = "centos"
    RHEL          Distribution = "rhel"
    Fedora        Distribution = "fedora"
    SUSE          Distribution = "suse"
    Arch          Distribution = "arch"
    Alpine        Distribution = "alpine"
    UnknownDistro Distribution = "unknown"
)
```

## 下载器

下载器组件提供带进度跟踪和重试机制的文件下载功能。

### NewDownloader

```go
func NewDownloader() *Downloader
```

使用默认配置创建新的下载器。

### Download

```go
func (d *Downloader) Download(ctx context.Context, options DownloadOptions) (*DownloadResult, error)
```

使用指定选项下载文件。

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `options` (DownloadOptions): 下载配置

**返回:**
- `*DownloadResult`: 下载结果信息
- `error`: 如果下载失败返回错误

### DownloadWithRetry

```go
func (d *Downloader) DownloadWithRetry(ctx context.Context, options DownloadOptions, maxRetries int) (*DownloadResult, error)
```

使用重试机制下载文件。

**示例:**
```go
downloader := platform.NewDownloader()
ctx := context.Background()

options := platform.DownloadOptions{
    URL:         "https://nodejs.org/dist/v18.17.0/node-v18.17.0-linux-x64.tar.xz",
    Destination: "/tmp/node-v18.17.0-linux-x64.tar.xz",
    Timeout:     30 * time.Minute,
    UserAgent:   "Go-NPM-SDK/1.0",
    Headers: map[string]string{
        "Accept": "application/octet-stream",
    },
    Progress: func(downloaded, total int64) {
        percentage := float64(downloaded) / float64(total) * 100
        fmt.Printf("下载进度: %.2f%%\n", percentage)
    },
}

result, err := downloader.DownloadWithRetry(ctx, options, 3)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("下载了%d字节，耗时%v\n", result.Size, result.Duration)
```

### DownloadOptions

```go
type DownloadOptions struct {
    URL         string                           `json:"url"`
    Destination string                           `json:"destination"`
    Timeout     time.Duration                    `json:"timeout,omitempty"`
    UserAgent   string                           `json:"user_agent,omitempty"`
    Headers     map[string]string                `json:"headers,omitempty"`
    Progress    func(downloaded, total int64)    `json:"-"`
}
```

### DownloadResult

```go
type DownloadResult struct {
    FilePath string        `json:"file_path"`
    Size     int64         `json:"size"`
    Duration time.Duration `json:"duration"`
    Success  bool          `json:"success"`
    Error    string        `json:"error,omitempty"`
}
```

## Node.js下载器

专门用于Node.js发布版本的下载器，具有平台特定的URL生成功能。

### NewNodeJSDownloader

```go
func NewNodeJSDownloader() *NodeJSDownloader
```

创建新的Node.js下载器。

### GetDownloadURL

```go
func (n *NodeJSDownloader) GetDownloadURL(version string, platform Platform, arch Architecture) string
```

为特定Node.js版本和平台生成下载URL。

**参数:**
- `version` (string): Node.js版本 (例如 "18.17.0")
- `platform` (Platform): 目标平台
- `arch` (Architecture): 目标架构

**返回:**
- `string`: 指定版本和平台的下载URL

### GetLatestVersion

```go
func (n *NodeJSDownloader) GetLatestVersion(ctx context.Context) (string, error)
```

从官方API检索最新的Node.js版本。

### Download

```go
func (n *NodeJSDownloader) Download(ctx context.Context, version string, platform Platform, arch Architecture, destination string) (*DownloadResult, error)
```

为目标平台下载特定的Node.js版本。

**示例:**
```go
nodeDownloader := platform.NewNodeJSDownloader()
ctx := context.Background()

// 获取最新版本
latest, err := nodeDownloader.GetLatestVersion(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("最新Node.js版本: %s\n", latest)

// 获取下载URL
detector := platform.NewDetector()
info, _ := detector.Detect()

url := nodeDownloader.GetDownloadURL("18.17.0", info.Platform, info.Architecture)
fmt.Printf("下载URL: %s\n", url)

// 下载Node.js
result, err := nodeDownloader.Download(ctx, "18.17.0", info.Platform, info.Architecture, "/tmp/node.tar.xz")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("下载Node.js: %s (%d字节)\n", result.FilePath, result.Size)
```

## 平台特定URL

Node.js下载器生成平台特定的URL：

### Windows
- **AMD64**: `node-v{version}-win-x64.zip`
- **I386**: `node-v{version}-win-x86.zip`

### macOS
- **AMD64**: `node-v{version}-darwin-x64.tar.gz`
- **ARM64**: `node-v{version}-darwin-arm64.tar.gz`

### Linux
- **AMD64**: `node-v{version}-linux-x64.tar.xz`
- **ARM64**: `node-v{version}-linux-arm64.tar.xz`
- **ARM**: `node-v{version}-linux-armv7l.tar.xz`
- **I386**: `node-v{version}-linux-x86.tar.xz`

## 错误处理

平台包提供特定的错误类型：

```go
type PlatformError struct {
    Platform string
    Reason   string
    Err      error
}

type DownloadError struct {
    URL    string
    Reason string
    Err    error
}
```

**示例:**
```go
result, err := downloader.Download(ctx, options)
if err != nil {
    if platformErr, ok := err.(*platform.PlatformError); ok {
        fmt.Printf("平台错误在%s: %s\n", platformErr.Platform, platformErr.Reason)
    } else if downloadErr, ok := err.(*platform.DownloadError); ok {
        fmt.Printf("下载错误%s: %s\n", downloadErr.URL, downloadErr.Reason)
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 最佳实践

1. **检查平台兼容性**: 在操作前始终验证平台支持
2. **使用适当的超时**: 为下载操作设置合理的超时
3. **处理进度**: 实现进度回调以获得更好的用户体验
4. **失败时重试**: 对网络操作使用重试机制
5. **验证URL**: 在尝试下载之前确保下载URL有效
6. **失败时清理**: 在失败时删除不完整的下载
7. **使用context**: 始终传递context以支持取消

## 集成示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/platform"
)

func main() {
    // 检测平台
    detector := platform.NewDetector()
    info, err := detector.Detect()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("检测到平台: %s\n", info.String())
    
    // 为当前平台下载Node.js
    nodeDownloader := platform.NewNodeJSDownloader()
    ctx := context.Background()
    
    result, err := nodeDownloader.Download(
        ctx,
        "18.17.0",
        info.Platform,
        info.Architecture,
        "/tmp/node.tar.xz",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("下载Node.js: %s (%d字节)\n", result.FilePath, result.Size)
}
```
