package npm

import (
	"context"
	"time"
)

// Client 定义npm客户端的核心接口
type Client interface {
	// 检查npm是否可用
	IsAvailable(ctx context.Context) bool

	// 安装npm
	Install(ctx context.Context) error

	// 获取npm版本
	Version(ctx context.Context) (string, error)

	// 项目初始化
	Init(ctx context.Context, options InitOptions) error

	// 安装包
	InstallPackage(ctx context.Context, pkg string, options InstallOptions) error

	// 卸载包
	UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error

	// 更新包
	UpdatePackage(ctx context.Context, pkg string) error

	// 列出已安装的包
	ListPackages(ctx context.Context, options ListOptions) ([]Package, error)

	// 运行脚本
	RunScript(ctx context.Context, script string, args ...string) error

	// 发布包
	Publish(ctx context.Context, options PublishOptions) error

	// 获取包信息
	GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error)

	// 搜索包
	Search(ctx context.Context, query string) ([]SearchResult, error)
}

// InitOptions 项目初始化选项
type InitOptions struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	License     string `json:"license,omitempty"`
	Private     bool   `json:"private,omitempty"`
	WorkingDir  string `json:"-"` // 工作目录，不序列化到package.json
	Force       bool   `json:"-"` // 强制覆盖，不序列化到package.json
}

// InstallOptions 安装选项
type InstallOptions struct {
	SaveDev       bool   `json:"save_dev,omitempty"`       // --save-dev
	SaveOptional  bool   `json:"save_optional,omitempty"`  // --save-optional
	SaveExact     bool   `json:"save_exact,omitempty"`     // --save-exact
	Global        bool   `json:"global,omitempty"`         // --global
	Production    bool   `json:"production,omitempty"`     // --production
	WorkingDir    string `json:"working_dir,omitempty"`    // 工作目录
	Registry      string `json:"registry,omitempty"`       // 自定义registry
	Force         bool   `json:"force,omitempty"`          // --force
	IgnoreScripts bool   `json:"ignore_scripts,omitempty"` // --ignore-scripts
}

// UninstallOptions 卸载选项
type UninstallOptions struct {
	SaveDev    bool   `json:"save_dev,omitempty"`    // --save-dev
	Global     bool   `json:"global,omitempty"`      // --global
	WorkingDir string `json:"working_dir,omitempty"` // 工作目录
}

// ListOptions 列表选项
type ListOptions struct {
	Global     bool   `json:"global,omitempty"`      // --global
	Depth      int    `json:"depth,omitempty"`       // --depth
	Production bool   `json:"production,omitempty"`  // --production
	WorkingDir string `json:"working_dir,omitempty"` // 工作目录
	JSON       bool   `json:"json,omitempty"`        // --json
}

// PublishOptions 发布选项
type PublishOptions struct {
	Tag        string `json:"tag,omitempty"`         // --tag
	Access     string `json:"access,omitempty"`      // --access (public/restricted)
	Registry   string `json:"registry,omitempty"`    // 自定义registry
	WorkingDir string `json:"working_dir,omitempty"` // 工作目录
	DryRun     bool   `json:"dry_run,omitempty"`     // --dry-run
}

// Package 表示一个npm包
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

// Repository 仓库信息
type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Bugs bug报告信息
type Bugs struct {
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// PackageInfo 包详细信息
type PackageInfo struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Keywords     []string               `json:"keywords"`
	Homepage     string                 `json:"homepage"`
	Repository   *Repository            `json:"repository"`
	Author       *Person                `json:"author"`
	License      string                 `json:"license"`
	Dependencies map[string]string      `json:"dependencies"`
	DevDeps      map[string]string      `json:"devDependencies"`
	Versions     map[string]interface{} `json:"versions"`
	Time         map[string]time.Time   `json:"time"`
	DistTags     map[string]string      `json:"dist-tags"`
}

// Person 人员信息
type Person struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Package     SearchPackage `json:"package"`
	Score       SearchScore   `json:"score"`
	SearchScore float64       `json:"searchScore"`
}

// SearchPackage 搜索包信息
type SearchPackage struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Keywords    []string          `json:"keywords"`
	Date        time.Time         `json:"date"`
	Links       map[string]string `json:"links"`
	Author      *Person           `json:"author"`
	Publisher   *Person           `json:"publisher"`
	Maintainers []*Person         `json:"maintainers"`
}

// SearchScore 搜索评分
type SearchScore struct {
	Final  float64     `json:"final"`
	Detail ScoreDetail `json:"detail"`
}

// ScoreDetail 评分详情
type ScoreDetail struct {
	Quality     float64 `json:"quality"`
	Popularity  float64 `json:"popularity"`
	Maintenance float64 `json:"maintenance"`
}

// CommandResult 命令执行结果
type CommandResult struct {
	Success  bool          `json:"success"`
	ExitCode int           `json:"exit_code"`
	Stdout   string        `json:"stdout"`
	Stderr   string        `json:"stderr"`
	Duration time.Duration `json:"duration"`
	Error    error         `json:"error,omitempty"`
}
