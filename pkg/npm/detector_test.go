package npm

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}

	if detector.timeout != 10*time.Second {
		t.Errorf("Expected default timeout 10s, got %v", detector.timeout)
	}
}

func TestDetectorSetTimeout(t *testing.T) {
	detector := NewDetector()
	timeout := 30 * time.Second
	detector.SetTimeout(timeout)

	if detector.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, detector.timeout)
	}
}

func TestDetectorIsValidVersion(t *testing.T) {
	detector := NewDetector()

	testCases := []struct {
		version string
		valid   bool
	}{
		{"8.19.2", true},
		{"6.14.15", true},
		{"10.0.0", true},
		{"1.2.3", true},
		{"", false},
		{"invalid", false},
		{"1.2", false},     // 不符合x.y.z格式
		{"abc.def", false}, // 不符合数字格式
		{"1.2.3.4", true},  // 符合开头的x.y.z格式
		{"0.0.1", true},
	}

	for _, tc := range testCases {
		result := detector.isValidVersion(tc.version)
		if result != tc.valid {
			t.Errorf("isValidVersion(%s) = %v, expected %v", tc.version, result, tc.valid)
		}
	}
}

func TestDetectorGetCommonNpmPaths(t *testing.T) {
	detector := NewDetector()
	paths := detector.getCommonNpmPaths()

	if len(paths) == 0 {
		t.Error("Expected non-empty paths list")
	}

	// 检查是否包含通用的npm路径
	foundGeneric := false
	for _, path := range paths {
		if path == "npm" {
			foundGeneric = true
			break
		}
	}

	if !foundGeneric {
		t.Error("Expected to find generic 'npm' in paths")
	}
}

func TestDetectorFindNpmPath(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 这个测试可能会失败如果系统没有npm，但我们测试方法是否能正常调用
	path, err := detector.findNpmPath(ctx)

	// 如果找到了npm，路径不应该为空
	if err == nil && path == "" {
		t.Error("Found npm but path is empty")
	}

	// 如果没找到npm，应该返回ErrNpmNotFound
	if err != nil && !IsNpmNotFound(err) {
		t.Logf("npm not found (expected in some environments): %v", err)
	}

	t.Logf("npm path result: path=%s, err=%v", path, err)
}

func TestDetectorGetNpmVersion(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试有效的npm路径（如果npm可用）
	if path, err := detector.findNpmPath(ctx); err == nil {
		version, err := detector.getNpmVersion(ctx, path)
		if err != nil {
			t.Logf("Failed to get npm version (expected if npm not available): %v", err)
		} else {
			if version == "" {
				t.Error("Got empty version string")
			}
			if !detector.isValidVersion(version) {
				t.Errorf("Got invalid version format: %s", version)
			}
			t.Logf("npm version: %s", version)
		}
	}

	// 测试无效的npm路径
	_, err := detector.getNpmVersion(ctx, "definitely-nonexistent-npm-command")
	if err == nil {
		t.Error("Expected error for nonexistent npm command")
	}
}

func TestDetectorDetect(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := detector.Detect(ctx)

	// 如果检测成功
	if err == nil {
		if info == nil {
			t.Fatal("Detect() returned nil info without error")
		}

		if info.Path == "" {
			t.Error("Expected non-empty npm path")
		}

		if info.Version == "" {
			t.Error("Expected non-empty npm version")
		}

		if !info.Available {
			t.Error("Expected Available to be true when detection succeeds")
		}

		t.Logf("npm detected: path=%s, version=%s, node_path=%s, node_version=%s",
			info.Path, info.Version, info.NodePath, info.NodeVersion)
	} else {
		// 如果检测失败，应该返回适当的错误
		t.Logf("npm detection failed (expected in some environments): %v", err)

		if info != nil && info.Available {
			t.Error("Expected Available to be false when detection fails")
		}
	}
}

func TestDetectorIsAvailable(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	available := detector.IsAvailable(ctx)
	t.Logf("npm available: %v", available)

	// 验证IsAvailable与Detect的一致性
	info, err := detector.Detect(ctx)
	if err == nil && info != nil {
		if available != info.Available {
			t.Errorf("IsAvailable() = %v, but Detect().Available = %v", available, info.Available)
		}
	}
}

func TestDetectorGetVersion(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	version, err := detector.GetVersion(ctx)

	if err == nil {
		if version == "" {
			t.Error("Got empty version string")
		}
		if !detector.isValidVersion(version) {
			t.Errorf("Got invalid version format: %s", version)
		}
		t.Logf("npm version: %s", version)
	} else {
		t.Logf("Failed to get npm version (expected if npm not available): %v", err)
	}
}

func TestDetectorGetNpmConfig(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试获取registry配置
	registry, err := detector.GetNpmConfig(ctx, "registry")
	if err == nil {
		if registry == "" {
			t.Error("Expected non-empty registry")
		}
		t.Logf("npm registry: %s", registry)
	} else {
		t.Logf("Failed to get npm config (expected if npm not available): %v", err)
	}

	// 测试获取不存在的配置
	_, err = detector.GetNpmConfig(ctx, "nonexistent-config-key")
	// 这可能成功（返回空值）或失败，取决于npm的行为
	t.Logf("nonexistent config result: %v", err)
}

func TestDetectorGetNpmConfigList(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := detector.GetNpmConfigList(ctx)
	if err == nil {
		if len(config) == 0 {
			t.Error("Expected non-empty config map")
		}

		// 检查是否包含常见的配置项
		if registry, exists := config["registry"]; exists {
			if registry == "" {
				t.Error("Expected non-empty registry value")
			}
			t.Logf("registry from config list: %s", registry)
		}

		t.Logf("npm config count: %d", len(config))
	} else {
		t.Logf("Failed to get npm config list (expected if npm not available): %v", err)
	}
}

func TestDetectorGetGlobalPackagesPath(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	path, err := detector.GetGlobalPackagesPath(ctx)
	if err == nil {
		if path == "" {
			t.Error("Expected non-empty global packages path")
		}
		t.Logf("global packages path: %s", path)
	} else {
		t.Logf("Failed to get global packages path (expected if npm not available): %v", err)
	}
}

func TestDetectorGetCacheDir(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cacheDir, err := detector.GetCacheDir(ctx)
	if err == nil {
		if cacheDir == "" {
			t.Error("Expected non-empty cache directory")
		}
		t.Logf("npm cache dir: %s", cacheDir)
	} else {
		t.Logf("Failed to get cache dir (expected if npm not available): %v", err)
	}
}

func TestDetectorGetRegistry(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registry, err := detector.GetRegistry(ctx)
	if err == nil {
		if registry == "" {
			t.Error("Expected non-empty registry")
		}
		// 检查是否是有效的URL格式
		if !strings.Contains(registry, "://") {
			t.Errorf("Registry doesn't look like a URL: %s", registry)
		}
		t.Logf("npm registry: %s", registry)
	} else {
		t.Logf("Failed to get registry (expected if npm not available): %v", err)
	}
}

func TestDetectorGetNodeInfo(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nodePath, nodeVersion, err := detector.getNodeInfo(ctx)
	if err == nil {
		t.Logf("Node.js detected: path=%s, version=%s", nodePath, nodeVersion)

		if nodeVersion != "" && !detector.isValidVersion(nodeVersion) {
			t.Errorf("Got invalid Node.js version format: %s", nodeVersion)
		}
	} else {
		t.Logf("Failed to get Node.js info (expected if Node.js not available): %v", err)
	}
}

func TestDetectorExecCommand(t *testing.T) {
	detector := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试简单命令
	output, err := detector.execCommand(ctx, "echo", "test")
	if err == nil {
		if !strings.Contains(output, "test") {
			t.Errorf("Expected output to contain 'test', got: %s", output)
		}
	} else {
		// 在某些环境中echo可能不可用
		t.Logf("echo command failed: %v", err)
	}

	// 测试不存在的命令
	_, err = detector.execCommand(ctx, "definitely-nonexistent-command")
	if err == nil {
		t.Error("Expected error for nonexistent command")
	}
}
