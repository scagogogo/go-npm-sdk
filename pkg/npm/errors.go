package npm

import (
	"errors"
	"fmt"
)

// 预定义错误
var (
	// ErrNpmNotFound npm未找到
	ErrNpmNotFound = errors.New("npm not found")

	// ErrNpmNotInstalled npm未安装
	ErrNpmNotInstalled = errors.New("npm is not installed")

	// ErrInvalidPackageName 无效的包名
	ErrInvalidPackageName = errors.New("invalid package name")

	// ErrPackageNotFound 包未找到
	ErrPackageNotFound = errors.New("package not found")

	// ErrPackageAlreadyExists 包已存在
	ErrPackageAlreadyExists = errors.New("package already exists")

	// ErrInvalidVersion 无效版本
	ErrInvalidVersion = errors.New("invalid version")

	// ErrNetworkError 网络错误
	ErrNetworkError = errors.New("network error")

	// ErrPermissionDenied 权限被拒绝
	ErrPermissionDenied = errors.New("permission denied")

	// ErrInvalidWorkingDirectory 无效的工作目录
	ErrInvalidWorkingDirectory = errors.New("invalid working directory")

	// ErrCommandTimeout 命令超时
	ErrCommandTimeout = errors.New("command timeout")

	// ErrInvalidPackageJSON 无效的package.json
	ErrInvalidPackageJSON = errors.New("invalid package.json")

	// ErrRegistryError registry错误
	ErrRegistryError = errors.New("registry error")

	// ErrAuthenticationFailed 认证失败
	ErrAuthenticationFailed = errors.New("authentication failed")

	// ErrUnsupportedPlatform 不支持的平台
	ErrUnsupportedPlatform = errors.New("unsupported platform")
)

// NpmError npm操作错误
type NpmError struct {
	Op       string // 操作名称
	Package  string // 包名（如果适用）
	ExitCode int    // 退出码
	Stdout   string // 标准输出
	Stderr   string // 标准错误
	Err      error  // 原始错误
}

func (e *NpmError) Error() string {
	if e.Package != "" {
		return fmt.Sprintf("npm %s failed for package '%s': %v", e.Op, e.Package, e.Err)
	}
	return fmt.Sprintf("npm %s failed: %v", e.Op, e.Err)
}

func (e *NpmError) Unwrap() error {
	return e.Err
}

// NewNpmError 创建npm错误
func NewNpmError(op, pkg string, exitCode int, stdout, stderr string, err error) *NpmError {
	return &NpmError{
		Op:       op,
		Package:  pkg,
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
		Err:      err,
	}
}

// InstallError 安装错误
type InstallError struct {
	Package string
	Reason  string
	Err     error
}

func (e *InstallError) Error() string {
	return fmt.Sprintf("failed to install package '%s': %s", e.Package, e.Reason)
}

func (e *InstallError) Unwrap() error {
	return e.Err
}

// NewInstallError 创建安装错误
func NewInstallError(pkg, reason string, err error) *InstallError {
	return &InstallError{
		Package: pkg,
		Reason:  reason,
		Err:     err,
	}
}

// UninstallError 卸载错误
type UninstallError struct {
	Package string
	Reason  string
	Err     error
}

func (e *UninstallError) Error() string {
	return fmt.Sprintf("failed to uninstall package '%s': %s", e.Package, e.Reason)
}

func (e *UninstallError) Unwrap() error {
	return e.Err
}

// NewUninstallError 创建卸载错误
func NewUninstallError(pkg, reason string, err error) *UninstallError {
	return &UninstallError{
		Package: pkg,
		Reason:  reason,
		Err:     err,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Value   string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s' with value '%s': %s", e.Field, e.Value, e.Reason)
}

// NewValidationError 创建验证错误
func NewValidationError(field, value, reason string) *ValidationError {
	return &ValidationError{
		Field:  field,
		Value:  value,
		Reason: reason,
	}
}

// PlatformError 平台相关错误
type PlatformError struct {
	Platform string
	Reason   string
	Err      error
}

func (e *PlatformError) Error() string {
	return fmt.Sprintf("platform error on %s: %s", e.Platform, e.Reason)
}

func (e *PlatformError) Unwrap() error {
	return e.Err
}

// NewPlatformError 创建平台错误
func NewPlatformError(platform, reason string, err error) *PlatformError {
	return &PlatformError{
		Platform: platform,
		Reason:   reason,
		Err:      err,
	}
}

// DownloadError 下载错误
type DownloadError struct {
	URL    string
	Reason string
	Err    error
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("failed to download from %s: %s", e.URL, e.Reason)
}

func (e *DownloadError) Unwrap() error {
	return e.Err
}

// NewDownloadError 创建下载错误
func NewDownloadError(url, reason string, err error) *DownloadError {
	return &DownloadError{
		URL:    url,
		Reason: reason,
		Err:    err,
	}
}

// IsNpmNotFound 检查是否为npm未找到错误
func IsNpmNotFound(err error) bool {
	return errors.Is(err, ErrNpmNotFound)
}

// IsPackageNotFound 检查是否为包未找到错误
func IsPackageNotFound(err error) bool {
	return errors.Is(err, ErrPackageNotFound)
}

// IsNetworkError 检查是否为网络错误
func IsNetworkError(err error) bool {
	return errors.Is(err, ErrNetworkError)
}

// IsPermissionDenied 检查是否为权限错误
func IsPermissionDenied(err error) bool {
	return errors.Is(err, ErrPermissionDenied)
}

// IsTimeout 检查是否为超时错误
func IsTimeout(err error) bool {
	return errors.Is(err, ErrCommandTimeout)
}

// IsUnsupportedPlatform 检查是否为不支持的平台错误
func IsUnsupportedPlatform(err error) bool {
	return errors.Is(err, ErrUnsupportedPlatform)
}
