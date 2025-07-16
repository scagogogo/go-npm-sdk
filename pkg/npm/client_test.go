package npm

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClientWithPath(t *testing.T) {
	npmPath := "npm"
	client, err := NewClientWithPath(npmPath)
	if err != nil {
		t.Fatalf("NewClientWithPath() failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewClientWithPath() returned nil client")
	}
}

func TestClientIsAvailable(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 这个测试可能会失败如果系统没有安装npm
	// 但我们仍然测试方法是否能正常调用
	available := client.IsAvailable(ctx)
	t.Logf("npm available: %v", available)
}

func TestClientVersion(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 只有在npm可用时才测试版本获取
	if !client.IsAvailable(ctx) {
		t.Skip("npm not available, skipping version test")
	}

	version, err := client.Version(ctx)
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}

	if version == "" {
		t.Fatal("Version() returned empty version")
	}

	t.Logf("npm version: %s", version)
}

func TestValidationErrors(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试空包名验证
	err = client.InstallPackage(ctx, "", InstallOptions{})
	if err == nil {
		t.Fatal("InstallPackage() should fail with empty package name")
	}

	var validationErr *ValidationError
	if !IsValidationError(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got: %T", err)
	}

	// 测试空脚本名验证
	err = client.RunScript(ctx, "")
	if err == nil {
		t.Fatal("RunScript() should fail with empty script name")
	}

	if !IsValidationError(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got: %T", err)
	}
}

func TestInitOptions(t *testing.T) {
	options := InitOptions{
		Name:        "test-project",
		Version:     "1.0.0",
		Description: "Test project",
		Author:      "Test Author",
		License:     "MIT",
		Private:     true,
		Force:       true,
	}

	if options.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", options.Name)
	}

	if !options.Private {
		t.Error("Expected private to be true")
	}

	if !options.Force {
		t.Error("Expected force to be true")
	}
}

func TestInstallOptions(t *testing.T) {
	options := InstallOptions{
		SaveDev:       true,
		SaveOptional:  false,
		SaveExact:     true,
		Global:        false,
		Production:    false,
		Registry:      "https://registry.npmjs.org/",
		Force:         true,
		IgnoreScripts: true,
	}

	if !options.SaveDev {
		t.Error("Expected SaveDev to be true")
	}

	if options.SaveOptional {
		t.Error("Expected SaveOptional to be false")
	}

	if !options.SaveExact {
		t.Error("Expected SaveExact to be true")
	}

	if options.Registry != "https://registry.npmjs.org/" {
		t.Errorf("Expected registry 'https://registry.npmjs.org/', got '%s'", options.Registry)
	}
}

func TestUninstallOptions(t *testing.T) {
	options := UninstallOptions{
		SaveDev:    true,
		Global:     false,
		WorkingDir: "/tmp/test",
	}

	if !options.SaveDev {
		t.Error("Expected SaveDev to be true")
	}

	if options.Global {
		t.Error("Expected Global to be false")
	}

	if options.WorkingDir != "/tmp/test" {
		t.Errorf("Expected WorkingDir '/tmp/test', got '%s'", options.WorkingDir)
	}
}

func TestListOptions(t *testing.T) {
	options := ListOptions{
		Global:     true,
		Depth:      2,
		Production: true,
		JSON:       true,
		WorkingDir: "/tmp/test",
	}

	if !options.Global {
		t.Error("Expected Global to be true")
	}

	if options.Depth != 2 {
		t.Errorf("Expected Depth 2, got %d", options.Depth)
	}

	if !options.Production {
		t.Error("Expected Production to be true")
	}

	if !options.JSON {
		t.Error("Expected JSON to be true")
	}
}

func TestPublishOptions(t *testing.T) {
	options := PublishOptions{
		Tag:        "beta",
		Access:     "public",
		Registry:   "https://registry.npmjs.org/",
		WorkingDir: "/tmp/test",
		DryRun:     true,
	}

	if options.Tag != "beta" {
		t.Errorf("Expected Tag 'beta', got '%s'", options.Tag)
	}

	if options.Access != "public" {
		t.Errorf("Expected Access 'public', got '%s'", options.Access)
	}

	if !options.DryRun {
		t.Error("Expected DryRun to be true")
	}
}

func TestClientInit(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test init with basic options
	options := InitOptions{
		Name:        "test-project",
		Version:     "1.0.0",
		Description: "Test project",
		Author:      "Test Author",
		License:     "MIT",
		Force:       true,
	}

	// This will fail if npm is not available, but we test the validation
	err = client.Init(ctx, options)
	// We don't assert success here since npm might not be available
	// but we test that the method can be called without panic
	t.Logf("Init result: %v", err)
}

func TestClientInstallPackage(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test install with empty package name (should fail validation)
	err = client.InstallPackage(ctx, "", InstallOptions{})
	if err == nil {
		t.Error("Expected error for empty package name")
	}

	var validationErr *ValidationError
	if !IsValidationError(err, &validationErr) {
		t.Errorf("Expected ValidationError, got: %T", err)
	}

	// Test install with valid package name
	options := InstallOptions{
		SaveDev:    true,
		SaveExact:  true,
		WorkingDir: "/tmp",
	}

	err = client.InstallPackage(ctx, "lodash", options)
	// We don't assert success here since npm might not be available
	t.Logf("InstallPackage result: %v", err)
}

func TestClientUninstallPackage(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test uninstall with empty package name (should fail validation)
	err = client.UninstallPackage(ctx, "", UninstallOptions{})
	if err == nil {
		t.Error("Expected error for empty package name")
	}

	// Test uninstall with valid package name
	options := UninstallOptions{
		SaveDev:    true,
		WorkingDir: "/tmp",
	}

	err = client.UninstallPackage(ctx, "lodash", options)
	// We don't assert success here since npm might not be available
	t.Logf("UninstallPackage result: %v", err)
}

func TestClientUpdatePackage(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test update with empty package name (should fail validation)
	err = client.UpdatePackage(ctx, "")
	if err == nil {
		t.Error("Expected error for empty package name")
	}

	// Test update with valid package name
	err = client.UpdatePackage(ctx, "lodash")
	// We don't assert success here since npm might not be available
	t.Logf("UpdatePackage result: %v", err)
}

func TestClientRunScript(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test run script with empty script name (should fail validation)
	err = client.RunScript(ctx, "")
	if err == nil {
		t.Error("Expected error for empty script name")
	}

	// Test run script with valid script name
	err = client.RunScript(ctx, "test", "--verbose")
	// We don't assert success here since npm might not be available
	t.Logf("RunScript result: %v", err)
}

func TestClientGetPackageInfo(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test get package info with empty package name (should fail validation)
	_, err = client.GetPackageInfo(ctx, "")
	if err == nil {
		t.Error("Expected error for empty package name")
	}

	// Test get package info with valid package name
	_, err = client.GetPackageInfo(ctx, "lodash")
	// We don't assert success here since npm might not be available
	t.Logf("GetPackageInfo result: %v", err)
}

func TestClientSearch(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// Test search with empty query (should fail validation)
	_, err = client.Search(ctx, "")
	if err == nil {
		t.Error("Expected error for empty search query")
	}

	// Test search with valid query
	_, err = client.Search(ctx, "react")
	// We don't assert success here since npm might not be available
	t.Logf("Search result: %v", err)
}

func TestClientPublish(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试基本发布选项
	options := PublishOptions{
		Tag:        "latest",
		Access:     "public",
		Registry:   "https://registry.npmjs.org/",
		WorkingDir: "/tmp/test-project",
		DryRun:     true, // 使用dry-run避免实际发布
	}

	err = client.Publish(ctx, options)
	// 我们不期望成功，因为没有真实的项目
	t.Logf("Publish result: %v", err)

	// 测试空选项
	emptyOptions := PublishOptions{}
	err = client.Publish(ctx, emptyOptions)
	t.Logf("Publish with empty options result: %v", err)
}

func TestClientListPackagesWithOptions(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试不同的列表选项
	testCases := []struct {
		name    string
		options ListOptions
	}{
		{
			name: "global packages",
			options: ListOptions{
				Global: true,
				Depth:  0,
			},
		},
		{
			name: "with depth",
			options: ListOptions{
				Depth: 2,
				JSON:  true,
			},
		},
		{
			name: "production only",
			options: ListOptions{
				Production: true,
				JSON:       false,
			},
		},
		{
			name: "all options",
			options: ListOptions{
				Global:     false,
				Depth:      1,
				Production: false,
				JSON:       true,
				WorkingDir: "/tmp",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packages, err := client.ListPackages(ctx, tc.options)
			// 结果取决于系统状态，我们主要测试不会panic
			t.Logf("ListPackages %s result: %d packages, error: %v", tc.name, len(packages), err)
		})
	}
}

func TestClientWithTimeout(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// 测试超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// 这些操作应该因为超时而失败
	_, err = client.Version(ctx)
	if err == nil {
		t.Log("Version call completed before timeout (possible)")
	} else {
		t.Logf("Version call timed out as expected: %v", err)
	}
}

func TestClientCancelledContext(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// 测试已取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	err = client.InstallPackage(ctx, "lodash", InstallOptions{})
	if err == nil {
		t.Log("InstallPackage completed despite cancelled context (possible)")
	} else {
		t.Logf("InstallPackage failed with cancelled context: %v", err)
	}
}

func TestClientEdgeCases(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试特殊字符的包名
	specialPackages := []string{
		"@scope/package",
		"package-with-dashes",
		"package_with_underscores",
		"package.with.dots",
		"123numeric-start",
	}

	for _, pkg := range specialPackages {
		t.Run("special_package_"+pkg, func(t *testing.T) {
			err := client.InstallPackage(ctx, pkg, InstallOptions{})
			// 我们不期望成功，主要测试不会panic
			t.Logf("Install %s result: %v", pkg, err)
		})
	}

	// 测试非常长的包名
	longPackageName := strings.Repeat("a", 1000)
	err = client.InstallPackage(ctx, longPackageName, InstallOptions{})
	t.Logf("Install long package name result: %v", err)

	// 测试空白字符的包名
	whitespacePackages := []string{
		" ",
		"\t",
		"\n",
		"  package  ",
	}

	for _, pkg := range whitespacePackages {
		err := client.InstallPackage(ctx, pkg, InstallOptions{})
		if err == nil {
			t.Errorf("Expected error for whitespace package name '%s'", pkg)
		}
	}
}

func TestClientInstallOptionsValidation(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试冲突的选项
	conflictingOptions := InstallOptions{
		SaveDev:      true,
		Production:   true, // 这两个选项冲突
		SaveOptional: true,
		SaveExact:    true,
		Global:       true,
	}

	err = client.InstallPackage(ctx, "test-package", conflictingOptions)
	t.Logf("Install with conflicting options result: %v", err)

	// 测试无效的registry URL
	invalidRegistryOptions := InstallOptions{
		Registry: "not-a-valid-url",
	}

	err = client.InstallPackage(ctx, "test-package", invalidRegistryOptions)
	t.Logf("Install with invalid registry result: %v", err)
}

func TestClientInitOptionsValidation(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试无效的包名
	invalidNameOptions := InitOptions{
		Name:    "Invalid Package Name", // 包含空格
		Version: "1.0.0",
		Force:   true,
	}

	err = client.Init(ctx, invalidNameOptions)
	t.Logf("Init with invalid name result: %v", err)

	// 测试无效的版本
	invalidVersionOptions := InitOptions{
		Name:    "valid-name",
		Version: "not-a-version",
		Force:   true,
	}

	err = client.Init(ctx, invalidVersionOptions)
	t.Logf("Init with invalid version result: %v", err)

	// 测试无效的许可证
	invalidLicenseOptions := InitOptions{
		Name:    "valid-name",
		Version: "1.0.0",
		License: "INVALID-LICENSE",
		Force:   true,
	}

	err = client.Init(ctx, invalidLicenseOptions)
	t.Logf("Init with invalid license result: %v", err)
}

func TestClientRunScriptWithArgs(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试带参数的脚本运行
	testCases := []struct {
		script string
		args   []string
	}{
		{"test", []string{"--verbose"}},
		{"build", []string{"--production", "--minify"}},
		{"start", []string{"--port", "3000"}},
		{"lint", []string{"--fix", "--ext", ".js,.ts"}},
	}

	for _, tc := range testCases {
		t.Run("script_"+tc.script, func(t *testing.T) {
			err := client.RunScript(ctx, tc.script, tc.args...)
			t.Logf("RunScript %s with args %v result: %v", tc.script, tc.args, err)
		})
	}

	// 测试特殊字符的脚本名
	specialScripts := []string{
		"pre:build",
		"post:install",
		"test:unit",
		"build:prod",
	}

	for _, script := range specialScripts {
		err := client.RunScript(ctx, script)
		t.Logf("RunScript %s result: %v", script, err)
	}
}

func TestClientSearchWithComplexQueries(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试复杂的搜索查询
	complexQueries := []string{
		"react hooks",
		"@types/node",
		"webpack plugin",
		"testing framework",
		"database orm",
		"ui component library",
	}

	for _, query := range complexQueries {
		t.Run("search_"+strings.ReplaceAll(query, " ", "_"), func(t *testing.T) {
			results, err := client.Search(ctx, query)
			t.Logf("Search '%s' result: %d results, error: %v", query, len(results), err)
		})
	}

	// 测试特殊字符的查询
	specialQueries := []string{
		"@scope/package",
		"package-name",
		"package_name",
		"package.name",
	}

	for _, query := range specialQueries {
		_, err := client.Search(ctx, query)
		t.Logf("Search special query '%s' result: %v", query, err)
	}
}

func TestClientGetPackageInfoEdgeCases(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()

	// 测试不存在的包
	nonexistentPackages := []string{
		"definitely-does-not-exist-package-12345",
		"@nonexistent/scope",
		"package-with-very-long-name-that-probably-does-not-exist",
	}

	for _, pkg := range nonexistentPackages {
		t.Run("nonexistent_"+pkg, func(t *testing.T) {
			info, err := client.GetPackageInfo(ctx, pkg)
			if err == nil {
				t.Logf("Unexpectedly found package %s: %+v", pkg, info)
			} else {
				t.Logf("Expected error for nonexistent package %s: %v", pkg, err)
			}
		})
	}

	// 测试特殊格式的包名
	specialPackages := []string{
		"@types/node",
		"@babel/core",
		"@angular/core",
	}

	for _, pkg := range specialPackages {
		t.Run("special_"+pkg, func(t *testing.T) {
			info, err := client.GetPackageInfo(ctx, pkg)
			t.Logf("GetPackageInfo %s result: info=%v, error=%v", pkg, info != nil, err)
		})
	}
}

// Helper function to check if error is ValidationError
func IsValidationError(err error, target **ValidationError) bool {
	if validationErr, ok := err.(*ValidationError); ok {
		if target != nil {
			*target = validationErr
		}
		return true
	}
	return false
}
