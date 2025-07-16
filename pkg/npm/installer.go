package npm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/scagogogo/go-npm-sdk/pkg/platform"
)

// InstallMethod 安装方法
type InstallMethod string

const (
	PackageManager    InstallMethod = "package_manager"
	OfficialInstaller InstallMethod = "official_installer"
	Portable          InstallMethod = "portable"
	Manual            InstallMethod = "manual"
)

// NpmInstallOptions npm安装选项
type NpmInstallOptions struct {
	Method      InstallMethod `json:"method"`
	Version     string        `json:"version"`      // 指定版本，空表示最新版
	InstallPath string        `json:"install_path"` // 安装路径（便携版使用）
	Force       bool          `json:"force"`        // 强制安装
	Global      bool          `json:"global"`       // 全局安装
	Progress    func(string)  `json:"-"`            // 进度回调
}

// InstallResult 安装结果
type InstallResult struct {
	Success  bool          `json:"success"`
	Method   InstallMethod `json:"method"`
	Version  string        `json:"version"`
	Path     string        `json:"path"`
	Duration time.Duration `json:"duration"`
	Error    error         `json:"error,omitempty"`
}

// Installer npm安装器
type Installer struct {
	detector     *Detector
	downloader   *platform.NodeJSDownloader
	platformInfo *platform.Info
}

// NewInstaller 创建npm安装器
func NewInstaller() (*Installer, error) {
	detector := platform.NewDetector()
	info, err := detector.Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to detect platform: %w", err)
	}

	return &Installer{
		detector:     NewDetector(),
		downloader:   platform.NewNodeJSDownloader(),
		platformInfo: info,
	}, nil
}

// Install 安装npm
func (i *Installer) Install(ctx context.Context, options NpmInstallOptions) (*InstallResult, error) {
	startTime := time.Now()

	// 如果已安装且不强制安装，直接返回
	if !options.Force && i.detector.IsAvailable(ctx) {
		info, _ := i.detector.Detect(ctx)
		return &InstallResult{
			Success:  true,
			Method:   Manual,
			Version:  info.Version,
			Path:     info.Path,
			Duration: time.Since(startTime),
		}, nil
	}

	// 根据安装方法进行安装
	var result *InstallResult
	var err error

	switch options.Method {
	case PackageManager:
		result, err = i.installViaPackageManager(ctx, options)
	case OfficialInstaller:
		result, err = i.installViaOfficialInstaller(ctx, options)
	case Portable:
		result, err = i.installPortable(ctx, options)
	default:
		// 自动选择最佳安装方法
		result, err = i.installAuto(ctx, options)
	}

	if result != nil {
		result.Duration = time.Since(startTime)
	}

	return result, err
}

// installAuto 自动选择安装方法
func (i *Installer) installAuto(ctx context.Context, options NpmInstallOptions) (*InstallResult, error) {
	// 优先尝试包管理器
	if i.hasPackageManager() {
		if result, err := i.installViaPackageManager(ctx, options); err == nil {
			return result, nil
		}
	}

	// 然后尝试便携版
	if options.InstallPath != "" {
		return i.installPortable(ctx, options)
	}

	// 最后尝试官方安装程序
	return i.installViaOfficialInstaller(ctx, options)
}

// installViaPackageManager 通过包管理器安装
func (i *Installer) installViaPackageManager(ctx context.Context, options NpmInstallOptions) (*InstallResult, error) {
	if options.Progress != nil {
		options.Progress("正在通过包管理器安装Node.js/npm...")
	}

	var cmd *exec.Cmd
	var packageName string

	switch i.platformInfo.Platform {
	case platform.Windows:
		// 尝试Chocolatey
		if i.hasCommand("choco") {
			packageName = "nodejs"
			cmd = exec.CommandContext(ctx, "choco", "install", packageName, "-y")
		} else if i.hasCommand("winget") {
			packageName = "OpenJS.NodeJS"
			cmd = exec.CommandContext(ctx, "winget", "install", packageName)
		} else {
			return nil, fmt.Errorf("no package manager found on Windows")
		}

	case platform.MacOS:
		// 尝试Homebrew
		if i.hasCommand("brew") {
			packageName = "node"
			cmd = exec.CommandContext(ctx, "brew", "install", packageName)
		} else if i.hasCommand("port") {
			packageName = "nodejs18"
			cmd = exec.CommandContext(ctx, "sudo", "port", "install", packageName)
		} else {
			return nil, fmt.Errorf("no package manager found on macOS")
		}

	case platform.Linux:
		packageName = "nodejs npm"
		switch i.platformInfo.Distribution {
		case platform.Ubuntu, platform.Debian:
			cmd = exec.CommandContext(ctx, "sudo", "apt-get", "update")
			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("failed to update package list: %w", err)
			}
			cmd = exec.CommandContext(ctx, "sudo", "apt-get", "install", "-y", "nodejs", "npm")
		case platform.CentOS, platform.RHEL, platform.Fedora:
			if i.hasCommand("dnf") {
				cmd = exec.CommandContext(ctx, "sudo", "dnf", "install", "-y", "nodejs", "npm")
			} else {
				cmd = exec.CommandContext(ctx, "sudo", "yum", "install", "-y", "nodejs", "npm")
			}
		case platform.Arch:
			cmd = exec.CommandContext(ctx, "sudo", "pacman", "-S", "--noconfirm", "nodejs", "npm")
		case platform.Alpine:
			cmd = exec.CommandContext(ctx, "sudo", "apk", "add", "nodejs", "npm")
		case platform.SUSE:
			cmd = exec.CommandContext(ctx, "sudo", "zypper", "install", "-y", "nodejs", "npm")
		default:
			return nil, fmt.Errorf("unsupported Linux distribution: %s", i.platformInfo.Distribution)
		}

	default:
		return nil, NewPlatformError(string(i.platformInfo.Platform), "unsupported platform for package manager installation", nil)
	}

	if options.Progress != nil {
		options.Progress(fmt.Sprintf("执行安装命令: %s", cmd.String()))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &InstallResult{
			Success: false,
			Method:  PackageManager,
			Error:   fmt.Errorf("package manager installation failed: %w\nOutput: %s", err, string(output)),
		}, err
	}

	// 验证安装
	if !i.detector.IsAvailable(ctx) {
		return &InstallResult{
			Success: false,
			Method:  PackageManager,
			Error:   fmt.Errorf("npm not available after package manager installation"),
		}, fmt.Errorf("npm not available after installation")
	}

	info, _ := i.detector.Detect(ctx)
	return &InstallResult{
		Success: true,
		Method:  PackageManager,
		Version: info.Version,
		Path:    info.Path,
	}, nil
}

// installViaOfficialInstaller 通过官方安装程序安装
func (i *Installer) installViaOfficialInstaller(ctx context.Context, options NpmInstallOptions) (*InstallResult, error) {
	if options.Progress != nil {
		options.Progress("正在下载官方安装程序...")
	}

	// 获取版本
	version := options.Version
	if version == "" {
		var err error
		version, err = i.downloader.GetLatestVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", err)
		}
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "nodejs-installer-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 下载安装程序
	progress := func(downloaded, total int64) {
		if options.Progress != nil {
			percent := float64(downloaded) / float64(total) * 100
			options.Progress(fmt.Sprintf("下载进度: %.1f%%", percent))
		}
	}

	result, err := i.downloader.DownloadNodeJS(ctx, version, i.platformInfo, tempDir, progress)
	if err != nil {
		return &InstallResult{
			Success: false,
			Method:  OfficialInstaller,
			Error:   fmt.Errorf("failed to download Node.js: %w", err),
		}, err
	}

	if options.Progress != nil {
		options.Progress("正在安装Node.js...")
	}

	// 执行安装
	if err := i.executeInstaller(ctx, result.FilePath); err != nil {
		return &InstallResult{
			Success: false,
			Method:  OfficialInstaller,
			Error:   fmt.Errorf("failed to execute installer: %w", err),
		}, err
	}

	// 验证安装
	if !i.detector.IsAvailable(ctx) {
		return &InstallResult{
			Success: false,
			Method:  OfficialInstaller,
			Error:   fmt.Errorf("npm not available after installation"),
		}, fmt.Errorf("npm not available after installation")
	}

	info, _ := i.detector.Detect(ctx)
	return &InstallResult{
		Success: true,
		Method:  OfficialInstaller,
		Version: info.Version,
		Path:    info.Path,
	}, nil
}

// installPortable 安装便携版
func (i *Installer) installPortable(ctx context.Context, options NpmInstallOptions) (*InstallResult, error) {
	if options.InstallPath == "" {
		return nil, fmt.Errorf("install path is required for portable installation")
	}

	if options.Progress != nil {
		options.Progress("正在下载便携版Node.js...")
	}

	// 获取版本
	version := options.Version
	if version == "" {
		var err error
		version, err = i.downloader.GetLatestVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", err)
		}
	}

	// 下载便携版
	progress := func(downloaded, total int64) {
		if options.Progress != nil {
			percent := float64(downloaded) / float64(total) * 100
			options.Progress(fmt.Sprintf("下载进度: %.1f%%", percent))
		}
	}

	tempDir, err := os.MkdirTemp("", "nodejs-portable-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	result, err := i.downloader.DownloadNodeJS(ctx, version, i.platformInfo, tempDir, progress)
	if err != nil {
		return &InstallResult{
			Success: false,
			Method:  Portable,
			Error:   fmt.Errorf("failed to download portable Node.js: %w", err),
		}, err
	}

	if options.Progress != nil {
		options.Progress("正在解压便携版...")
	}

	// 解压到目标目录
	if err := i.extractPortable(result.FilePath, options.InstallPath); err != nil {
		return &InstallResult{
			Success: false,
			Method:  Portable,
			Error:   fmt.Errorf("failed to extract portable Node.js: %w", err),
		}, err
	}

	// 获取npm路径
	npmPath := i.getPortableNpmPath(options.InstallPath)

	return &InstallResult{
		Success: true,
		Method:  Portable,
		Version: version,
		Path:    npmPath,
	}, nil
}

// hasPackageManager 检查是否有包管理器
func (i *Installer) hasPackageManager() bool {
	switch i.platformInfo.Platform {
	case platform.Windows:
		return i.hasCommand("choco") || i.hasCommand("winget")
	case platform.MacOS:
		return i.hasCommand("brew") || i.hasCommand("port")
	case platform.Linux:
		return true // Linux通常都有包管理器
	default:
		return false
	}
}

// hasCommand 检查命令是否存在
func (i *Installer) hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// executeInstaller 执行安装程序
func (i *Installer) executeInstaller(ctx context.Context, installerPath string) error {
	var cmd *exec.Cmd

	switch i.platformInfo.Platform {
	case platform.Windows:
		if strings.HasSuffix(installerPath, ".msi") {
			cmd = exec.CommandContext(ctx, "msiexec", "/i", installerPath, "/quiet")
		} else {
			cmd = exec.CommandContext(ctx, installerPath, "/S")
		}
	case platform.MacOS:
		if strings.HasSuffix(installerPath, ".pkg") {
			cmd = exec.CommandContext(ctx, "sudo", "installer", "-pkg", installerPath, "-target", "/")
		} else {
			return fmt.Errorf("unsupported installer format for macOS")
		}
	default:
		return fmt.Errorf("official installer not supported on %s", i.platformInfo.Platform)
	}

	return cmd.Run()
}

// extractPortable 解压便携版
func (i *Installer) extractPortable(archivePath, destPath string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(archivePath, ".zip"):
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "Expand-Archive", "-Path", archivePath, "-DestinationPath", destPath)
		} else {
			cmd = exec.Command("unzip", "-q", archivePath, "-d", destPath)
		}
	case strings.HasSuffix(archivePath, ".tar.gz"):
		cmd = exec.Command("tar", "-xzf", archivePath, "-C", destPath)
	case strings.HasSuffix(archivePath, ".tar.xz"):
		cmd = exec.Command("tar", "-xJf", archivePath, "-C", destPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}

	return cmd.Run()
}

// getPortableNpmPath 获取便携版npm路径
func (i *Installer) getPortableNpmPath(installPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(installPath, "npm.cmd")
	}
	return filepath.Join(installPath, "bin", "npm")
}
