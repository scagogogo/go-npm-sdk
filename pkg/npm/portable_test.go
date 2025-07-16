package npm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewPortableManager(t *testing.T) {
	tempDir := t.TempDir()

	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	if manager == nil {
		t.Fatal("NewPortableManager() returned nil")
	}

	if manager.baseDir != tempDir {
		t.Errorf("Expected baseDir %s, got %s", tempDir, manager.baseDir)
	}

	if manager.downloader == nil {
		t.Error("Expected downloader to be initialized")
	}

	if manager.platformInfo == nil {
		t.Error("Expected platformInfo to be initialized")
	}
}

func TestNewPortableManagerWithEmptyDir(t *testing.T) {
	manager, err := NewPortableManager("")
	if err != nil {
		t.Fatalf("NewPortableManager() with empty dir failed: %v", err)
	}

	// 应该使用默认目录
	if manager.baseDir == "" {
		t.Error("Expected non-empty baseDir when empty string provided")
	}

	// 检查默认目录格式
	expectedSuffix := filepath.Join(".go-npm-sdk", "portable")
	if !strings.HasSuffix(manager.baseDir, expectedSuffix) {
		t.Errorf("Expected baseDir to end with '%s', got '%s'", expectedSuffix, manager.baseDir)
	}
}

func TestPortableConfig(t *testing.T) {
	config := &PortableConfig{
		Version:     "18.17.0",
		InstallPath: "/tmp/node-v18.17.0",
		NodePath:    "/tmp/node-v18.17.0/bin/node",
		NpmPath:     "/tmp/node-v18.17.0/bin/npm",
		InstallDate: "2023-01-01",
	}

	if config.Version != "18.17.0" {
		t.Errorf("Expected version '18.17.0', got '%s'", config.Version)
	}

	if config.InstallPath != "/tmp/node-v18.17.0" {
		t.Errorf("Expected install path '/tmp/node-v18.17.0', got '%s'", config.InstallPath)
	}
}

func TestPortableManagerSaveLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	// 创建测试配置
	installPath := filepath.Join(tempDir, "node-v18.17.0")
	err = os.MkdirAll(installPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}

	config := &PortableConfig{
		Version:     "18.17.0",
		InstallPath: installPath,
		NodePath:    manager.getNodePath(installPath),
		NpmPath:     manager.getNpmPath(installPath),
		InstallDate: time.Now().Format(time.RFC3339),
	}

	// 保存配置
	err = manager.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// 加载配置
	loadedConfig, err := manager.LoadConfig(installPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if loadedConfig.Version != config.Version {
		t.Errorf("Expected version '%s', got '%s'", config.Version, loadedConfig.Version)
	}

	if loadedConfig.InstallPath != config.InstallPath {
		t.Errorf("Expected install path '%s', got '%s'", config.InstallPath, loadedConfig.InstallPath)
	}

	if loadedConfig.NodePath != config.NodePath {
		t.Errorf("Expected node path '%s', got '%s'", config.NodePath, loadedConfig.NodePath)
	}

	if loadedConfig.NpmPath != config.NpmPath {
		t.Errorf("Expected npm path '%s', got '%s'", config.NpmPath, loadedConfig.NpmPath)
	}
}

func TestPortableManagerGetNodePath(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	installPath := "/tmp/node-v18.17.0"
	nodePath := manager.getNodePath(installPath)

	if runtime.GOOS == "windows" {
		expected := filepath.Join(installPath, "node.exe")
		if nodePath != expected {
			t.Errorf("Expected Windows node path '%s', got '%s'", expected, nodePath)
		}
	} else {
		expected := filepath.Join(installPath, "bin", "node")
		if nodePath != expected {
			t.Errorf("Expected Unix node path '%s', got '%s'", expected, nodePath)
		}
	}
}

func TestPortableManagerGetNpmPath(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	installPath := "/tmp/node-v18.17.0"
	npmPath := manager.getNpmPath(installPath)

	if runtime.GOOS == "windows" {
		expected := filepath.Join(installPath, "npm.cmd")
		if npmPath != expected {
			t.Errorf("Expected Windows npm path '%s', got '%s'", expected, npmPath)
		}
	} else {
		expected := filepath.Join(installPath, "bin", "npm")
		if npmPath != expected {
			t.Errorf("Expected Unix npm path '%s', got '%s'", expected, npmPath)
		}
	}
}

func TestPortableManagerList(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	// 初始状态应该为空
	configs, err := manager.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(configs) != 0 {
		t.Errorf("Expected empty list initially, got %d items", len(configs))
	}

	// 创建一些测试配置
	versions := []string{"16.20.0", "18.17.0"}
	for _, version := range versions {
		installPath := filepath.Join(tempDir, fmt.Sprintf("node-v%s", version))
		err = os.MkdirAll(installPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create install directory: %v", err)
		}

		config := &PortableConfig{
			Version:     version,
			InstallPath: installPath,
			NodePath:    manager.getNodePath(installPath),
			NpmPath:     manager.getNpmPath(installPath),
			InstallDate: time.Now().Format(time.RFC3339),
		}

		err = manager.SaveConfig(config)
		if err != nil {
			t.Fatalf("SaveConfig() failed: %v", err)
		}
	}

	// 再次列出配置
	configs, err = manager.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(configs) != 2 {
		t.Errorf("Expected 2 configs, got %d", len(configs))
	}

	// 验证版本
	foundVersions := make(map[string]bool)
	for _, config := range configs {
		foundVersions[config.Version] = true
	}

	for _, version := range versions {
		if !foundVersions[version] {
			t.Errorf("Expected to find version %s in list", version)
		}
	}
}

func TestPortableManagerGetConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	version := "18.17.0"
	installPath := filepath.Join(tempDir, fmt.Sprintf("node-v%s", version))
	err = os.MkdirAll(installPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}

	// 保存配置
	config := &PortableConfig{
		Version:     version,
		InstallPath: installPath,
		NodePath:    manager.getNodePath(installPath),
		NpmPath:     manager.getNpmPath(installPath),
		InstallDate: time.Now().Format(time.RFC3339),
	}

	err = manager.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// 获取配置
	retrievedConfig, err := manager.GetConfig(version)
	if err != nil {
		t.Fatalf("GetConfig() failed: %v", err)
	}

	if retrievedConfig.Version != version {
		t.Errorf("Expected version '%s', got '%s'", version, retrievedConfig.Version)
	}

	// 测试不存在的版本
	_, err = manager.GetConfig("nonexistent-version")
	if err == nil {
		t.Error("Expected error for nonexistent version")
	}
}

func TestPortableManagerIsVersionInstalled(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	version := "18.17.0"

	// 初始状态应该未安装
	if manager.IsVersionInstalled(version) {
		t.Error("Expected version to not be installed initially")
	}

	// 创建安装目录
	installPath := filepath.Join(tempDir, fmt.Sprintf("node-v%s", version))
	err = os.MkdirAll(installPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}

	// 现在应该被检测为已安装
	if !manager.IsVersionInstalled(version) {
		t.Error("Expected version to be installed after creating directory")
	}
}

func TestPortableManagerUninstall(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	version := "18.17.0"
	installPath := filepath.Join(tempDir, fmt.Sprintf("node-v%s", version))

	// 测试卸载不存在的版本
	err = manager.Uninstall(version)
	if err == nil {
		t.Error("Expected error when uninstalling nonexistent version")
	}

	// 创建安装目录和配置
	err = os.MkdirAll(installPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}

	config := &PortableConfig{
		Version:     version,
		InstallPath: installPath,
		NodePath:    manager.getNodePath(installPath),
		NpmPath:     manager.getNpmPath(installPath),
		InstallDate: time.Now().Format(time.RFC3339),
	}

	err = manager.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// 验证安装存在
	if !manager.IsVersionInstalled(version) {
		t.Error("Expected version to be installed")
	}

	// 卸载
	err = manager.Uninstall(version)
	if err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	// 验证已卸载
	if manager.IsVersionInstalled(version) {
		t.Error("Expected version to be uninstalled")
	}

	// 验证目录已删除
	if _, err := os.Stat(installPath); err == nil {
		t.Error("Expected install directory to be removed")
	}
}

func TestPortableManagerGetDefaultPath(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	defaultPath := manager.GetDefaultPath()
	expected := filepath.Join(tempDir, "default")

	if defaultPath != expected {
		t.Errorf("Expected default path '%s', got '%s'", expected, defaultPath)
	}
}

func TestPortableManagerRunCommand(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	// 测试简单命令
	err = manager.runCommand("echo", "test")
	if err != nil {
		// echo可能在某些环境中不可用
		t.Logf("runCommand failed (expected in some environments): %v", err)
	}

	// 测试不存在的命令
	err = manager.runCommand("definitely-nonexistent-command")
	if err == nil {
		t.Error("Expected error for nonexistent command")
	}
}

// Mock client for testing CreateClient
type MockNpmClient struct{}

func (m *MockNpmClient) IsAvailable(ctx context.Context) bool { return true }
func (m *MockNpmClient) Install(ctx context.Context) error    { return nil }
func (m *MockNpmClient) Version(ctx context.Context) (string, error) {
	return "8.19.2", nil
}
func (m *MockNpmClient) Init(ctx context.Context, options InitOptions) error { return nil }
func (m *MockNpmClient) InstallPackage(ctx context.Context, pkg string, options InstallOptions) error {
	return nil
}
func (m *MockNpmClient) UninstallPackage(ctx context.Context, pkg string, options UninstallOptions) error {
	return nil
}
func (m *MockNpmClient) UpdatePackage(ctx context.Context, pkg string) error { return nil }
func (m *MockNpmClient) ListPackages(ctx context.Context, options ListOptions) ([]Package, error) {
	return []Package{}, nil
}
func (m *MockNpmClient) RunScript(ctx context.Context, script string, args ...string) error {
	return nil
}
func (m *MockNpmClient) Publish(ctx context.Context, options PublishOptions) error { return nil }
func (m *MockNpmClient) GetPackageInfo(ctx context.Context, pkg string) (*PackageInfo, error) {
	return &PackageInfo{}, nil
}
func (m *MockNpmClient) Search(ctx context.Context, query string) ([]SearchResult, error) {
	return []SearchResult{}, nil
}

func TestPortableManagerCreateClient(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewPortableManager(tempDir)
	if err != nil {
		t.Fatalf("NewPortableManager() failed: %v", err)
	}

	version := "18.17.0"

	// 测试不存在的版本
	_, err = manager.CreateClient(version)
	if err == nil {
		t.Error("Expected error for nonexistent version")
	}

	// 创建配置
	installPath := filepath.Join(tempDir, fmt.Sprintf("node-v%s", version))
	err = os.MkdirAll(installPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}

	config := &PortableConfig{
		Version:     version,
		InstallPath: installPath,
		NodePath:    manager.getNodePath(installPath),
		NpmPath:     manager.getNpmPath(installPath),
		InstallDate: time.Now().Format(time.RFC3339),
	}

	err = manager.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// 现在应该能创建客户端
	client, err := manager.CreateClient(version)
	if err != nil {
		t.Fatalf("CreateClient() failed: %v", err)
	}

	if client == nil {
		t.Error("Expected non-nil client")
	}
}
