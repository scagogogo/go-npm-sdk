package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/scagogogo/go-npm-sdk/pkg/utils"
)

// client npm客户端实现
type client struct {
	npmPath   string
	executor  *utils.Executor
	detector  *Detector
	installer *Installer
}

// NewClient 创建新的npm客户端
func NewClient() (Client, error) {
	detector := NewDetector()
	installer, err := NewInstaller()
	if err != nil {
		return nil, fmt.Errorf("failed to create installer: %w", err)
	}

	return &client{
		npmPath:   "npm",
		executor:  utils.NewExecutor(),
		detector:  detector,
		installer: installer,
	}, nil
}

// NewClientWithPath 使用指定路径创建npm客户端
func NewClientWithPath(npmPath string) (Client, error) {
	detector := NewDetector()
	installer, err := NewInstaller()
	if err != nil {
		return nil, fmt.Errorf("failed to create installer: %w", err)
	}

	return &client{
		npmPath:   npmPath,
		executor:  utils.NewExecutor(),
		detector:  detector,
		installer: installer,
	}, nil
}

// IsAvailable 检查npm是否可用
func (c *client) IsAvailable(ctx context.Context) bool {
	result, err := c.executor.ExecuteSimple(ctx, c.npmPath, "--version")
	return err == nil && result.Success
}

// Install 安装npm
func (c *client) Install(ctx context.Context) error {
	options := NpmInstallOptions{
		Method: PackageManager, // 默认使用包管理器
		Force:  false,
	}

	result, err := c.installer.Install(ctx, options)
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	// 更新npm路径
	if result.Path != "" {
		c.npmPath = result.Path
	}

	return nil
}

// Version 获取npm版本
func (c *client) Version(ctx context.Context) (string, error) {
	result, err := c.executor.ExecuteSimple(ctx, c.npmPath, "--version")
	if err != nil {
		return "", NewNpmError("version", "", result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return "", NewNpmError("version", "", result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("failed to get version"))
	}

	return strings.TrimSpace(result.Stdout), nil
}

// Init 项目初始化
func (c *client) Init(ctx context.Context, options InitOptions) error {
	args := []string{"init"}

	// 构建参数
	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}
	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}
	if options.Description != "" {
		args = append(args, "--description", options.Description)
	}
	if options.Author != "" {
		args = append(args, "--author", options.Author)
	}
	if options.License != "" {
		args = append(args, "--license", options.License)
	}
	if options.Private {
		args = append(args, "--private")
	}
	if options.Force {
		args = append(args, "--yes")
	}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		WorkingDir:    options.WorkingDir,
		CaptureOutput: true,
		Timeout:       2 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return NewNpmError("init", "", result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return NewNpmError("init", "", result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm init failed"))
	}

	return nil
}

// InstallPackage 安装包
func (c *client) InstallPackage(ctx context.Context, pkg string, options InstallOptions) error {
	if pkg == "" {
		return NewValidationError("package", pkg, "package name cannot be empty")
	}

	args := []string{"install", pkg}

	// 构建参数
	if options.SaveDev {
		args = append(args, "--save-dev")
	}
	if options.SaveOptional {
		args = append(args, "--save-optional")
	}
	if options.SaveExact {
		args = append(args, "--save-exact")
	}
	if options.Global {
		args = append(args, "--global")
	}
	if options.Production {
		args = append(args, "--production")
	}
	if options.Registry != "" {
		args = append(args, "--registry", options.Registry)
	}
	if options.Force {
		args = append(args, "--force")
	}
	if options.IgnoreScripts {
		args = append(args, "--ignore-scripts")
	}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		WorkingDir:    options.WorkingDir,
		CaptureOutput: true,
		Timeout:       10 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return NewInstallError(pkg, "execution failed", NewNpmError("install", pkg, result.ExitCode, result.Stdout, result.Stderr, err))
	}

	if !result.Success {
		return NewInstallError(pkg, "npm install failed", NewNpmError("install", pkg, result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("install failed")))
	}

	return nil
}

// UninstallPackage 卸载包
func (c *client) UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error {
	if pkg == "" {
		return NewValidationError("package", pkg, "package name cannot be empty")
	}

	args := []string{"uninstall", pkg}

	// 构建参数
	if options.SaveDev {
		args = append(args, "--save-dev")
	}
	if options.Global {
		args = append(args, "--global")
	}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		WorkingDir:    options.WorkingDir,
		CaptureOutput: true,
		Timeout:       5 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return NewUninstallError(pkg, "execution failed", NewNpmError("uninstall", pkg, result.ExitCode, result.Stdout, result.Stderr, err))
	}

	if !result.Success {
		return NewUninstallError(pkg, "npm uninstall failed", NewNpmError("uninstall", pkg, result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("uninstall failed")))
	}

	return nil
}

// UpdatePackage 更新包
func (c *client) UpdatePackage(ctx context.Context, pkg string) error {
	if pkg == "" {
		return NewValidationError("package", pkg, "package name cannot be empty")
	}

	args := []string{"update", pkg}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		CaptureOutput: true,
		Timeout:       10 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return NewNpmError("update", pkg, result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return NewNpmError("update", pkg, result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm update failed"))
	}

	return nil
}

// ListPackages 列出已安装的包
func (c *client) ListPackages(ctx context.Context, options ListOptions) ([]Package, error) {
	args := []string{"list"}

	// 构建参数
	if options.Global {
		args = append(args, "--global")
	}
	if options.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", options.Depth))
	}
	if options.Production {
		args = append(args, "--production")
	}
	if options.JSON {
		args = append(args, "--json")
	}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		WorkingDir:    options.WorkingDir,
		CaptureOutput: true,
		Timeout:       2 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return nil, NewNpmError("list", "", result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return nil, NewNpmError("list", "", result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm list failed"))
	}

	// 解析JSON输出
	if options.JSON {
		return c.parseListJSON(result.Stdout)
	}

	// 解析文本输出
	return c.parseListText(result.Stdout)
}

// RunScript 运行脚本
func (c *client) RunScript(ctx context.Context, script string, args ...string) error {
	if script == "" {
		return NewValidationError("script", script, "script name cannot be empty")
	}

	cmdArgs := []string{"run", script}
	if len(args) > 0 {
		cmdArgs = append(cmdArgs, "--")
		cmdArgs = append(cmdArgs, args...)
	}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          cmdArgs,
		CaptureOutput: true,
		StreamOutput:  true,
		Timeout:       30 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return NewNpmError("run", script, result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return NewNpmError("run", script, result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm run failed"))
	}

	return nil
}

// parseListJSON 解析JSON格式的list输出
func (c *client) parseListJSON(output string) ([]Package, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	var packages []Package

	// 解析dependencies
	if deps, ok := data["dependencies"].(map[string]interface{}); ok {
		for name, info := range deps {
			if pkgInfo, ok := info.(map[string]interface{}); ok {
				pkg := Package{
					Name: name,
				}
				if version, ok := pkgInfo["version"].(string); ok {
					pkg.Version = version
				}
				packages = append(packages, pkg)
			}
		}
	}

	return packages, nil
}

// parseListText 解析文本格式的list输出
func (c *client) parseListText(output string) ([]Package, error) {
	lines := strings.Split(output, "\n")
	var packages []Package

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "npm") {
			continue
		}

		// 简单解析，实际实现可能需要更复杂的逻辑
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := strings.TrimPrefix(parts[0], "├──")
			name = strings.TrimPrefix(name, "└──")
			name = strings.TrimSpace(name)

			if strings.Contains(name, "@") {
				nameParts := strings.Split(name, "@")
				if len(nameParts) >= 2 {
					packages = append(packages, Package{
						Name:    nameParts[0],
						Version: nameParts[len(nameParts)-1],
					})
				}
			}
		}
	}

	return packages, nil
}

// Publish 发布包
func (c *client) Publish(ctx context.Context, options PublishOptions) error {
	args := []string{"publish"}

	// 构建参数
	if options.Tag != "" {
		args = append(args, "--tag", options.Tag)
	}
	if options.Access != "" {
		args = append(args, "--access", options.Access)
	}
	if options.Registry != "" {
		args = append(args, "--registry", options.Registry)
	}
	if options.DryRun {
		args = append(args, "--dry-run")
	}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		WorkingDir:    options.WorkingDir,
		CaptureOutput: true,
		Timeout:       10 * time.Minute,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return NewNpmError("publish", "", result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return NewNpmError("publish", "", result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm publish failed"))
	}

	return nil
}

// GetPackageInfo 获取包信息
func (c *client) GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error) {
	if pkg == "" {
		return nil, NewValidationError("package", pkg, "package name cannot be empty")
	}

	args := []string{"view", pkg, "--json"}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		CaptureOutput: true,
		Timeout:       30 * time.Second,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return nil, NewNpmError("view", pkg, result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return nil, NewNpmError("view", pkg, result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm view failed"))
	}

	var info PackageInfo
	if err := json.Unmarshal([]byte(result.Stdout), &info); err != nil {
		return nil, fmt.Errorf("failed to parse package info: %w", err)
	}

	return &info, nil
}

// Search 搜索包
func (c *client) Search(ctx context.Context, query string) ([]SearchResult, error) {
	if query == "" {
		return nil, NewValidationError("query", query, "search query cannot be empty")
	}

	args := []string{"search", query, "--json"}

	executeOptions := utils.ExecuteOptions{
		Command:       c.npmPath,
		Args:          args,
		CaptureOutput: true,
		Timeout:       30 * time.Second,
	}

	result, err := c.executor.Execute(ctx, executeOptions)
	if err != nil {
		return nil, NewNpmError("search", query, result.ExitCode, result.Stdout, result.Stderr, err)
	}

	if !result.Success {
		return nil, NewNpmError("search", query, result.ExitCode, result.Stdout, result.Stderr, fmt.Errorf("npm search failed"))
	}

	var results []SearchResult
	if err := json.Unmarshal([]byte(result.Stdout), &results); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return results, nil
}
