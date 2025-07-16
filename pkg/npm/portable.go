package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/scagogogo/go-npm-sdk/pkg/platform"
)

// PortableManager 便携版管理器
type PortableManager struct {
	downloader   *platform.NodeJSDownloader
	platformInfo *platform.Info
	baseDir      string
}

// PortableConfig 便携版配置
type PortableConfig struct {
	Version     string `json:"version"`
	InstallPath string `json:"install_path"`
	NodePath    string `json:"node_path"`
	NpmPath     string `json:"npm_path"`
	InstallDate string `json:"install_date"`
}

// NewPortableManager 创建便携版管理器
func NewPortableManager(baseDir string) (*PortableManager, error) {
	detector := platform.NewDetector()
	info, err := detector.Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to detect platform: %w", err)
	}

	if baseDir == "" {
		// 使用默认目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".go-npm-sdk", "portable")
	}

	return &PortableManager{
		downloader:   platform.NewNodeJSDownloader(),
		platformInfo: info,
		baseDir:      baseDir,
	}, nil
}

// Install 安装便携版Node.js/npm
func (pm *PortableManager) Install(ctx context.Context, version string, progress func(string)) (*PortableConfig, error) {
	if progress != nil {
		progress("开始安装便携版Node.js...")
	}

	// 获取版本
	if version == "" {
		var err error
		version, err = pm.downloader.GetLatestVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", err)
		}
	}

	// 检查是否已安装
	installPath := filepath.Join(pm.baseDir, fmt.Sprintf("node-v%s", version))
	if config, err := pm.LoadConfig(installPath); err == nil {
		if progress != nil {
			progress(fmt.Sprintf("版本 %s 已安装", version))
		}
		return config, nil
	}

	// 创建安装目录
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create install directory: %w", err)
	}

	// 下载Node.js
	if progress != nil {
		progress("正在下载Node.js...")
	}

	downloadProgress := func(downloaded, total int64) {
		if progress != nil {
			percent := float64(downloaded) / float64(total) * 100
			progress(fmt.Sprintf("下载进度: %.1f%%", percent))
		}
	}

	tempDir, err := os.MkdirTemp("", "nodejs-portable-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	result, err := pm.downloader.DownloadNodeJS(ctx, version, pm.platformInfo, tempDir, downloadProgress)
	if err != nil {
		return nil, fmt.Errorf("failed to download Node.js: %w", err)
	}

	// 解压
	if progress != nil {
		progress("正在解压...")
	}

	if err := pm.extractArchive(result.FilePath, installPath); err != nil {
		return nil, fmt.Errorf("failed to extract archive: %w", err)
	}

	// 创建配置
	config := &PortableConfig{
		Version:     version,
		InstallPath: installPath,
		NodePath:    pm.getNodePath(installPath),
		NpmPath:     pm.getNpmPath(installPath),
		InstallDate: result.Duration.String(),
	}

	// 保存配置
	if err := pm.SaveConfig(config); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	if progress != nil {
		progress(fmt.Sprintf("便携版Node.js v%s 安装完成", version))
	}

	return config, nil
}

// Uninstall 卸载便携版
func (pm *PortableManager) Uninstall(version string) error {
	installPath := filepath.Join(pm.baseDir, fmt.Sprintf("node-v%s", version))

	// 检查是否存在
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	// 删除安装目录
	return os.RemoveAll(installPath)
}

// List 列出已安装的版本
func (pm *PortableManager) List() ([]*PortableConfig, error) {
	if _, err := os.Stat(pm.baseDir); os.IsNotExist(err) {
		return []*PortableConfig{}, nil
	}

	entries, err := os.ReadDir(pm.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read base directory: %w", err)
	}

	var configs []*PortableConfig
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "node-v") {
			continue
		}

		installPath := filepath.Join(pm.baseDir, entry.Name())
		if config, err := pm.LoadConfig(installPath); err == nil {
			configs = append(configs, config)
		}
	}

	return configs, nil
}

// GetConfig 获取指定版本的配置
func (pm *PortableManager) GetConfig(version string) (*PortableConfig, error) {
	installPath := filepath.Join(pm.baseDir, fmt.Sprintf("node-v%s", version))
	return pm.LoadConfig(installPath)
}

// LoadConfig 加载配置
func (pm *PortableManager) LoadConfig(installPath string) (*PortableConfig, error) {
	configPath := filepath.Join(installPath, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config PortableConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存配置
func (pm *PortableManager) SaveConfig(config *PortableConfig) error {
	configPath := filepath.Join(config.InstallPath, "config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// extractArchive 解压归档文件
func (pm *PortableManager) extractArchive(archivePath, destPath string) error {
	// 根据文件扩展名选择解压方法
	switch {
	case strings.HasSuffix(archivePath, ".zip"):
		return pm.extractZip(archivePath, destPath)
	case strings.HasSuffix(archivePath, ".tar.gz"):
		return pm.extractTarGz(archivePath, destPath)
	case strings.HasSuffix(archivePath, ".tar.xz"):
		return pm.extractTarXz(archivePath, destPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}
}

// extractZip 解压ZIP文件
func (pm *PortableManager) extractZip(archivePath, destPath string) error {
	if runtime.GOOS == "windows" {
		// 使用PowerShell解压
		cmd := fmt.Sprintf(`Expand-Archive -Path "%s" -DestinationPath "%s" -Force`, archivePath, destPath)
		return pm.runCommand("powershell", "-Command", cmd)
	} else {
		// 使用unzip命令
		return pm.runCommand("unzip", "-q", archivePath, "-d", destPath)
	}
}

// extractTarGz 解压tar.gz文件
func (pm *PortableManager) extractTarGz(archivePath, destPath string) error {
	return pm.runCommand("tar", "-xzf", archivePath, "-C", destPath, "--strip-components=1")
}

// extractTarXz 解压tar.xz文件
func (pm *PortableManager) extractTarXz(archivePath, destPath string) error {
	return pm.runCommand("tar", "-xJf", archivePath, "-C", destPath, "--strip-components=1")
}

// runCommand 运行命令
func (pm *PortableManager) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// getNodePath 获取Node.js可执行文件路径
func (pm *PortableManager) getNodePath(installPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(installPath, "node.exe")
	}
	return filepath.Join(installPath, "bin", "node")
}

// getNpmPath 获取npm可执行文件路径
func (pm *PortableManager) getNpmPath(installPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(installPath, "npm.cmd")
	}
	return filepath.Join(installPath, "bin", "npm")
}

// CreateClient 为指定版本创建npm客户端
func (pm *PortableManager) CreateClient(version string) (Client, error) {
	config, err := pm.GetConfig(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get config for version %s: %w", version, err)
	}

	// 创建使用便携版npm的客户端
	return NewClientWithPath(config.NpmPath)
}

// SetAsDefault 将指定版本设置为默认版本
func (pm *PortableManager) SetAsDefault(version string) error {
	config, err := pm.GetConfig(version)
	if err != nil {
		return fmt.Errorf("failed to get config for version %s: %w", version, err)
	}

	// 创建符号链接或更新PATH
	defaultPath := filepath.Join(pm.baseDir, "default")

	// 删除现有的默认链接
	os.RemoveAll(defaultPath)

	// 创建新的符号链接
	if runtime.GOOS == "windows" {
		// Windows上创建目录链接
		return pm.runCommand("mklink", "/D", defaultPath, config.InstallPath)
	} else {
		// Unix系统上创建符号链接
		return os.Symlink(config.InstallPath, defaultPath)
	}
}

// GetDefaultPath 获取默认版本的路径
func (pm *PortableManager) GetDefaultPath() string {
	return filepath.Join(pm.baseDir, "default")
}

// IsVersionInstalled 检查版本是否已安装
func (pm *PortableManager) IsVersionInstalled(version string) bool {
	installPath := filepath.Join(pm.baseDir, fmt.Sprintf("node-v%s", version))
	_, err := os.Stat(installPath)
	return err == nil
}
