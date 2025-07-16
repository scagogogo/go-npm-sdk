package npm

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/scagogogo/go-npm-sdk/pkg/platform"
)

func TestNewInstaller(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	if installer == nil {
		t.Fatal("NewInstaller() returned nil")
	}

	if installer.detector == nil {
		t.Error("Expected detector to be initialized")
	}

	if installer.downloader == nil {
		t.Error("Expected downloader to be initialized")
	}

	if installer.platformInfo == nil {
		t.Error("Expected platformInfo to be initialized")
	}
}

func TestInstallMethodConstants(t *testing.T) {
	// 测试安装方法常量
	if PackageManager != "package_manager" {
		t.Errorf("Expected PackageManager to be 'package_manager', got '%s'", PackageManager)
	}

	if OfficialInstaller != "official_installer" {
		t.Errorf("Expected OfficialInstaller to be 'official_installer', got '%s'", OfficialInstaller)
	}

	if Portable != "portable" {
		t.Errorf("Expected Portable to be 'portable', got '%s'", Portable)
	}

	if Manual != "manual" {
		t.Errorf("Expected Manual to be 'manual', got '%s'", Manual)
	}
}

func TestNpmInstallOptions(t *testing.T) {
	options := NpmInstallOptions{
		Method:      PackageManager,
		Version:     "18.17.0",
		InstallPath: "/tmp/node",
		Force:       true,
		Global:      false,
		Progress:    func(msg string) { t.Logf("Progress: %s", msg) },
	}

	if options.Method != PackageManager {
		t.Errorf("Expected method PackageManager, got %s", options.Method)
	}

	if options.Version != "18.17.0" {
		t.Errorf("Expected version '18.17.0', got '%s'", options.Version)
	}

	if !options.Force {
		t.Error("Expected Force to be true")
	}

	if options.Global {
		t.Error("Expected Global to be false")
	}

	// 测试进度回调
	if options.Progress != nil {
		options.Progress("test message")
	}
}

func TestInstallResult(t *testing.T) {
	result := &InstallResult{
		Success:  true,
		Method:   PackageManager,
		Version:  "8.19.2",
		Path:     "/usr/local/bin/npm",
		Duration: 30 * time.Second,
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.Method != PackageManager {
		t.Errorf("Expected method PackageManager, got %s", result.Method)
	}

	if result.Duration != 30*time.Second {
		t.Errorf("Expected duration 30s, got %v", result.Duration)
	}
}

func TestInstallerHasCommand(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	// 测试存在的命令
	if !installer.hasCommand("echo") {
		// echo在某些系统上可能不存在，这是正常的
		t.Logf("echo command not found (normal on some systems)")
	}

	// 测试不存在的命令
	if installer.hasCommand("definitely-nonexistent-command-12345") {
		t.Error("Expected nonexistent command to return false")
	}
}

func TestInstallerHasPackageManager(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	hasPackageManager := installer.hasPackageManager()
	t.Logf("Has package manager: %v", hasPackageManager)

	// 根据平台验证结果
	switch installer.platformInfo.Platform {
	case platform.Windows:
		// Windows上可能有choco或winget
		t.Logf("Windows platform detected")
	case platform.MacOS:
		// macOS上可能有brew或port
		t.Logf("macOS platform detected")
	case platform.Linux:
		// Linux通常都有包管理器
		if !hasPackageManager {
			t.Log("Linux platform but no package manager detected (unusual but possible)")
		}
	default:
		t.Logf("Unknown platform: %s", installer.platformInfo.Platform)
	}
}

func TestInstallerGetPortableNpmPath(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	installPath := "/tmp/node"
	npmPath := installer.getPortableNpmPath(installPath)

	if npmPath == "" {
		t.Error("Expected non-empty npm path")
	}

	// 验证路径格式
	switch installer.platformInfo.Platform {
	case platform.Windows:
		if !strings.HasSuffix(npmPath, "npm.cmd") {
			t.Errorf("Expected Windows npm path to end with 'npm.cmd', got '%s'", npmPath)
		}
	default:
		if !strings.Contains(npmPath, "bin/npm") {
			t.Errorf("Expected Unix npm path to contain 'bin/npm', got '%s'", npmPath)
		}
	}

	t.Logf("Portable npm path: %s", npmPath)
}

// MockDetector for testing
type MockDetector struct {
	available bool
	version   string
	path      string
}

func (m *MockDetector) IsAvailable(ctx context.Context) bool {
	return m.available
}

func (m *MockDetector) Detect(ctx context.Context) (*NpmInfo, error) {
	if !m.available {
		return nil, ErrNpmNotFound
	}
	return &NpmInfo{
		Version:   m.version,
		Path:      m.path,
		Available: true,
	}, nil
}

func TestInstallerInstallWithExistingNpm(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	// 由于detector是私有字段，我们无法直接替换
	// 但我们可以测试安装选项的处理逻辑
	_ = &MockDetector{
		available: true,
		version:   "8.19.2",
		path:      "/usr/local/bin/npm",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := NpmInstallOptions{
		Method: PackageManager,
		Force:  false, // 不强制安装
	}

	// 这个测试主要验证方法调用不会panic
	result, err := installer.Install(ctx, options)

	// 结果取决于系统是否真的有npm和包管理器
	t.Logf("Install result: success=%v, error=%v", result != nil && result.Success, err)

	if result != nil {
		t.Logf("Install method used: %s", result.Method)
		t.Logf("Install duration: %v", result.Duration)
	}
}

func TestInstallerInstallWithForce(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := NpmInstallOptions{
		Method: PackageManager,
		Force:  true, // 强制安装
	}

	// 这个测试主要验证强制安装的逻辑
	result, err := installer.Install(ctx, options)

	t.Logf("Force install result: success=%v, error=%v", result != nil && result.Success, err)

	if result != nil {
		t.Logf("Force install method used: %s", result.Method)
	}
}

func TestInstallerInstallPortable(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := NpmInstallOptions{
		Method:      Portable,
		Version:     "18.17.0",
		InstallPath: "/tmp/test-node",
	}

	// 这个测试会尝试下载，但可能因为网络问题失败
	result, err := installer.installPortable(ctx, options)

	// 我们不期望成功，因为这需要真实的下载
	t.Logf("Portable install result: success=%v, error=%v", result != nil && result.Success, err)

	if err != nil {
		// 检查错误类型
		if result != nil && !result.Success {
			t.Logf("Expected failure for portable install without network: %v", result.Error)
		}
	}
}

func TestInstallerInstallAuto(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := NpmInstallOptions{
		Method: "", // 空方法，应该触发自动选择
	}

	// 测试自动选择安装方法
	result, err := installer.installAuto(ctx, options)

	t.Logf("Auto install result: success=%v, error=%v", result != nil && result.Success, err)

	if result != nil {
		t.Logf("Auto selected method: %s", result.Method)
	}
}

func TestInstallerValidation(t *testing.T) {
	installer, err := NewInstaller()
	if err != nil {
		t.Fatalf("NewInstaller() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试便携版安装但没有指定路径
	options := NpmInstallOptions{
		Method:      Portable,
		InstallPath: "", // 空路径
	}

	result, err := installer.installPortable(ctx, options)
	if err == nil {
		t.Error("Expected error for portable install without path")
	}

	if result != nil && result.Success {
		t.Error("Expected failure for portable install without path")
	}

	t.Logf("Validation test result: %v", err)
}
