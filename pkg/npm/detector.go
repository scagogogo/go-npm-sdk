package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// NpmInfo npm信息
type NpmInfo struct {
	Version     string `json:"version"`
	Path        string `json:"path"`
	NodePath    string `json:"node_path"`
	NodeVersion string `json:"node_version"`
	Available   bool   `json:"available"`
}

// Detector npm检测器
type Detector struct {
	timeout time.Duration
}

// NewDetector 创建npm检测器
func NewDetector() *Detector {
	return &Detector{
		timeout: 10 * time.Second,
	}
}

// SetTimeout 设置超时时间
func (d *Detector) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

// Detect 检测npm是否可用
func (d *Detector) Detect(ctx context.Context) (*NpmInfo, error) {
	info := &NpmInfo{}

	// 检测npm路径
	npmPath, err := d.findNpmPath(ctx)
	if err != nil {
		return info, err
	}
	info.Path = npmPath

	// 检测npm版本
	version, err := d.getNpmVersion(ctx, npmPath)
	if err != nil {
		return info, err
	}
	info.Version = version
	info.Available = true

	// 检测Node.js路径和版本
	nodePath, nodeVersion, err := d.getNodeInfo(ctx)
	if err == nil {
		info.NodePath = nodePath
		info.NodeVersion = nodeVersion
	}

	return info, nil
}

// IsAvailable 检查npm是否可用
func (d *Detector) IsAvailable(ctx context.Context) bool {
	info, err := d.Detect(ctx)
	return err == nil && info.Available
}

// GetVersion 获取npm版本
func (d *Detector) GetVersion(ctx context.Context) (string, error) {
	npmPath, err := d.findNpmPath(ctx)
	if err != nil {
		return "", err
	}

	return d.getNpmVersion(ctx, npmPath)
}

// findNpmPath 查找npm路径
func (d *Detector) findNpmPath(ctx context.Context) (string, error) {
	// 首先尝试直接执行npm
	if _, err := d.execCommand(ctx, "npm", "--version"); err == nil {
		// 获取npm的实际路径
		if npmPath, err := d.execCommand(ctx, "which", "npm"); err == nil {
			return strings.TrimSpace(npmPath), nil
		}
		// 在Windows上使用where命令
		if runtime.GOOS == "windows" {
			if npmPath, err := d.execCommand(ctx, "where", "npm"); err == nil {
				lines := strings.Split(strings.TrimSpace(npmPath), "\n")
				if len(lines) > 0 {
					return lines[0], nil
				}
			}
		}
		return "npm", nil // 如果能执行但找不到路径，返回命令名
	}

	// 尝试常见的安装路径
	commonPaths := d.getCommonNpmPaths()
	for _, path := range commonPaths {
		if d.isExecutable(path) {
			return path, nil
		}
	}

	return "", ErrNpmNotFound
}

// getCommonNpmPaths 获取常见的npm安装路径
func (d *Detector) getCommonNpmPaths() []string {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		paths = []string{
			"C:\\Program Files\\nodejs\\npm.cmd",
			"C:\\Program Files (x86)\\nodejs\\npm.cmd",
			"C:\\Users\\%USERNAME%\\AppData\\Roaming\\npm\\npm.cmd",
			"npm.cmd",
			"npm",
		}
	case "darwin":
		paths = []string{
			"/usr/local/bin/npm",
			"/opt/homebrew/bin/npm",
			"/usr/bin/npm",
			"npm",
		}
	case "linux":
		paths = []string{
			"/usr/local/bin/npm",
			"/usr/bin/npm",
			"/opt/node/bin/npm",
			"npm",
		}
	default:
		paths = []string{"npm"}
	}

	return paths
}

// isExecutable 检查文件是否可执行
func (d *Detector) isExecutable(path string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, path, "--version")
	} else {
		cmd = exec.CommandContext(ctx, "test", "-x", path)
	}

	return cmd.Run() == nil
}

// getNpmVersion 获取npm版本
func (d *Detector) getNpmVersion(ctx context.Context, npmPath string) (string, error) {
	output, err := d.execCommand(ctx, npmPath, "--version")
	if err != nil {
		return "", fmt.Errorf("failed to get npm version: %w", err)
	}

	version := strings.TrimSpace(output)

	// 验证版本格式
	if !d.isValidVersion(version) {
		return "", fmt.Errorf("invalid version format: %s", version)
	}

	return version, nil
}

// getNodeInfo 获取Node.js信息
func (d *Detector) getNodeInfo(ctx context.Context) (string, string, error) {
	// 获取Node.js路径
	var nodePath string
	if path, err := d.execCommand(ctx, "which", "node"); err == nil {
		nodePath = strings.TrimSpace(path)
	} else if runtime.GOOS == "windows" {
		if path, err := d.execCommand(ctx, "where", "node"); err == nil {
			lines := strings.Split(strings.TrimSpace(path), "\n")
			if len(lines) > 0 {
				nodePath = lines[0]
			}
		}
	}

	// 获取Node.js版本
	version, err := d.execCommand(ctx, "node", "--version")
	if err != nil {
		return nodePath, "", err
	}

	nodeVersion := strings.TrimSpace(version)
	// 移除版本号前的'v'
	if strings.HasPrefix(nodeVersion, "v") {
		nodeVersion = nodeVersion[1:]
	}

	return nodePath, nodeVersion, nil
}

// execCommand 执行命令
func (d *Detector) execCommand(ctx context.Context, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// isValidVersion 验证版本格式
func (d *Detector) isValidVersion(version string) bool {
	// 匹配语义化版本格式 (x.y.z)
	pattern := `^\d+\.\d+\.\d+`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

// GetNpmConfig 获取npm配置
func (d *Detector) GetNpmConfig(ctx context.Context, key string) (string, error) {
	npmPath, err := d.findNpmPath(ctx)
	if err != nil {
		return "", err
	}

	output, err := d.execCommand(ctx, npmPath, "config", "get", key)
	if err != nil {
		return "", fmt.Errorf("failed to get npm config %s: %w", key, err)
	}

	return strings.TrimSpace(output), nil
}

// GetNpmConfigList 获取npm配置列表
func (d *Detector) GetNpmConfigList(ctx context.Context) (map[string]string, error) {
	npmPath, err := d.findNpmPath(ctx)
	if err != nil {
		return nil, err
	}

	output, err := d.execCommand(ctx, npmPath, "config", "list", "--json")
	if err != nil {
		return nil, fmt.Errorf("failed to get npm config list: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(output), &config); err != nil {
		return nil, fmt.Errorf("failed to parse npm config: %w", err)
	}

	result := make(map[string]string)
	for key, value := range config {
		if str, ok := value.(string); ok {
			result[key] = str
		} else {
			result[key] = fmt.Sprintf("%v", value)
		}
	}

	return result, nil
}

// GetGlobalPackagesPath 获取全局包安装路径
func (d *Detector) GetGlobalPackagesPath(ctx context.Context) (string, error) {
	path, err := d.GetNpmConfig(ctx, "prefix")
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		return filepath.Join(path, "node_modules"), nil
	}
	return filepath.Join(path, "lib", "node_modules"), nil
}

// GetCacheDir 获取npm缓存目录
func (d *Detector) GetCacheDir(ctx context.Context) (string, error) {
	return d.GetNpmConfig(ctx, "cache")
}

// GetRegistry 获取当前registry
func (d *Detector) GetRegistry(ctx context.Context) (string, error) {
	return d.GetNpmConfig(ctx, "registry")
}
