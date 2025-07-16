package npm

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// DependencyType 依赖类型
type DependencyType string

const (
	Production DependencyType = "dependencies"
	Development DependencyType = "devDependencies"
	Optional DependencyType = "optionalDependencies"
	Peer DependencyType = "peerDependencies"
)

// DependencyManager 依赖管理器
type DependencyManager struct {
	client     Client
	packageJSON *PackageJSON
	workingDir string
}

// DependencyInfo 依赖信息
type DependencyInfo struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Type         DependencyType `json:"type"`
	Installed    bool           `json:"installed"`
	Latest       string         `json:"latest,omitempty"`
	Description  string         `json:"description,omitempty"`
}

// DependencyOperation 依赖操作结果
type DependencyOperation struct {
	Success     bool              `json:"success"`
	Operation   string            `json:"operation"`
	Package     string            `json:"package"`
	Version     string            `json:"version"`
	Type        DependencyType    `json:"type"`
	Error       error             `json:"error,omitempty"`
	Changes     []string          `json:"changes,omitempty"`
}

// NewDependencyManager 创建依赖管理器
func NewDependencyManager(client Client, workingDir string) (*DependencyManager, error) {
	packageJSONPath := filepath.Join(workingDir, "package.json")
	packageJSON := NewPackageJSON(packageJSONPath)
	
	return &DependencyManager{
		client:      client,
		packageJSON: packageJSON,
		workingDir:  workingDir,
	}, nil
}

// LoadPackageJSON 加载package.json
func (dm *DependencyManager) LoadPackageJSON() error {
	return dm.packageJSON.Load()
}

// SavePackageJSON 保存package.json
func (dm *DependencyManager) SavePackageJSON() error {
	return dm.packageJSON.Save()
}

// Add 添加依赖
func (dm *DependencyManager) Add(ctx context.Context, packageName, version string, depType DependencyType) (*DependencyOperation, error) {
	operation := &DependencyOperation{
		Operation: "add",
		Package:   packageName,
		Version:   version,
		Type:      depType,
	}

	// 验证包名
	if packageName == "" {
		operation.Error = NewValidationError("package", packageName, "package name cannot be empty")
		return operation, operation.Error
	}

	// 如果没有指定版本，获取最新版本
	if version == "" {
		packageInfo, err := dm.client.GetPackageInfo(ctx, packageName)
		if err != nil {
			operation.Error = fmt.Errorf("failed to get package info: %w", err)
			return operation, operation.Error
		}
		version = "^" + packageInfo.Version
		operation.Version = version
	}

	// 安装包
	installOptions := InstallOptions{
		WorkingDir: dm.workingDir,
		SaveDev:    depType == Development,
		SaveOptional: depType == Optional,
	}

	if err := dm.client.InstallPackage(ctx, packageName, installOptions); err != nil {
		operation.Error = fmt.Errorf("failed to install package: %w", err)
		return operation, operation.Error
	}

	// 更新package.json
	if err := dm.LoadPackageJSON(); err == nil {
		switch depType {
		case Production:
			dm.packageJSON.AddDependency(packageName, version)
		case Development:
			dm.packageJSON.AddDevDependency(packageName, version)
		case Optional:
			dm.packageJSON.AddOptionalDependency(packageName, version)
		case Peer:
			dm.packageJSON.AddPeerDependency(packageName, version)
		}
		
		if err := dm.SavePackageJSON(); err != nil {
			operation.Changes = append(operation.Changes, fmt.Sprintf("Warning: failed to update package.json: %v", err))
		} else {
			operation.Changes = append(operation.Changes, "Updated package.json")
		}
	}

	operation.Success = true
	operation.Changes = append(operation.Changes, fmt.Sprintf("Installed %s@%s", packageName, version))
	
	return operation, nil
}

// Remove 移除依赖
func (dm *DependencyManager) Remove(ctx context.Context, packageName string) (*DependencyOperation, error) {
	operation := &DependencyOperation{
		Operation: "remove",
		Package:   packageName,
	}

	// 验证包名
	if packageName == "" {
		operation.Error = NewValidationError("package", packageName, "package name cannot be empty")
		return operation, operation.Error
	}

	// 检查依赖类型
	if err := dm.LoadPackageJSON(); err == nil {
		if dm.packageJSON.HasDependency(packageName) {
			operation.Type = Production
		} else if dm.packageJSON.HasDevDependency(packageName) {
			operation.Type = Development
		}
	}

	// 卸载包
	uninstallOptions := UninstallOptions{
		WorkingDir: dm.workingDir,
		SaveDev:    operation.Type == Development,
	}

	if err := dm.client.UninstallPackage(ctx, packageName, uninstallOptions); err != nil {
		operation.Error = fmt.Errorf("failed to uninstall package: %w", err)
		return operation, operation.Error
	}

	// 更新package.json
	if err := dm.LoadPackageJSON(); err == nil {
		dm.packageJSON.RemoveDependency(packageName)
		dm.packageJSON.RemoveDevDependency(packageName)
		dm.packageJSON.RemoveOptionalDependency(packageName)
		dm.packageJSON.RemovePeerDependency(packageName)
		
		if err := dm.SavePackageJSON(); err != nil {
			operation.Changes = append(operation.Changes, fmt.Sprintf("Warning: failed to update package.json: %v", err))
		} else {
			operation.Changes = append(operation.Changes, "Updated package.json")
		}
	}

	operation.Success = true
	operation.Changes = append(operation.Changes, fmt.Sprintf("Removed %s", packageName))
	
	return operation, nil
}

// Update 更新依赖
func (dm *DependencyManager) Update(ctx context.Context, packageName string) (*DependencyOperation, error) {
	operation := &DependencyOperation{
		Operation: "update",
		Package:   packageName,
	}

	// 验证包名
	if packageName == "" {
		operation.Error = NewValidationError("package", packageName, "package name cannot be empty")
		return operation, operation.Error
	}

	// 获取当前版本
	if err := dm.LoadPackageJSON(); err == nil {
		deps := dm.packageJSON.GetDependencies()
		devDeps := dm.packageJSON.GetDevDependencies()
		
		if version, exists := deps[packageName]; exists {
			operation.Version = version
			operation.Type = Production
		} else if version, exists := devDeps[packageName]; exists {
			operation.Version = version
			operation.Type = Development
		}
	}

	// 更新包
	if err := dm.client.UpdatePackage(ctx, packageName); err != nil {
		operation.Error = fmt.Errorf("failed to update package: %w", err)
		return operation, operation.Error
	}

	operation.Success = true
	operation.Changes = append(operation.Changes, fmt.Sprintf("Updated %s", packageName))
	
	return operation, nil
}

// List 列出所有依赖
func (dm *DependencyManager) List(ctx context.Context) ([]*DependencyInfo, error) {
	var dependencies []*DependencyInfo

	// 从package.json读取依赖
	if err := dm.LoadPackageJSON(); err != nil {
		return nil, fmt.Errorf("failed to load package.json: %w", err)
	}

	// 获取已安装的包列表
	installedPackages, err := dm.client.ListPackages(ctx, ListOptions{
		WorkingDir: dm.workingDir,
		Depth:      0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}

	installedMap := make(map[string]bool)
	for _, pkg := range installedPackages {
		installedMap[pkg.Name] = true
	}

	// 处理生产依赖
	for name, version := range dm.packageJSON.GetDependencies() {
		dep := &DependencyInfo{
			Name:      name,
			Version:   version,
			Type:      Production,
			Installed: installedMap[name],
		}
		dependencies = append(dependencies, dep)
	}

	// 处理开发依赖
	for name, version := range dm.packageJSON.GetDevDependencies() {
		dep := &DependencyInfo{
			Name:      name,
			Version:   version,
			Type:      Development,
			Installed: installedMap[name],
		}
		dependencies = append(dependencies, dep)
	}

	// 处理可选依赖
	for name, version := range dm.packageJSON.GetOptionalDependencies() {
		dep := &DependencyInfo{
			Name:      name,
			Version:   version,
			Type:      Optional,
			Installed: installedMap[name],
		}
		dependencies = append(dependencies, dep)
	}

	// 处理同级依赖
	for name, version := range dm.packageJSON.GetPeerDependencies() {
		dep := &DependencyInfo{
			Name:      name,
			Version:   version,
			Type:      Peer,
			Installed: installedMap[name],
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

// CheckOutdated 检查过期的依赖
func (dm *DependencyManager) CheckOutdated(ctx context.Context) ([]*DependencyInfo, error) {
	dependencies, err := dm.List(ctx)
	if err != nil {
		return nil, err
	}

	var outdated []*DependencyInfo

	for _, dep := range dependencies {
		if !dep.Installed {
			continue
		}

		// 获取最新版本信息
		packageInfo, err := dm.client.GetPackageInfo(ctx, dep.Name)
		if err != nil {
			continue // 跳过无法获取信息的包
		}

		dep.Latest = packageInfo.Version
		dep.Description = packageInfo.Description

		// 简单的版本比较（实际应该使用更复杂的语义化版本比较）
		if dep.Latest != strings.TrimPrefix(dep.Version, "^") && 
		   dep.Latest != strings.TrimPrefix(dep.Version, "~") {
			outdated = append(outdated, dep)
		}
	}

	return outdated, nil
}

// Install 安装所有依赖
func (dm *DependencyManager) Install(ctx context.Context) error {
	installOptions := InstallOptions{
		WorkingDir: dm.workingDir,
	}

	// 不指定包名时，npm install会安装package.json中的所有依赖
	return dm.client.InstallPackage(ctx, "", installOptions)
}

// Clean 清理node_modules并重新安装
func (dm *DependencyManager) Clean(ctx context.Context) error {
	// 这里可以添加删除node_modules的逻辑
	// 然后重新安装所有依赖
	return dm.Install(ctx)
}

// GetDependencyTree 获取依赖树
func (dm *DependencyManager) GetDependencyTree(ctx context.Context) ([]Package, error) {
	return dm.client.ListPackages(ctx, ListOptions{
		WorkingDir: dm.workingDir,
		Depth:      -1, // 获取完整依赖树
		JSON:       true,
	})
}

// Audit 安全审计
func (dm *DependencyManager) Audit(ctx context.Context) error {
	// 这里可以实现npm audit功能
	// 目前返回未实现错误
	return fmt.Errorf("audit functionality not implemented yet")
}

// Fix 修复安全漏洞
func (dm *DependencyManager) Fix(ctx context.Context) error {
	// 这里可以实现npm audit fix功能
	// 目前返回未实现错误
	return fmt.Errorf("fix functionality not implemented yet")
}
