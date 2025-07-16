# NPM包API

`pkg/npm`包提供全面的npm功能，包括客户端操作、安装管理、package.json处理和依赖解析。

## 核心组件

### 客户端

主要的npm客户端实现。

#### NewClient

```go
func NewClient() (Client, error)
```

使用自动npm检测创建新的npm客户端。

#### NewClientWithPath

```go
func NewClientWithPath(npmPath string) (Client, error)
```

使用特定的npm可执行文件路径创建新的npm客户端。

### 安装器

管理跨不同平台和方法的npm安装。

#### NewInstaller

```go
func NewInstaller() (*Installer, error)
```

使用平台检测创建新的npm安装器。

**示例:**
```go
installer, err := npm.NewInstaller()
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
options := npm.NpmInstallOptions{
    Method:  npm.PackageManager, // 或 npm.OfficialInstaller, npm.Portable
    Version: "latest",
    Force:   false,
}

result, err := installer.Install(ctx, options)
if err != nil {
    log.Fatalf("安装失败: %v", err)
}

fmt.Printf("npm安装成功: %s\n", result.Version)
```

#### 安装方法

安装器支持多种安装方法：

- `PackageManager`: 使用系统包管理器 (apt, yum, brew等)
- `OfficialInstaller`: 下载并运行官方Node.js安装程序
- `Portable`: 下载便携版Node.js/npm
- `Manual`: 手动安装指导

### 检测器

检测npm可用性和版本信息。

#### NewDetector

```go
func NewDetector() *Detector
```

创建新的npm检测器。

#### IsAvailable

```go
func (d *Detector) IsAvailable(ctx context.Context) bool
```

检查系统上是否有npm可用。

#### Detect

```go
func (d *Detector) Detect(ctx context.Context) (*NpmInfo, error)
```

检测npm安装详情。

**示例:**
```go
detector := npm.NewDetector()
ctx := context.Background()

if detector.IsAvailable(ctx) {
    info, err := detector.Detect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("npm位置: %s\n", info.Path)
    fmt.Printf("npm版本: %s\n", info.Version)
} else {
    fmt.Println("未找到npm")
}
```

### 便携版管理器

管理便携版npm安装和版本。

#### NewPortableManager

```go
func NewPortableManager(baseDir string) (*PortableManager, error)
```

使用指定的基础目录创建新的便携版管理器。

#### Install

```go
func (pm *PortableManager) Install(ctx context.Context, version string) (*PortableConfig, error)
```

安装便携版npm版本。

#### List

```go
func (pm *PortableManager) List() ([]*PortableConfig, error)
```

列出所有已安装的便携版本。

#### Uninstall

```go
func (pm *PortableManager) Uninstall(version string) error
```

卸载便携版npm版本。

#### CreateClient

```go
func (pm *PortableManager) CreateClient(version string) (Client, error)
```

为特定便携版本创建客户端。

**示例:**
```go
manager, err := npm.NewPortableManager("/opt/npm-portable")
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

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

// 使用客户端
version, err := client.Version(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("便携版npm版本: %s\n", version)
```

### 包管理器

处理package.json文件操作。

#### NewPackageJSON

```go
func NewPackageJSON(filePath string) *PackageJSON
```

为指定文件创建新的package.json管理器。

#### Load

```go
func (p *PackageJSON) Load() error
```

从磁盘加载package.json。

#### Save

```go
func (p *PackageJSON) Save() error
```

将package.json保存到磁盘。

#### 基本操作

```go
// 获取器
func (p *PackageJSON) GetName() string
func (p *PackageJSON) GetVersion() string
func (p *PackageJSON) GetDescription() string
func (p *PackageJSON) GetAuthor() string
func (p *PackageJSON) GetLicense() string

// 设置器
func (p *PackageJSON) SetName(name string)
func (p *PackageJSON) SetVersion(version string)
func (p *PackageJSON) SetDescription(description string)
func (p *PackageJSON) SetAuthor(author string)
func (p *PackageJSON) SetLicense(license string)
```

#### 依赖管理

```go
// 依赖
func (p *PackageJSON) AddDependency(name, version string)
func (p *PackageJSON) RemoveDependency(name string)
func (p *PackageJSON) GetDependencies() map[string]string

// 开发依赖
func (p *PackageJSON) AddDevDependency(name, version string)
func (p *PackageJSON) RemoveDevDependency(name string)
func (p *PackageJSON) GetDevDependencies() map[string]string

// 对等依赖
func (p *PackageJSON) AddPeerDependency(name, version string)
func (p *PackageJSON) RemovePeerDependency(name string)
func (p *PackageJSON) GetPeerDependencies() map[string]string

// 可选依赖
func (p *PackageJSON) AddOptionalDependency(name, version string)
func (p *PackageJSON) RemoveOptionalDependency(name string)
func (p *PackageJSON) GetOptionalDependencies() map[string]string
```

#### 脚本管理

```go
func (p *PackageJSON) AddScript(name, command string)
func (p *PackageJSON) RemoveScript(name string)
func (p *PackageJSON) GetScripts() map[string]string
func (p *PackageJSON) HasScript(name string) bool
```

#### 仓库和元数据

```go
func (p *PackageJSON) SetRepository(repo *Repository)
func (p *PackageJSON) GetRepository() *Repository
func (p *PackageJSON) SetRepositoryURL(url string)

func (p *PackageJSON) SetBugs(bugs *Bugs)
func (p *PackageJSON) GetBugs() *Bugs
func (p *PackageJSON) SetBugsURL(url string)

func (p *PackageJSON) SetHomepage(homepage string)
func (p *PackageJSON) GetHomepage() string

func (p *PackageJSON) SetKeywords(keywords []string)
func (p *PackageJSON) GetKeywords() []string
func (p *PackageJSON) AddKeyword(keyword string)
```

#### 验证

```go
func (p *PackageJSON) Validate() error
```

验证package.json结构和内容。

**示例:**
```go
pkg := npm.NewPackageJSON("./package.json")

// 加载现有package.json
if err := pkg.Load(); err != nil {
    log.Fatal(err)
}

// 修改包信息
pkg.SetName("my-awesome-package")
pkg.SetVersion("1.0.0")
pkg.SetDescription("一个很棒的包")

// 添加依赖
pkg.AddDependency("lodash", "^4.17.21")
pkg.AddDevDependency("typescript", "^4.5.0")

// 添加脚本
pkg.AddScript("build", "tsc")
pkg.AddScript("test", "jest")

// 设置仓库
repo := &npm.Repository{
    Type: "git",
    URL:  "https://github.com/user/repo.git",
}
pkg.SetRepository(repo)

// 验证并保存
if err := pkg.Validate(); err != nil {
    log.Fatal(err)
}

if err := pkg.Save(); err != nil {
    log.Fatal(err)
}
```

### 依赖管理器

处理依赖解析和冲突检测。

#### NewDependencyManager

```go
func NewDependencyManager() *DependencyManager
```

创建新的依赖管理器。

#### ResolveDependencies

```go
func (dm *DependencyManager) ResolveDependencies(ctx context.Context, deps map[string]string) (*DependencyTree, error)
```

解析依赖并构建依赖树。

#### CheckConflicts

```go
func (dm *DependencyManager) CheckConflicts(tree *DependencyTree) []Conflict
```

检查树中的依赖冲突。

#### DetectCircularDependencies

```go
func (dm *DependencyManager) DetectCircularDependencies(tree *DependencyTree) [][]string
```

检测树中的循环依赖。

**示例:**
```go
manager := npm.NewDependencyManager()
ctx := context.Background()

dependencies := map[string]string{
    "lodash": "^4.17.21",
    "axios":  "^0.24.0",
    "react":  "^17.0.0",
}

// 解析依赖
tree, err := manager.ResolveDependencies(ctx, dependencies)
if err != nil {
    log.Fatal(err)
}

// 检查冲突
conflicts := manager.CheckConflicts(tree)
if len(conflicts) > 0 {
    fmt.Println("发现依赖冲突:")
    for _, conflict := range conflicts {
        fmt.Printf("- %s: %s vs %s\n", conflict.Package, conflict.Version1, conflict.Version2)
    }
}

// 检查循环依赖
circular := manager.DetectCircularDependencies(tree)
if len(circular) > 0 {
    fmt.Println("发现循环依赖:")
    for _, cycle := range circular {
        fmt.Printf("- %s\n", strings.Join(cycle, " -> "))
    }
}
```

## 常量

### 安装方法

```go
const (
    PackageManager    = "package_manager"
    OfficialInstaller = "official_installer"
    Portable          = "portable"
    Manual            = "manual"
)
```

### 错误常量

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

```go
func IsNpmNotFound(err error) bool
func IsPackageNotFound(err error) bool
func IsNetworkError(err error) bool
func IsPlatformError(err error) bool
```

这些函数帮助识别特定错误条件以便更好地处理错误。

## 最佳实践

1. **使用适当的管理器**: 根据您的用例选择正确的组件
2. **优雅地处理错误**: 使用错误检查函数进行特定错误处理
3. **验证输入**: 始终验证包名和版本
4. **使用context**: 传递context以进行超时和取消控制
5. **管理便携版本**: 使用便携版管理器进行隔离的npm环境
6. **备份package.json**: 在修改package.json文件之前始终备份
