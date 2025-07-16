package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Platform 平台类型
type Platform string

const (
	Windows Platform = "windows"
	MacOS   Platform = "darwin"
	Linux   Platform = "linux"
	Unknown Platform = "unknown"
)

// Architecture 架构类型
type Architecture string

const (
	AMD64 Architecture = "amd64"
	ARM64 Architecture = "arm64"
	I386  Architecture = "386"
	ARM   Architecture = "arm"
)

// Distribution Linux发行版
type Distribution string

const (
	Ubuntu   Distribution = "ubuntu"
	Debian   Distribution = "debian"
	CentOS   Distribution = "centos"
	RHEL     Distribution = "rhel"
	Fedora   Distribution = "fedora"
	SUSE     Distribution = "suse"
	Arch     Distribution = "arch"
	Alpine   Distribution = "alpine"
	UnknownDistro Distribution = "unknown"
)

// Info 平台信息
type Info struct {
	Platform     Platform     `json:"platform"`
	Architecture Architecture `json:"architecture"`
	Distribution Distribution `json:"distribution,omitempty"`
	Version      string       `json:"version,omitempty"`
	Kernel       string       `json:"kernel,omitempty"`
}

// Detector 平台检测器
type Detector struct{}

// NewDetector 创建新的平台检测器
func NewDetector() *Detector {
	return &Detector{}
}

// Detect 检测当前平台信息
func (d *Detector) Detect() (*Info, error) {
	info := &Info{
		Platform:     Platform(runtime.GOOS),
		Architecture: Architecture(runtime.GOARCH),
	}

	// 检测Linux发行版
	if info.Platform == Linux {
		dist, version, err := d.detectLinuxDistribution()
		if err == nil {
			info.Distribution = dist
			info.Version = version
		}
	}

	// 检测系统版本
	if version, err := d.detectSystemVersion(); err == nil {
		if info.Version == "" {
			info.Version = version
		}
	}

	// 检测内核版本
	if kernel, err := d.detectKernelVersion(); err == nil {
		info.Kernel = kernel
	}

	return info, nil
}

// detectLinuxDistribution 检测Linux发行版
func (d *Detector) detectLinuxDistribution() (Distribution, string, error) {
	// 尝试读取 /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		return d.parseOSRelease(string(data))
	}

	// 尝试读取 /etc/lsb-release
	if data, err := os.ReadFile("/etc/lsb-release"); err == nil {
		return d.parseLSBRelease(string(data))
	}

	// 检查特定发行版文件
	distFiles := map[string]Distribution{
		"/etc/redhat-release": CentOS,
		"/etc/centos-release": CentOS,
		"/etc/fedora-release": Fedora,
		"/etc/debian_version": Debian,
		"/etc/arch-release":   Arch,
		"/etc/alpine-release": Alpine,
	}

	for file, dist := range distFiles {
		if _, err := os.Stat(file); err == nil {
			if data, err := os.ReadFile(file); err == nil {
				return dist, strings.TrimSpace(string(data)), nil
			}
			return dist, "", nil
		}
	}

	return UnknownDistro, "", fmt.Errorf("unable to detect Linux distribution")
}

// parseOSRelease 解析 /etc/os-release 文件
func (d *Detector) parseOSRelease(content string) (Distribution, string, error) {
	lines := strings.Split(content, "\n")
	var id, version string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}

	dist := d.mapDistributionID(id)
	return dist, version, nil
}

// parseLSBRelease 解析 /etc/lsb-release 文件
func (d *Detector) parseLSBRelease(content string) (Distribution, string, error) {
	lines := strings.Split(content, "\n")
	var id, version string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "DISTRIB_ID=") {
			id = strings.ToLower(strings.TrimPrefix(line, "DISTRIB_ID="))
		} else if strings.HasPrefix(line, "DISTRIB_RELEASE=") {
			version = strings.TrimPrefix(line, "DISTRIB_RELEASE=")
		}
	}

	dist := d.mapDistributionID(id)
	return dist, version, nil
}

// mapDistributionID 映射发行版ID到Distribution类型
func (d *Detector) mapDistributionID(id string) Distribution {
	id = strings.ToLower(id)
	switch {
	case strings.Contains(id, "ubuntu"):
		return Ubuntu
	case strings.Contains(id, "debian"):
		return Debian
	case strings.Contains(id, "centos"):
		return CentOS
	case strings.Contains(id, "rhel") || strings.Contains(id, "redhat"):
		return RHEL
	case strings.Contains(id, "fedora"):
		return Fedora
	case strings.Contains(id, "suse"):
		return SUSE
	case strings.Contains(id, "arch"):
		return Arch
	case strings.Contains(id, "alpine"):
		return Alpine
	default:
		return UnknownDistro
	}
}

// detectSystemVersion 检测系统版本
func (d *Detector) detectSystemVersion() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return d.detectWindowsVersion()
	case "darwin":
		return d.detectMacOSVersion()
	case "linux":
		return d.detectLinuxVersion()
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// detectWindowsVersion 检测Windows版本
func (d *Detector) detectWindowsVersion() (string, error) {
	cmd := exec.Command("cmd", "/c", "ver")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// detectMacOSVersion 检测macOS版本
func (d *Detector) detectMacOSVersion() (string, error) {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// detectLinuxVersion 检测Linux版本
func (d *Detector) detectLinuxVersion() (string, error) {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// detectKernelVersion 检测内核版本
func (d *Detector) detectKernelVersion() (string, error) {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// IsWindows 检查是否为Windows平台
func (info *Info) IsWindows() bool {
	return info.Platform == Windows
}

// IsMacOS 检查是否为macOS平台
func (info *Info) IsMacOS() bool {
	return info.Platform == MacOS
}

// IsLinux 检查是否为Linux平台
func (info *Info) IsLinux() bool {
	return info.Platform == Linux
}

// IsARM 检查是否为ARM架构
func (info *Info) IsARM() bool {
	return info.Architecture == ARM64 || info.Architecture == ARM
}

// IsX86 检查是否为x86架构
func (info *Info) IsX86() bool {
	return info.Architecture == AMD64 || info.Architecture == I386
}

// String 返回平台信息的字符串表示
func (info *Info) String() string {
	result := fmt.Sprintf("%s/%s", info.Platform, info.Architecture)
	if info.Distribution != "" && info.Distribution != UnknownDistro {
		result += fmt.Sprintf(" (%s", info.Distribution)
		if info.Version != "" {
			result += fmt.Sprintf(" %s", info.Version)
		}
		result += ")"
	} else if info.Version != "" {
		result += fmt.Sprintf(" (%s)", info.Version)
	}
	return result
}
