# 类型与错误API

本文档涵盖Go NPM SDK中使用的所有数据类型、结构和错误处理机制。

## 核心类型

### 客户端接口

npm操作的主要接口：

```go
type Client interface {
    IsAvailable(ctx context.Context) bool
    Install(ctx context.Context) error
    Version(ctx context.Context) (string, error)
    Init(ctx context.Context, options InitOptions) error
    InstallPackage(ctx context.Context, pkg string, options InstallOptions) error
    UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error
    UpdatePackage(ctx context.Context, pkg string) error
    ListPackages(ctx context.Context, options ListOptions) ([]Package, error)
    RunScript(ctx context.Context, script string, args ...string) error
    Publish(ctx context.Context, options PublishOptions) error
    GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)
    Search(ctx context.Context, query string) ([]SearchResult, error)
}
```

## 选项类型

### InitOptions

项目初始化配置：

```go
type InitOptions struct {
    Name        string `json:"name,omitempty"`
    Version     string `json:"version,omitempty"`
    Description string `json:"description,omitempty"`
    Author      string `json:"author,omitempty"`
    License     string `json:"license,omitempty"`
    Private     bool   `json:"private,omitempty"`
    WorkingDir  string `json:"-"`
    Force       bool   `json:"-"`
}
```

### InstallOptions

包安装配置：

```go
type InstallOptions struct {
    SaveDev       bool              `json:"save_dev,omitempty"`
    SaveOptional  bool              `json:"save_optional,omitempty"`
    SaveExact     bool              `json:"save_exact,omitempty"`
    Global        bool              `json:"global,omitempty"`
    Production    bool              `json:"production,omitempty"`
    WorkingDir    string            `json:"working_dir,omitempty"`
    Registry      string            `json:"registry,omitempty"`
    Force         bool              `json:"force,omitempty"`
    IgnoreScripts bool              `json:"ignore_scripts,omitempty"`
}
```

### UninstallOptions

包卸载配置：

```go
type UninstallOptions struct {
    SaveDev    bool   `json:"save_dev,omitempty"`
    Global     bool   `json:"global,omitempty"`
    WorkingDir string `json:"working_dir,omitempty"`
}
```

### ListOptions

包列表配置：

```go
type ListOptions struct {
    Global     bool   `json:"global,omitempty"`
    Depth      int    `json:"depth,omitempty"`
    Production bool   `json:"production,omitempty"`
    WorkingDir string `json:"working_dir,omitempty"`
    JSON       bool   `json:"json,omitempty"`
}
```

### PublishOptions

包发布配置：

```go
type PublishOptions struct {
    Tag        string `json:"tag,omitempty"`
    Access     string `json:"access,omitempty"`
    Registry   string `json:"registry,omitempty"`
    WorkingDir string `json:"working_dir,omitempty"`
    DryRun     bool   `json:"dry_run,omitempty"`
}
```

## 数据类型

### Package

表示npm包：

```go
type Package struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description,omitempty"`
    Dependencies map[string]string `json:"dependencies,omitempty"`
    DevDeps      map[string]string `json:"devDependencies,omitempty"`
    OptionalDeps map[string]string `json:"optionalDependencies,omitempty"`
    PeerDeps     map[string]string `json:"peerDependencies,omitempty"`
    Scripts      map[string]string `json:"scripts,omitempty"`
    Keywords     []string          `json:"keywords,omitempty"`
    Author       string            `json:"author,omitempty"`
    License      string            `json:"license,omitempty"`
    Homepage     string            `json:"homepage,omitempty"`
    Repository   *Repository       `json:"repository,omitempty"`
    Bugs         *Bugs             `json:"bugs,omitempty"`
    Main         string            `json:"main,omitempty"`
    Private      bool              `json:"private,omitempty"`
}
```

### Repository

仓库信息：

```go
type Repository struct {
    Type string `json:"type"`
    URL  string `json:"url"`
}
```

### Bugs

错误报告信息：

```go
type Bugs struct {
    URL   string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}
```

### PackageInfo

来自注册表的详细包信息：

```go
type PackageInfo struct {
    Name         string                 `json:"name"`
    Version      string                 `json:"version"`
    Description  string                 `json:"description,omitempty"`
    Keywords     []string               `json:"keywords,omitempty"`
    Homepage     string                 `json:"homepage,omitempty"`
    Repository   *Repository            `json:"repository,omitempty"`
    Author       *Person                `json:"author,omitempty"`
    License      string                 `json:"license,omitempty"`
    Dependencies map[string]string      `json:"dependencies,omitempty"`
    DevDeps      map[string]string      `json:"devDependencies,omitempty"`
    Versions     map[string]interface{} `json:"versions,omitempty"`
    Time         map[string]time.Time   `json:"time,omitempty"`
    DistTags     map[string]string      `json:"dist-tags,omitempty"`
}
```

### Person

人员信息（作者、维护者等）：

```go
type Person struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    URL   string `json:"url,omitempty"`
}
```

### SearchResult

来自npm注册表的搜索结果：

```go
type SearchResult struct {
    Package     SearchPackage `json:"package"`
    Score       SearchScore   `json:"score"`
    SearchScore float64       `json:"searchScore"`
}
```

### SearchPackage

搜索结果中的包信息：

```go
type SearchPackage struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Keywords    []string          `json:"keywords,omitempty"`
    Date        time.Time         `json:"date"`
    Links       map[string]string `json:"links,omitempty"`
    Author      *Person           `json:"author,omitempty"`
    Publisher   *Person           `json:"publisher,omitempty"`
    Maintainers []*Person         `json:"maintainers,omitempty"`
}
```

### SearchScore

搜索结果的评分信息：

```go
type SearchScore struct {
    Final  float64     `json:"final"`
    Detail ScoreDetail `json:"detail"`
}

type ScoreDetail struct {
    Quality     float64 `json:"quality"`
    Popularity  float64 `json:"popularity"`
    Maintenance float64 `json:"maintenance"`
}
```

### CommandResult

命令执行结果：

```go
type CommandResult struct {
    Success  bool          `json:"success"`
    ExitCode int           `json:"exit_code"`
    Stdout   string        `json:"stdout,omitempty"`
    Stderr   string        `json:"stderr,omitempty"`
    Duration time.Duration `json:"duration"`
    Error    error         `json:"-"`
}
```

## 安装类型

### NpmInstallOptions

npm安装选项：

```go
type NpmInstallOptions struct {
    Method      string                `json:"method"`
    Version     string                `json:"version,omitempty"`
    InstallPath string                `json:"install_path,omitempty"`
    Force       bool                  `json:"force,omitempty"`
    Global      bool                  `json:"global,omitempty"`
    Progress    func(message string)  `json:"-"`
}
```

### InstallResult

npm安装结果：

```go
type InstallResult struct {
    Success  bool          `json:"success"`
    Method   string        `json:"method"`
    Version  string        `json:"version"`
    Path     string        `json:"path"`
    Duration time.Duration `json:"duration"`
    Error    string        `json:"error,omitempty"`
}
```

### NpmInfo

npm安装信息：

```go
type NpmInfo struct {
    Version   string `json:"version"`
    Path      string `json:"path"`
    Available bool   `json:"available"`
}
```

### PortableConfig

便携版npm安装配置：

```go
type PortableConfig struct {
    Version     string `json:"version"`
    InstallPath string `json:"install_path"`
    NodePath    string `json:"node_path"`
    NpmPath     string `json:"npm_path"`
    InstallDate string `json:"install_date"`
}
```

## 错误类型

### NpmError

npm操作的基础错误类型：

```go
type NpmError struct {
    Op       string `json:"op"`
    Package  string `json:"package,omitempty"`
    ExitCode int    `json:"exit_code"`
    Stdout   string `json:"stdout,omitempty"`
    Stderr   string `json:"stderr,omitempty"`
    Err      error  `json:"-"`
}

func (e *NpmError) Error() string
func (e *NpmError) Unwrap() error
```

### ValidationError

验证失败的错误：

```go
type ValidationError struct {
    Field  string `json:"field"`
    Value  string `json:"value"`
    Reason string `json:"reason"`
}

func (e *ValidationError) Error() string
```

### InstallError

安装失败的错误：

```go
type InstallError struct {
    Package string `json:"package"`
    Reason  string `json:"reason"`
    Err     error  `json:"-"`
}

func (e *InstallError) Error() string
func (e *InstallError) Unwrap() error
```

### UninstallError

卸载失败的错误：

```go
type UninstallError struct {
    Package string `json:"package"`
    Reason  string `json:"reason"`
    Err     error  `json:"-"`
}

func (e *UninstallError) Error() string
func (e *UninstallError) Unwrap() error
```

### PlatformError

平台相关问题的错误：

```go
type PlatformError struct {
    Platform string `json:"platform"`
    Reason   string `json:"reason"`
    Err      error  `json:"-"`
}

func (e *PlatformError) Error() string
func (e *PlatformError) Unwrap() error
```

### DownloadError

下载失败的错误：

```go
type DownloadError struct {
    URL    string `json:"url"`
    Reason string `json:"reason"`
    Err    error  `json:"-"`
}

func (e *DownloadError) Error() string
func (e *DownloadError) Unwrap() error
```

## 错误常量

预定义的错误常量：

```go
var (
    ErrNpmNotFound         = errors.New("npm not found")
    ErrPackageNotFound     = errors.New("package not found")
    ErrInvalidVersion      = errors.New("invalid version")
    ErrInvalidPackageName  = errors.New("invalid package name")
    ErrPermissionDenied    = errors.New("permission denied")
    ErrUnsupportedPlatform = errors.New("unsupported platform")
)
```

## 错误检查函数

错误类型检查的辅助函数：

```go
func IsNpmNotFound(err error) bool
func IsPackageNotFound(err error) bool
func IsDownloadError(err error) bool
func IsPlatformError(err error) bool
```

**使用示例:**
```go
err := client.InstallPackage(ctx, "nonexistent-package", npm.InstallOptions{})
if err != nil {
    if npm.IsPackageNotFound(err) {
        fmt.Println("注册表中未找到包")
    } else if npm.IsNpmNotFound(err) {
        fmt.Println("npm未安装")
    } else if npmErr, ok := err.(*npm.NpmError); ok {
        fmt.Printf("npm命令失败: %s\n", npmErr.Stderr)
    } else {
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

## 构造函数

创建错误类型的辅助函数：

```go
func NewNpmError(op, pkg string, exitCode int, stdout, stderr string, err error) *NpmError
func NewValidationError(field, value, reason string) *ValidationError
func NewInstallError(pkg, reason string, err error) *InstallError
func NewUninstallError(pkg, reason string, err error) *UninstallError
func NewPlatformError(platform, reason string, err error) *PlatformError
func NewDownloadError(url, reason string, err error) *DownloadError
```

## 最佳实践

1. **错误处理**: 始终使用提供的辅助函数检查特定错误类型
2. **验证**: 对输入验证失败使用ValidationError
3. **上下文**: 在错误消息中包含相关上下文信息
4. **包装**: 使用错误包装来保留原始错误链
5. **结构化数据**: 使用结构化错误类型而不是纯字符串以便更好地处理错误
6. **JSON序列化**: 大多数类型支持JSON序列化以用于API响应
7. **空值检查**: 在处理可选字段时始终检查空指针
