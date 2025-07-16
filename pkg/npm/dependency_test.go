package npm

import (
	"context"
	"path/filepath"
	"testing"
)

// MockClient 用于测试的模拟客户端
type MockClient struct {
	packages  map[string]*PackageInfo
	installed map[string]bool
}

func NewMockClient() *MockClient {
	return &MockClient{
		packages:  make(map[string]*PackageInfo),
		installed: make(map[string]bool),
	}
}

func (m *MockClient) IsAvailable(ctx context.Context) bool {
	return true
}

func (m *MockClient) Install(ctx context.Context) error {
	return nil
}

func (m *MockClient) Version(ctx context.Context) (string, error) {
	return "8.0.0", nil
}

func (m *MockClient) Init(ctx context.Context, options InitOptions) error {
	return nil
}

func (m *MockClient) InstallPackage(ctx context.Context, pkg string, options InstallOptions) error {
	if pkg == "" {
		return nil // 安装所有依赖
	}
	m.installed[pkg] = true
	return nil
}

func (m *MockClient) UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error {
	delete(m.installed, pkg)
	return nil
}

func (m *MockClient) UpdatePackage(ctx context.Context, pkg string) error {
	return nil
}

func (m *MockClient) ListPackages(ctx context.Context, options ListOptions) ([]Package, error) {
	var packages []Package
	for name := range m.installed {
		packages = append(packages, Package{
			Name:    name,
			Version: "1.0.0",
		})
	}
	return packages, nil
}

func (m *MockClient) RunScript(ctx context.Context, script string, args ...string) error {
	return nil
}

func (m *MockClient) Publish(ctx context.Context, options PublishOptions) error {
	return nil
}

func (m *MockClient) GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error) {
	if info, exists := m.packages[pkg]; exists {
		return info, nil
	}

	// 返回默认信息
	return &PackageInfo{
		Name:        pkg,
		Version:     "1.0.0",
		Description: "Test package",
		License:     "MIT",
	}, nil
}

func (m *MockClient) Search(ctx context.Context, query string) ([]SearchResult, error) {
	return []SearchResult{}, nil
}

func (m *MockClient) AddPackage(name, version, description string) {
	m.packages[name] = &PackageInfo{
		Name:        name,
		Version:     version,
		Description: description,
		License:     "MIT",
	}
}

func TestNewDependencyManager(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	if dm == nil {
		t.Fatal("NewDependencyManager() returned nil")
	}

	if dm.workingDir != tempDir {
		t.Errorf("Expected working dir %s, got %s", tempDir, dm.workingDir)
	}
}

func TestDependencyManagerAdd(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	// 创建package.json
	packageJSONPath := filepath.Join(tempDir, "package.json")
	pkg := NewPackageJSON(packageJSONPath)
	pkg.SetName("test-project")
	pkg.SetVersion("1.0.0")
	if err := pkg.Save(); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	ctx := context.Background()

	// 添加生产依赖
	operation, err := dm.Add(ctx, "lodash", "^4.17.21", Production)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if !operation.Success {
		t.Error("Expected operation to succeed")
	}

	if operation.Package != "lodash" {
		t.Errorf("Expected package 'lodash', got '%s'", operation.Package)
	}

	if operation.Type != Production {
		t.Errorf("Expected type Production, got %s", operation.Type)
	}

	// 验证包已安装
	if !client.installed["lodash"] {
		t.Error("Expected lodash to be installed")
	}

	// 添加开发依赖
	operation, err = dm.Add(ctx, "jest", "^27.0.0", Development)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if !operation.Success {
		t.Error("Expected operation to succeed")
	}

	if operation.Type != Development {
		t.Errorf("Expected type Development, got %s", operation.Type)
	}
}

func TestDependencyManagerRemove(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	// 创建package.json
	packageJSONPath := filepath.Join(tempDir, "package.json")
	pkg := NewPackageJSON(packageJSONPath)
	pkg.SetName("test-project")
	pkg.SetVersion("1.0.0")
	pkg.AddDependency("lodash", "^4.17.21")
	if err := pkg.Save(); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// 模拟已安装的包
	client.installed["lodash"] = true

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	ctx := context.Background()

	// 移除依赖
	operation, err := dm.Remove(ctx, "lodash")
	if err != nil {
		t.Fatalf("Remove() failed: %v", err)
	}

	if !operation.Success {
		t.Error("Expected operation to succeed")
	}

	if operation.Package != "lodash" {
		t.Errorf("Expected package 'lodash', got '%s'", operation.Package)
	}

	// 验证包已卸载
	if client.installed["lodash"] {
		t.Error("Expected lodash to be uninstalled")
	}
}

func TestDependencyManagerUpdate(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	// 创建package.json
	packageJSONPath := filepath.Join(tempDir, "package.json")
	pkg := NewPackageJSON(packageJSONPath)
	pkg.SetName("test-project")
	pkg.SetVersion("1.0.0")
	pkg.AddDependency("lodash", "^4.17.20")
	if err := pkg.Save(); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	ctx := context.Background()

	// 更新依赖
	operation, err := dm.Update(ctx, "lodash")
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if !operation.Success {
		t.Error("Expected operation to succeed")
	}

	if operation.Package != "lodash" {
		t.Errorf("Expected package 'lodash', got '%s'", operation.Package)
	}

	if operation.Type != Production {
		t.Errorf("Expected type Production, got %s", operation.Type)
	}
}

func TestDependencyManagerList(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	// 创建package.json
	packageJSONPath := filepath.Join(tempDir, "package.json")
	pkg := NewPackageJSON(packageJSONPath)
	pkg.SetName("test-project")
	pkg.SetVersion("1.0.0")
	pkg.AddDependency("lodash", "^4.17.21")
	pkg.AddDevDependency("jest", "^27.0.0")
	if err := pkg.Save(); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// 模拟已安装的包
	client.installed["lodash"] = true
	client.installed["jest"] = true

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	ctx := context.Background()

	// 列出依赖
	dependencies, err := dm.List(ctx)
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(dependencies))
	}

	// 检查依赖信息
	var lodashFound, jestFound bool
	for _, dep := range dependencies {
		switch dep.Name {
		case "lodash":
			lodashFound = true
			if dep.Type != Production {
				t.Errorf("Expected lodash type Production, got %s", dep.Type)
			}
			if !dep.Installed {
				t.Error("Expected lodash to be installed")
			}
		case "jest":
			jestFound = true
			if dep.Type != Development {
				t.Errorf("Expected jest type Development, got %s", dep.Type)
			}
			if !dep.Installed {
				t.Error("Expected jest to be installed")
			}
		}
	}

	if !lodashFound {
		t.Error("Expected to find lodash dependency")
	}

	if !jestFound {
		t.Error("Expected to find jest dependency")
	}
}

func TestDependencyManagerCheckOutdated(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	// 添加包信息
	client.AddPackage("lodash", "4.17.21", "Lodash utility library")
	client.AddPackage("jest", "28.0.0", "JavaScript testing framework")

	// 创建package.json（使用旧版本）
	packageJSONPath := filepath.Join(tempDir, "package.json")
	pkg := NewPackageJSON(packageJSONPath)
	pkg.SetName("test-project")
	pkg.SetVersion("1.0.0")
	pkg.AddDependency("lodash", "^4.17.20") // 旧版本
	pkg.AddDevDependency("jest", "^27.0.0") // 旧版本
	if err := pkg.Save(); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// 模拟已安装的包
	client.installed["lodash"] = true
	client.installed["jest"] = true

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	ctx := context.Background()

	// 检查过期依赖
	outdated, err := dm.CheckOutdated(ctx)
	if err != nil {
		t.Fatalf("CheckOutdated() failed: %v", err)
	}

	if len(outdated) != 2 {
		t.Errorf("Expected 2 outdated dependencies, got %d", len(outdated))
	}

	// 检查过期信息
	for _, dep := range outdated {
		switch dep.Name {
		case "lodash":
			if dep.Latest != "4.17.21" {
				t.Errorf("Expected lodash latest version '4.17.21', got '%s'", dep.Latest)
			}
		case "jest":
			if dep.Latest != "28.0.0" {
				t.Errorf("Expected jest latest version '28.0.0', got '%s'", dep.Latest)
			}
		}
	}
}

func TestDependencyManagerValidation(t *testing.T) {
	client := NewMockClient()
	tempDir := t.TempDir()

	dm, err := NewDependencyManager(client, tempDir)
	if err != nil {
		t.Fatalf("NewDependencyManager() failed: %v", err)
	}

	ctx := context.Background()

	// 测试空包名
	operation, err := dm.Add(ctx, "", "1.0.0", Production)
	if err == nil {
		t.Error("Expected error for empty package name")
	}

	if operation.Success {
		t.Error("Expected operation to fail")
	}

	// 测试移除空包名
	operation, err = dm.Remove(ctx, "")
	if err == nil {
		t.Error("Expected error for empty package name")
	}

	if operation.Success {
		t.Error("Expected operation to fail")
	}
}
