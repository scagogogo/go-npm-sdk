package npm

import (
	"errors"
	"testing"
)

func TestNpmError(t *testing.T) {
	// 测试基本的NpmError创建
	err := NewNpmError("install", "lodash", 1, "stdout content", "stderr content", errors.New("underlying error"))

	if err.Op != "install" {
		t.Errorf("Expected Op 'install', got '%s'", err.Op)
	}

	if err.Package != "lodash" {
		t.Errorf("Expected Package 'lodash', got '%s'", err.Package)
	}

	if err.ExitCode != 1 {
		t.Errorf("Expected ExitCode 1, got %d", err.ExitCode)
	}

	if err.Stdout != "stdout content" {
		t.Errorf("Expected Stdout 'stdout content', got '%s'", err.Stdout)
	}

	if err.Stderr != "stderr content" {
		t.Errorf("Expected Stderr 'stderr content', got '%s'", err.Stderr)
	}

	if err.Err == nil {
		t.Error("Expected underlying error to be set")
	}
}

func TestNpmErrorError(t *testing.T) {
	testCases := []struct {
		name     string
		err      *NpmError
		expected string
	}{
		{
			name: "with package",
			err: &NpmError{
				Op:       "install",
				Package:  "lodash",
				ExitCode: 1,
				Err:      errors.New("command failed"),
			},
			expected: "npm install failed for package 'lodash': command failed",
		},
		{
			name: "without package",
			err: &NpmError{
				Op:       "version",
				ExitCode: 127,
				Err:      errors.New("command not found"),
			},
			expected: "npm version failed: command not found",
		},
		{
			name: "with zero exit code",
			err: &NpmError{
				Op:  "list",
				Err: errors.New("parsing error"),
			},
			expected: "npm list failed: parsing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.Error()
			if result != tc.expected {
				t.Errorf("Expected error message '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestNpmErrorUnwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	npmErr := &NpmError{
		Op:  "test",
		Err: underlyingErr,
	}

	unwrapped := npmErr.Unwrap()
	if unwrapped != underlyingErr {
		t.Errorf("Expected unwrapped error to be the underlying error")
	}

	// 测试没有underlying error的情况
	npmErrNoUnderlying := &NpmError{Op: "test"}
	if npmErrNoUnderlying.Unwrap() != nil {
		t.Error("Expected Unwrap() to return nil when no underlying error")
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("package", "invalid-name", "package name is invalid")

	if err.Field != "package" {
		t.Errorf("Expected Field 'package', got '%s'", err.Field)
	}

	if err.Value != "invalid-name" {
		t.Errorf("Expected Value 'invalid-name', got '%s'", err.Value)
	}

	if err.Reason != "package name is invalid" {
		t.Errorf("Expected Reason 'package name is invalid', got '%s'", err.Reason)
	}
}

func TestValidationErrorError(t *testing.T) {
	testCases := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "with value",
			err: &ValidationError{
				Field:  "version",
				Value:  "invalid",
				Reason: "invalid version format",
			},
			expected: "validation failed for field 'version' with value 'invalid': invalid version format",
		},
		{
			name: "without value",
			err: &ValidationError{
				Field:  "name",
				Value:  "",
				Reason: "name is required",
			},
			expected: "validation failed for field 'name' with value '': name is required",
		},
		{
			name: "empty reason",
			err: &ValidationError{
				Field: "test",
				Value: "value",
			},
			expected: "validation failed for field 'test' with value 'value': ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.Error()
			if result != tc.expected {
				t.Errorf("Expected error message '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestInstallError(t *testing.T) {
	underlyingErr := NewNpmError("install", "lodash", 1, "", "permission denied", errors.New("EACCES"))
	err := NewInstallError("lodash", "permission denied", underlyingErr)

	if err.Package != "lodash" {
		t.Errorf("Expected Package 'lodash', got '%s'", err.Package)
	}

	if err.Reason != "permission denied" {
		t.Errorf("Expected Reason 'permission denied', got '%s'", err.Reason)
	}

	if err.Err != underlyingErr {
		t.Error("Expected underlying error to be set correctly")
	}
}

func TestInstallErrorError(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &InstallError{
		Package: "react",
		Reason:  "network timeout",
		Err:     underlyingErr,
	}

	expected := "failed to install package 'react': network timeout"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, result)
	}
}

func TestInstallErrorUnwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	installErr := &InstallError{
		Package: "test",
		Reason:  "test reason",
		Err:     underlyingErr,
	}

	unwrapped := installErr.Unwrap()
	if unwrapped != underlyingErr {
		t.Error("Expected unwrapped error to be the underlying error")
	}
}

func TestUninstallError(t *testing.T) {
	underlyingErr := NewNpmError("uninstall", "lodash", 1, "", "package not found", errors.New("ENOENT"))
	err := NewUninstallError("lodash", "package not found", underlyingErr)

	if err.Package != "lodash" {
		t.Errorf("Expected Package 'lodash', got '%s'", err.Package)
	}

	if err.Reason != "package not found" {
		t.Errorf("Expected Reason 'package not found', got '%s'", err.Reason)
	}

	if err.Err != underlyingErr {
		t.Error("Expected underlying error to be set correctly")
	}
}

func TestUninstallErrorError(t *testing.T) {
	err := &UninstallError{
		Package: "vue",
		Reason:  "dependency conflict",
	}

	expected := "failed to uninstall package 'vue': dependency conflict"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, result)
	}
}

func TestPlatformError(t *testing.T) {
	underlyingErr := errors.New("unsupported architecture")
	err := NewPlatformError("linux/arm", "platform not supported", underlyingErr)

	if err.Platform != "linux/arm" {
		t.Errorf("Expected Platform 'linux/arm', got '%s'", err.Platform)
	}

	if err.Reason != "platform not supported" {
		t.Errorf("Expected Reason 'platform not supported', got '%s'", err.Reason)
	}

	if err.Err != underlyingErr {
		t.Error("Expected underlying error to be set correctly")
	}
}

func TestPlatformErrorError(t *testing.T) {
	testCases := []struct {
		name     string
		err      *PlatformError
		expected string
	}{
		{
			name: "with underlying error",
			err: &PlatformError{
				Platform: "windows/arm64",
				Reason:   "not supported",
				Err:      errors.New("arch error"),
			},
			expected: "platform error on windows/arm64: not supported",
		},
		{
			name: "without underlying error",
			err: &PlatformError{
				Platform: "darwin/amd64",
				Reason:   "configuration issue",
			},
			expected: "platform error on darwin/amd64: configuration issue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.Error()
			if result != tc.expected {
				t.Errorf("Expected error message '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestDownloadError(t *testing.T) {
	underlyingErr := errors.New("connection timeout")
	err := NewDownloadError("https://registry.npmjs.org", "download failed", underlyingErr)

	if err.URL != "https://registry.npmjs.org" {
		t.Errorf("Expected URL 'https://registry.npmjs.org', got '%s'", err.URL)
	}

	if err.Reason != "download failed" {
		t.Errorf("Expected Reason 'download failed', got '%s'", err.Reason)
	}

	if err.Err != underlyingErr {
		t.Error("Expected underlying error to be set correctly")
	}
}

func TestDownloadErrorError(t *testing.T) {
	err := &DownloadError{
		URL:    "https://example.com/package.tgz",
		Reason: "connection refused",
		Err:    errors.New("dial tcp: connection refused"),
	}

	expected := "failed to download from https://example.com/package.tgz: connection refused"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, result)
	}
}

func TestErrorTypeChecking(t *testing.T) {
	// 测试IsNpmNotFound
	notFoundErr := &NpmError{Op: "version", ExitCode: 127, Err: errors.New("command not found")}
	if !IsNpmNotFound(notFoundErr) {
		t.Error("Expected IsNpmNotFound to return true for command not found error")
	}

	otherErr := &NpmError{Op: "install", ExitCode: 1, Err: errors.New("other error")}
	if IsNpmNotFound(otherErr) {
		t.Error("Expected IsNpmNotFound to return false for other errors")
	}

	// 测试IsPackageNotFound
	packageNotFoundErr := &NpmError{Op: "view", ExitCode: 1, Stderr: "404 Not Found"}
	if !IsPackageNotFound(packageNotFoundErr) {
		t.Error("Expected IsPackageNotFound to return true for 404 error")
	}

	// 测试非相关错误
	genericErr := errors.New("generic error")
	if IsNpmNotFound(genericErr) {
		t.Error("Expected IsNpmNotFound to return false for generic error")
	}
	if IsPackageNotFound(genericErr) {
		t.Error("Expected IsPackageNotFound to return false for generic error")
	}
}

func TestErrorConstants(t *testing.T) {
	// 测试预定义错误常量
	if ErrNpmNotFound == nil {
		t.Error("Expected ErrNpmNotFound to be defined")
	}

	if ErrPackageNotFound == nil {
		t.Error("Expected ErrPackageNotFound to be defined")
	}

	if ErrInvalidVersion == nil {
		t.Error("Expected ErrInvalidVersion to be defined")
	}

	if ErrInvalidPackageName == nil {
		t.Error("Expected ErrInvalidPackageName to be defined")
	}

	if ErrPermissionDenied == nil {
		t.Error("Expected ErrPermissionDenied to be defined")
	}

	if ErrUnsupportedPlatform == nil {
		t.Error("Expected ErrUnsupportedPlatform to be defined")
	}
}

func TestErrorWrapping(t *testing.T) {
	// 测试错误包装和解包
	originalErr := errors.New("original error")

	// NpmError包装
	npmErr := NewNpmError("test", "pkg", 1, "", "", originalErr)
	if !errors.Is(npmErr, originalErr) {
		t.Error("Expected errors.Is to find original error in NpmError")
	}

	// InstallError包装
	installErr := NewInstallError("pkg", "reason", npmErr)
	if !errors.Is(installErr, originalErr) {
		t.Error("Expected errors.Is to find original error through InstallError")
	}

	// 测试错误类型断言
	var targetNpmErr *NpmError
	if !errors.As(installErr, &targetNpmErr) {
		t.Error("Expected errors.As to find NpmError in InstallError")
	}

	if targetNpmErr.Op != "test" {
		t.Errorf("Expected Op 'test', got '%s'", targetNpmErr.Op)
	}
}
