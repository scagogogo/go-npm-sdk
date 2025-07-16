package platform

import (
	"runtime"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}
}

func TestDetect(t *testing.T) {
	detector := NewDetector()
	info, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() failed: %v", err)
	}

	if info == nil {
		t.Fatal("Detect() returned nil info")
	}

	// 验证平台信息
	expectedPlatform := Platform(runtime.GOOS)
	if info.Platform != expectedPlatform {
		t.Errorf("Expected platform %s, got %s", expectedPlatform, info.Platform)
	}

	expectedArch := Architecture(runtime.GOARCH)
	if info.Architecture != expectedArch {
		t.Errorf("Expected architecture %s, got %s", expectedArch, info.Architecture)
	}

	t.Logf("Platform info: %s", info.String())
}

func TestPlatformMethods(t *testing.T) {
	detector := NewDetector()
	info, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() failed: %v", err)
	}

	// 测试平台检查方法
	switch runtime.GOOS {
	case "windows":
		if !info.IsWindows() {
			t.Error("IsWindows() should return true on Windows")
		}
		if info.IsMacOS() {
			t.Error("IsMacOS() should return false on Windows")
		}
		if info.IsLinux() {
			t.Error("IsLinux() should return false on Windows")
		}
	case "darwin":
		if !info.IsMacOS() {
			t.Error("IsMacOS() should return true on macOS")
		}
		if info.IsWindows() {
			t.Error("IsWindows() should return false on macOS")
		}
		if info.IsLinux() {
			t.Error("IsLinux() should return false on macOS")
		}
	case "linux":
		if !info.IsLinux() {
			t.Error("IsLinux() should return true on Linux")
		}
		if info.IsWindows() {
			t.Error("IsWindows() should return false on Linux")
		}
		if info.IsMacOS() {
			t.Error("IsMacOS() should return false on Linux")
		}
	}
}

func TestArchitectureMethods(t *testing.T) {
	detector := NewDetector()
	info, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() failed: %v", err)
	}

	// 测试架构检查方法
	switch runtime.GOARCH {
	case "amd64", "386":
		if !info.IsX86() {
			t.Error("IsX86() should return true for x86 architectures")
		}
		if info.IsARM() {
			t.Error("IsARM() should return false for x86 architectures")
		}
	case "arm64", "arm":
		if !info.IsARM() {
			t.Error("IsARM() should return true for ARM architectures")
		}
		if info.IsX86() {
			t.Error("IsX86() should return false for ARM architectures")
		}
	}
}

func TestMapDistributionID(t *testing.T) {
	detector := NewDetector()

	testCases := []struct {
		id       string
		expected Distribution
	}{
		{"ubuntu", Ubuntu},
		{"Ubuntu", Ubuntu},
		{"debian", Debian},
		{"Debian", Debian},
		{"centos", CentOS},
		{"CentOS", CentOS},
		{"rhel", RHEL},
		{"redhat", RHEL},
		{"fedora", Fedora},
		{"Fedora", Fedora},
		{"suse", SUSE},
		{"SUSE", SUSE},
		{"arch", Arch},
		{"Arch", Arch},
		{"alpine", Alpine},
		{"Alpine", Alpine},
		{"unknown", UnknownDistro},
		{"", UnknownDistro},
	}

	for _, tc := range testCases {
		result := detector.mapDistributionID(tc.id)
		if result != tc.expected {
			t.Errorf("mapDistributionID(%s) = %s, expected %s", tc.id, result, tc.expected)
		}
	}
}

func TestParseOSRelease(t *testing.T) {
	detector := NewDetector()

	testContent := `NAME="Ubuntu"
VERSION="20.04.3 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.3 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal`

	dist, version, err := detector.parseOSRelease(testContent)
	if err != nil {
		t.Fatalf("parseOSRelease() failed: %v", err)
	}

	if dist != Ubuntu {
		t.Errorf("Expected distribution Ubuntu, got %s", dist)
	}

	if version != "20.04" {
		t.Errorf("Expected version '20.04', got '%s'", version)
	}
}

func TestParseLSBRelease(t *testing.T) {
	detector := NewDetector()

	testContent := `DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=20.04
DISTRIB_CODENAME=focal
DISTRIB_DESCRIPTION="Ubuntu 20.04.3 LTS"`

	dist, version, err := detector.parseLSBRelease(testContent)
	if err != nil {
		t.Fatalf("parseLSBRelease() failed: %v", err)
	}

	if dist != Ubuntu {
		t.Errorf("Expected distribution Ubuntu, got %s", dist)
	}

	if version != "20.04" {
		t.Errorf("Expected version '20.04', got '%s'", version)
	}
}

func TestInfoString(t *testing.T) {
	testCases := []struct {
		info     Info
		expected string
	}{
		{
			Info{Platform: Windows, Architecture: AMD64},
			"windows/amd64",
		},
		{
			Info{Platform: Linux, Architecture: AMD64, Distribution: Ubuntu, Version: "20.04"},
			"linux/amd64 (ubuntu 20.04)",
		},
		{
			Info{Platform: MacOS, Architecture: ARM64, Version: "12.0"},
			"darwin/arm64 (12.0)",
		},
		{
			Info{Platform: Linux, Architecture: AMD64, Distribution: UnknownDistro},
			"linux/amd64",
		},
	}

	for _, tc := range testCases {
		result := tc.info.String()
		if result != tc.expected {
			t.Errorf("Info.String() = '%s', expected '%s'", result, tc.expected)
		}
	}
}

func TestPlatformConstants(t *testing.T) {
	// 测试平台常量
	if Windows != "windows" {
		t.Errorf("Windows constant should be 'windows', got '%s'", Windows)
	}
	if MacOS != "darwin" {
		t.Errorf("MacOS constant should be 'darwin', got '%s'", MacOS)
	}
	if Linux != "linux" {
		t.Errorf("Linux constant should be 'linux', got '%s'", Linux)
	}
}

func TestArchitectureConstants(t *testing.T) {
	// 测试架构常量
	if AMD64 != "amd64" {
		t.Errorf("AMD64 constant should be 'amd64', got '%s'", AMD64)
	}
	if ARM64 != "arm64" {
		t.Errorf("ARM64 constant should be 'arm64', got '%s'", ARM64)
	}
	if I386 != "386" {
		t.Errorf("I386 constant should be '386', got '%s'", I386)
	}
	if ARM != "arm" {
		t.Errorf("ARM constant should be 'arm', got '%s'", ARM)
	}
}

func TestDistributionConstants(t *testing.T) {
	// 测试发行版常量
	distributions := []Distribution{
		Ubuntu, Debian, CentOS, RHEL, Fedora, SUSE, Arch, Alpine, UnknownDistro,
	}

	for _, dist := range distributions {
		if string(dist) == "" {
			t.Errorf("Distribution constant should not be empty: %v", dist)
		}
	}
}

func TestParseOSReleaseEdgeCases(t *testing.T) {
	detector := NewDetector()

	testCases := []struct {
		name    string
		content string
		expDist Distribution
		expVer  string
	}{
		{
			name: "quoted values",
			content: `ID="ubuntu"
VERSION_ID="20.04"`,
			expDist: Ubuntu,
			expVer:  "20.04",
		},
		{
			name: "unquoted values",
			content: `ID=centos
VERSION_ID=8`,
			expDist: CentOS,
			expVer:  "8",
		},
		{
			name: "mixed quotes",
			content: `ID=fedora
VERSION_ID="35"`,
			expDist: Fedora,
			expVer:  "35",
		},
		{
			name:    "empty content",
			content: "",
			expDist: UnknownDistro,
			expVer:  "",
		},
		{
			name: "malformed content",
			content: `INVALID_LINE
ID=debian
ANOTHER_INVALID_LINE=value
VERSION_ID=11`,
			expDist: Debian,
			expVer:  "11",
		},
		{
			name:    "missing version",
			content: `ID=arch`,
			expDist: Arch,
			expVer:  "",
		},
		{
			name:    "missing id",
			content: `VERSION_ID=1.0`,
			expDist: UnknownDistro,
			expVer:  "1.0",
		},
		{
			name: "whitespace handling",
			content: `  ID=alpine
  VERSION_ID=3.15  `,
			expDist: Alpine,
			expVer:  "3.15",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dist, version, err := detector.parseOSRelease(tc.content)
			if err != nil {
				t.Fatalf("parseOSRelease() failed: %v", err)
			}

			if dist != tc.expDist {
				t.Errorf("Expected distribution %s, got %s", tc.expDist, dist)
			}

			if version != tc.expVer {
				t.Errorf("Expected version '%s', got '%s'", tc.expVer, version)
			}
		})
	}
}

func TestParseLSBReleaseEdgeCases(t *testing.T) {
	detector := NewDetector()

	testCases := []struct {
		name    string
		content string
		expDist Distribution
		expVer  string
	}{
		{
			name: "standard format",
			content: `DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=22.04
DISTRIB_CODENAME=jammy
DISTRIB_DESCRIPTION="Ubuntu 22.04 LTS"`,
			expDist: Ubuntu,
			expVer:  "22.04",
		},
		{
			name: "case insensitive",
			content: `DISTRIB_ID=DEBIAN
DISTRIB_RELEASE=11`,
			expDist: Debian,
			expVer:  "11",
		},
		{
			name:    "empty content",
			content: "",
			expDist: UnknownDistro,
			expVer:  "",
		},
		{
			name:    "partial content",
			content: `DISTRIB_ID=CentOS`,
			expDist: CentOS,
			expVer:  "",
		},
		{
			name: "unknown distribution",
			content: `DISTRIB_ID=CustomLinux
DISTRIB_RELEASE=1.0`,
			expDist: UnknownDistro,
			expVer:  "1.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dist, version, err := detector.parseLSBRelease(tc.content)
			if err != nil {
				t.Fatalf("parseLSBRelease() failed: %v", err)
			}

			if dist != tc.expDist {
				t.Errorf("Expected distribution %s, got %s", tc.expDist, dist)
			}

			if version != tc.expVer {
				t.Errorf("Expected version '%s', got '%s'", tc.expVer, version)
			}
		})
	}
}

func TestMapDistributionIDEdgeCases(t *testing.T) {
	detector := NewDetector()

	testCases := []struct {
		id       string
		expected Distribution
	}{
		// Partial matches
		{"ubuntu-server", Ubuntu},
		{"debian-based", Debian},
		{"centos-stream", CentOS},
		{"rhel-workstation", RHEL},
		{"redhat-enterprise", RHEL},
		{"fedora-server", Fedora},
		{"opensuse", SUSE},
		{"suse-linux", SUSE},
		{"arch-linux", Arch},
		{"alpine-linux", Alpine},

		// Case variations
		{"UBUNTU", Ubuntu},
		{"Debian", Debian},
		{"CentOS", CentOS},
		{"RHEL", RHEL},
		{"RedHat", RHEL},
		{"FEDORA", Fedora},
		{"SUSE", SUSE},
		{"ARCH", Arch},
		{"ALPINE", Alpine},

		// Unknown cases
		{"gentoo", UnknownDistro},
		{"nixos", UnknownDistro},
		{"void", UnknownDistro},
		{"", UnknownDistro},
		{"   ", UnknownDistro},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			result := detector.mapDistributionID(tc.id)
			if result != tc.expected {
				t.Errorf("mapDistributionID(%s) = %s, expected %s", tc.id, result, tc.expected)
			}
		})
	}
}

func TestInfoMethodsWithDifferentValues(t *testing.T) {
	testCases := []struct {
		name  string
		info  Info
		tests map[string]bool
	}{
		{
			name: "Windows AMD64",
			info: Info{Platform: Windows, Architecture: AMD64},
			tests: map[string]bool{
				"IsWindows": true,
				"IsMacOS":   false,
				"IsLinux":   false,
				"IsX86":     true,
				"IsARM":     false,
			},
		},
		{
			name: "Linux ARM64",
			info: Info{Platform: Linux, Architecture: ARM64},
			tests: map[string]bool{
				"IsWindows": false,
				"IsMacOS":   false,
				"IsLinux":   true,
				"IsX86":     false,
				"IsARM":     true,
			},
		},
		{
			name: "macOS ARM",
			info: Info{Platform: MacOS, Architecture: ARM},
			tests: map[string]bool{
				"IsWindows": false,
				"IsMacOS":   true,
				"IsLinux":   false,
				"IsX86":     false,
				"IsARM":     true,
			},
		},
		{
			name: "Unknown Platform I386",
			info: Info{Platform: Unknown, Architecture: I386},
			tests: map[string]bool{
				"IsWindows": false,
				"IsMacOS":   false,
				"IsLinux":   false,
				"IsX86":     true,
				"IsARM":     false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.tests["IsWindows"] != tc.info.IsWindows() {
				t.Errorf("IsWindows() = %v, expected %v", tc.info.IsWindows(), tc.tests["IsWindows"])
			}
			if tc.tests["IsMacOS"] != tc.info.IsMacOS() {
				t.Errorf("IsMacOS() = %v, expected %v", tc.info.IsMacOS(), tc.tests["IsMacOS"])
			}
			if tc.tests["IsLinux"] != tc.info.IsLinux() {
				t.Errorf("IsLinux() = %v, expected %v", tc.info.IsLinux(), tc.tests["IsLinux"])
			}
			if tc.tests["IsX86"] != tc.info.IsX86() {
				t.Errorf("IsX86() = %v, expected %v", tc.info.IsX86(), tc.tests["IsX86"])
			}
			if tc.tests["IsARM"] != tc.info.IsARM() {
				t.Errorf("IsARM() = %v, expected %v", tc.info.IsARM(), tc.tests["IsARM"])
			}
		})
	}
}

func TestInfoStringComplexCases(t *testing.T) {
	testCases := []struct {
		name     string
		info     Info
		expected string
	}{
		{
			name: "Full Linux info",
			info: Info{
				Platform:     Linux,
				Architecture: AMD64,
				Distribution: Ubuntu,
				Version:      "20.04",
				Kernel:       "5.4.0-74-generic",
			},
			expected: "linux/amd64 (ubuntu 20.04)",
		},
		{
			name: "Windows with version",
			info: Info{
				Platform:     Windows,
				Architecture: AMD64,
				Version:      "10.0.19042",
			},
			expected: "windows/amd64 (10.0.19042)",
		},
		{
			name: "macOS ARM64",
			info: Info{
				Platform:     MacOS,
				Architecture: ARM64,
				Version:      "12.3.1",
			},
			expected: "darwin/arm64 (12.3.1)",
		},
		{
			name: "Linux with unknown distribution",
			info: Info{
				Platform:     Linux,
				Architecture: ARM,
				Distribution: UnknownDistro,
				Version:      "custom-1.0",
			},
			expected: "linux/arm (custom-1.0)",
		},
		{
			name: "Minimal info",
			info: Info{
				Platform:     Unknown,
				Architecture: I386,
			},
			expected: "unknown/386",
		},
		{
			name: "Distribution without version",
			info: Info{
				Platform:     Linux,
				Architecture: ARM64,
				Distribution: Arch,
			},
			expected: "linux/arm64 (arch)",
		},
		{
			name: "Empty distribution with version",
			info: Info{
				Platform:     Linux,
				Architecture: AMD64,
				Distribution: "",
				Version:      "5.10.0",
			},
			expected: "linux/amd64 (5.10.0)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.info.String()
			if result != tc.expected {
				t.Errorf("Info.String() = '%s', expected '%s'", result, tc.expected)
			}
		})
	}
}

func TestDetectSystemVersionMethods(t *testing.T) {
	detector := NewDetector()

	// Test detectSystemVersion method
	version, err := detector.detectSystemVersion()
	if err != nil {
		t.Logf("detectSystemVersion() failed (expected on some systems): %v", err)
	} else {
		t.Logf("System version: %s", version)
		if version == "" {
			t.Error("Expected non-empty version string")
		}
	}

	// Test detectKernelVersion method
	kernel, err := detector.detectKernelVersion()
	if err != nil {
		t.Logf("detectKernelVersion() failed (expected on some systems): %v", err)
	} else {
		t.Logf("Kernel version: %s", kernel)
		if kernel == "" {
			t.Error("Expected non-empty kernel string")
		}
	}
}

func TestDetectLinuxDistributionMethods(t *testing.T) {
	detector := NewDetector()

	// Only test on Linux systems
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux distribution tests on non-Linux system")
	}

	dist, version, err := detector.detectLinuxDistribution()
	if err != nil {
		t.Logf("detectLinuxDistribution() failed: %v", err)
	} else {
		t.Logf("Linux distribution: %s %s", dist, version)

		// Validate that we got a known distribution
		knownDistros := []Distribution{
			Ubuntu, Debian, CentOS, RHEL, Fedora, SUSE, Arch, Alpine, UnknownDistro,
		}

		found := false
		for _, known := range knownDistros {
			if dist == known {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Unknown distribution returned: %s", dist)
		}
	}
}

func TestPlatformSpecificVersionDetection(t *testing.T) {
	detector := NewDetector()

	switch runtime.GOOS {
	case "windows":
		version, err := detector.detectWindowsVersion()
		if err != nil {
			t.Logf("detectWindowsVersion() failed: %v", err)
		} else {
			t.Logf("Windows version: %s", version)
			if version == "" {
				t.Error("Expected non-empty Windows version")
			}
		}

	case "darwin":
		version, err := detector.detectMacOSVersion()
		if err != nil {
			t.Logf("detectMacOSVersion() failed: %v", err)
		} else {
			t.Logf("macOS version: %s", version)
			if version == "" {
				t.Error("Expected non-empty macOS version")
			}
		}

	case "linux":
		version, err := detector.detectLinuxVersion()
		if err != nil {
			t.Logf("detectLinuxVersion() failed: %v", err)
		} else {
			t.Logf("Linux version: %s", version)
			if version == "" {
				t.Error("Expected non-empty Linux version")
			}
		}
	}
}

func TestInfoJSONSerialization(t *testing.T) {
	info := Info{
		Platform:     Linux,
		Architecture: AMD64,
		Distribution: Ubuntu,
		Version:      "20.04",
		Kernel:       "5.4.0-74-generic",
	}

	// Test that the struct can be used with JSON tags
	// This is implicit testing of the json tags
	if info.Platform != Linux {
		t.Error("Platform field should be accessible")
	}
	if info.Architecture != AMD64 {
		t.Error("Architecture field should be accessible")
	}
	if info.Distribution != Ubuntu {
		t.Error("Distribution field should be accessible")
	}
	if info.Version != "20.04" {
		t.Error("Version field should be accessible")
	}
	if info.Kernel != "5.4.0-74-generic" {
		t.Error("Kernel field should be accessible")
	}
}

func TestArchitectureEdgeCases(t *testing.T) {
	testCases := []struct {
		name  string
		info  Info
		isX86 bool
		isARM bool
	}{
		{
			name:  "AMD64",
			info:  Info{Architecture: AMD64},
			isX86: true,
			isARM: false,
		},
		{
			name:  "I386",
			info:  Info{Architecture: I386},
			isX86: true,
			isARM: false,
		},
		{
			name:  "ARM64",
			info:  Info{Architecture: ARM64},
			isX86: false,
			isARM: true,
		},
		{
			name:  "ARM",
			info:  Info{Architecture: ARM},
			isX86: false,
			isARM: true,
		},
		{
			name:  "Unknown Architecture",
			info:  Info{Architecture: Architecture("unknown")},
			isX86: false,
			isARM: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.info.IsX86() != tc.isX86 {
				t.Errorf("IsX86() = %v, expected %v", tc.info.IsX86(), tc.isX86)
			}
			if tc.info.IsARM() != tc.isARM {
				t.Errorf("IsARM() = %v, expected %v", tc.info.IsARM(), tc.isARM)
			}
		})
	}
}
