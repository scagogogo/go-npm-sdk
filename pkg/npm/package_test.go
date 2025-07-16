package npm

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewPackageJSON(t *testing.T) {
	filePath := "/tmp/test-package.json"
	pkg := NewPackageJSON(filePath)

	if pkg == nil {
		t.Fatal("NewPackageJSON() returned nil")
	}

	if pkg.filePath != filePath {
		t.Errorf("Expected filePath %s, got %s", filePath, pkg.filePath)
	}
}

func TestPackageJSONBasicFields(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试名称
	pkg.SetName("test-package")
	if pkg.GetName() != "test-package" {
		t.Errorf("Expected name 'test-package', got '%s'", pkg.GetName())
	}

	// 测试版本
	pkg.SetVersion("1.0.0")
	if pkg.GetVersion() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", pkg.GetVersion())
	}

	// 测试描述
	pkg.SetDescription("Test package")
	if pkg.GetDescription() != "Test package" {
		t.Errorf("Expected description 'Test package', got '%s'", pkg.GetDescription())
	}

	// 测试作者
	pkg.SetAuthor("Test Author")
	if pkg.GetAuthor() != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", pkg.GetAuthor())
	}

	// 测试许可证
	pkg.SetLicense("MIT")
	if pkg.GetLicense() != "MIT" {
		t.Errorf("Expected license 'MIT', got '%s'", pkg.GetLicense())
	}

	// 测试私有标志
	pkg.SetPrivate(true)
	if !pkg.IsPrivate() {
		t.Error("Expected private to be true")
	}

	// 测试主入口
	pkg.SetMain("index.js")
	if pkg.GetMain() != "index.js" {
		t.Errorf("Expected main 'index.js', got '%s'", pkg.GetMain())
	}

	// 测试主页
	pkg.SetHomepage("https://example.com")
	if pkg.GetHomepage() != "https://example.com" {
		t.Errorf("Expected homepage 'https://example.com', got '%s'", pkg.GetHomepage())
	}
}

func TestPackageJSONKeywords(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 添加关键词
	pkg.AddKeyword("test")
	pkg.AddKeyword("npm")
	pkg.AddKeyword("package")

	keywords := pkg.GetKeywords()
	if len(keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(keywords))
	}

	// 检查关键词内容
	expected := []string{"test", "npm", "package"}
	for i, keyword := range keywords {
		if keyword != expected[i] {
			t.Errorf("Expected keyword '%s', got '%s'", expected[i], keyword)
		}
	}

	// 添加重复关键词
	pkg.AddKeyword("test")
	keywords = pkg.GetKeywords()
	if len(keywords) != 3 {
		t.Errorf("Expected 3 keywords after adding duplicate, got %d", len(keywords))
	}

	// 移除关键词
	pkg.RemoveKeyword("npm")
	keywords = pkg.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected 2 keywords after removal, got %d", len(keywords))
	}

	// 设置关键词
	newKeywords := []string{"golang", "sdk"}
	pkg.SetKeywords(newKeywords)
	keywords = pkg.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected 2 keywords after setting, got %d", len(keywords))
	}
}

func TestPackageJSONDependencies(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 添加依赖
	pkg.AddDependency("lodash", "^4.17.21")
	pkg.AddDependency("axios", "^0.24.0")

	deps := pkg.GetDependencies()
	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	if deps["lodash"] != "^4.17.21" {
		t.Errorf("Expected lodash version '^4.17.21', got '%s'", deps["lodash"])
	}

	// 检查依赖存在
	if !pkg.HasDependency("lodash") {
		t.Error("Expected to have lodash dependency")
	}

	if pkg.HasDependency("nonexistent") {
		t.Error("Expected not to have nonexistent dependency")
	}

	// 移除依赖
	pkg.RemoveDependency("axios")
	deps = pkg.GetDependencies()
	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency after removal, got %d", len(deps))
	}

	// 添加开发依赖
	pkg.AddDevDependency("jest", "^27.0.0")
	pkg.AddDevDependency("eslint", "^8.0.0")

	devDeps := pkg.GetDevDependencies()
	if len(devDeps) != 2 {
		t.Errorf("Expected 2 dev dependencies, got %d", len(devDeps))
	}

	if !pkg.HasDevDependency("jest") {
		t.Error("Expected to have jest dev dependency")
	}

	// 移除开发依赖
	pkg.RemoveDevDependency("eslint")
	devDeps = pkg.GetDevDependencies()
	if len(devDeps) != 1 {
		t.Errorf("Expected 1 dev dependency after removal, got %d", len(devDeps))
	}
}

func TestPackageJSONScripts(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 添加脚本
	pkg.AddScript("test", "jest")
	pkg.AddScript("build", "webpack")
	pkg.AddScript("start", "node index.js")

	scripts := pkg.GetScripts()
	if len(scripts) != 3 {
		t.Errorf("Expected 3 scripts, got %d", len(scripts))
	}

	if scripts["test"] != "jest" {
		t.Errorf("Expected test script 'jest', got '%s'", scripts["test"])
	}

	// 检查脚本存在
	if !pkg.HasScript("test") {
		t.Error("Expected to have test script")
	}

	if pkg.HasScript("nonexistent") {
		t.Error("Expected not to have nonexistent script")
	}

	// 移除脚本
	pkg.RemoveScript("build")
	scripts = pkg.GetScripts()
	if len(scripts) != 2 {
		t.Errorf("Expected 2 scripts after removal, got %d", len(scripts))
	}
}

func TestPackageJSONRepository(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 设置仓库URL
	pkg.SetRepositoryURL("https://github.com/user/repo.git")

	repo := pkg.GetRepository()
	if repo == nil {
		t.Fatal("Expected repository to be set")
	}

	if repo.URL != "https://github.com/user/repo.git" {
		t.Errorf("Expected repository URL 'https://github.com/user/repo.git', got '%s'", repo.URL)
	}

	if repo.Type != "git" {
		t.Errorf("Expected repository type 'git', got '%s'", repo.Type)
	}

	// 设置完整仓库信息
	newRepo := &Repository{
		Type: "git",
		URL:  "https://github.com/user/another-repo.git",
	}
	pkg.SetRepository(newRepo)

	repo = pkg.GetRepository()
	if repo.URL != "https://github.com/user/another-repo.git" {
		t.Errorf("Expected repository URL 'https://github.com/user/another-repo.git', got '%s'", repo.URL)
	}
}

func TestPackageJSONBugs(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 设置bugs URL
	pkg.SetBugsURL("https://github.com/user/repo/issues")

	bugs := pkg.GetBugs()
	if bugs == nil {
		t.Fatal("Expected bugs to be set")
	}

	if bugs.URL != "https://github.com/user/repo/issues" {
		t.Errorf("Expected bugs URL 'https://github.com/user/repo/issues', got '%s'", bugs.URL)
	}

	// 设置完整bugs信息
	newBugs := &Bugs{
		URL:   "https://github.com/user/repo/issues",
		Email: "bugs@example.com",
	}
	pkg.SetBugs(newBugs)

	bugs = pkg.GetBugs()
	if bugs.Email != "bugs@example.com" {
		t.Errorf("Expected bugs email 'bugs@example.com', got '%s'", bugs.Email)
	}
}

func TestPackageJSONValidation(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试空包名验证
	err := pkg.Validate()
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// 设置名称但没有版本
	pkg.SetName("test-package")
	err = pkg.Validate()
	if err == nil {
		t.Error("Expected validation error for empty version")
	}

	// 设置有效的名称和版本
	pkg.SetVersion("1.0.0")
	err = pkg.Validate()
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	// 测试无效包名
	pkg.SetName("invalid package name")
	err = pkg.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid package name")
	}

	// 测试以.开头的包名
	pkg.SetName(".invalid")
	err = pkg.Validate()
	if err == nil {
		t.Error("Expected validation error for package name starting with .")
	}
}

func TestPackageJSONSaveLoad(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "package.json")

	// 创建并设置package.json
	pkg := NewPackageJSON(filePath)
	pkg.SetName("test-package")
	pkg.SetVersion("1.0.0")
	pkg.SetDescription("Test package")
	pkg.AddDependency("lodash", "^4.17.21")
	pkg.AddScript("test", "jest")

	// 保存
	err := pkg.Save()
	if err != nil {
		t.Fatalf("Failed to save package.json: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("package.json file was not created")
	}

	// 创建新的实例并加载
	newPkg := NewPackageJSON(filePath)
	err = newPkg.Load()
	if err != nil {
		t.Fatalf("Failed to load package.json: %v", err)
	}

	// 验证数据
	if newPkg.GetName() != "test-package" {
		t.Errorf("Expected name 'test-package', got '%s'", newPkg.GetName())
	}

	if newPkg.GetVersion() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", newPkg.GetVersion())
	}

	if !newPkg.HasDependency("lodash") {
		t.Error("Expected to have lodash dependency")
	}

	if !newPkg.HasScript("test") {
		t.Error("Expected to have test script")
	}
}

func TestPackageJSONComplexDependencies(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试复杂的依赖版本格式
	complexDependencies := map[string]string{
		"exact-version":    "1.2.3",
		"caret-range":      "^1.2.3",
		"tilde-range":      "~1.2.3",
		"greater-than":     ">1.2.3",
		"less-than":        "<2.0.0",
		"range":            ">=1.2.3 <2.0.0",
		"git-url":          "git+https://github.com/user/repo.git",
		"github-shorthand": "user/repo",
		"file-path":        "file:../local-package",
		"tarball-url":      "https://example.com/package.tgz",
		"latest-tag":       "latest",
		"beta-tag":         "beta",
		"scoped-package":   "@scope/package",
	}

	for name, version := range complexDependencies {
		pkg.AddDependency(name, version)
	}

	deps := pkg.GetDependencies()
	if len(deps) != len(complexDependencies) {
		t.Errorf("Expected %d dependencies, got %d", len(complexDependencies), len(deps))
	}

	for name, expectedVersion := range complexDependencies {
		if actualVersion, exists := deps[name]; !exists {
			t.Errorf("Expected dependency %s to exist", name)
		} else if actualVersion != expectedVersion {
			t.Errorf("Expected %s version %s, got %s", name, expectedVersion, actualVersion)
		}
	}
}

func TestPackageJSONPeerDependencies(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试同级依赖
	peerDeps := map[string]string{
		"react":      ">=16.8.0",
		"react-dom":  ">=16.8.0",
		"typescript": ">=4.0.0",
	}

	for name, version := range peerDeps {
		pkg.AddPeerDependency(name, version)
	}

	retrievedPeerDeps := pkg.GetPeerDependencies()
	if len(retrievedPeerDeps) != len(peerDeps) {
		t.Errorf("Expected %d peer dependencies, got %d", len(peerDeps), len(retrievedPeerDeps))
	}

	for name, expectedVersion := range peerDeps {
		if actualVersion, exists := retrievedPeerDeps[name]; !exists {
			t.Errorf("Expected peer dependency %s to exist", name)
		} else if actualVersion != expectedVersion {
			t.Errorf("Expected %s version %s, got %s", name, expectedVersion, actualVersion)
		}
	}

	// 测试移除同级依赖
	pkg.RemovePeerDependency("react")
	retrievedPeerDeps = pkg.GetPeerDependencies()
	if len(retrievedPeerDeps) != len(peerDeps)-1 {
		t.Errorf("Expected %d peer dependencies after removal, got %d", len(peerDeps)-1, len(retrievedPeerDeps))
	}

	if _, exists := retrievedPeerDeps["react"]; exists {
		t.Error("Expected react peer dependency to be removed")
	}
}

func TestPackageJSONOptionalDependencies(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试可选依赖
	optionalDeps := map[string]string{
		"fsevents":  "^2.3.0",
		"node-sass": "^6.0.0",
		"sharp":     "^0.32.0",
	}

	for name, version := range optionalDeps {
		pkg.AddOptionalDependency(name, version)
	}

	retrievedOptionalDeps := pkg.GetOptionalDependencies()
	if len(retrievedOptionalDeps) != len(optionalDeps) {
		t.Errorf("Expected %d optional dependencies, got %d", len(optionalDeps), len(retrievedOptionalDeps))
	}

	// 测试移除可选依赖
	pkg.RemoveOptionalDependency("fsevents")
	retrievedOptionalDeps = pkg.GetOptionalDependencies()
	if len(retrievedOptionalDeps) != len(optionalDeps)-1 {
		t.Errorf("Expected %d optional dependencies after removal, got %d", len(optionalDeps)-1, len(retrievedOptionalDeps))
	}
}

func TestPackageJSONComplexScripts(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试复杂的脚本命令
	complexScripts := map[string]string{
		"build":         "webpack --mode production",
		"dev":           "webpack serve --mode development",
		"test":          "jest --coverage",
		"test:watch":    "jest --watch",
		"test:ci":       "jest --ci --coverage --watchAll=false",
		"lint":          "eslint src/**/*.{js,ts,tsx}",
		"lint:fix":      "eslint src/**/*.{js,ts,tsx} --fix",
		"format":        "prettier --write src/**/*.{js,ts,tsx,json,css,md}",
		"type-check":    "tsc --noEmit",
		"pre-commit":    "lint-staged",
		"postinstall":   "husky install",
		"clean":         "rimraf dist",
		"build:analyze": "npm run build && npx webpack-bundle-analyzer dist/stats.json",
		"deploy":        "npm run build && gh-pages -d dist",
		"start:prod":    "serve -s dist -l 3000",
		"docker:build":  "docker build -t myapp .",
		"docker:run":    "docker run -p 3000:3000 myapp",
	}

	for name, command := range complexScripts {
		pkg.AddScript(name, command)
	}

	scripts := pkg.GetScripts()
	if len(scripts) != len(complexScripts) {
		t.Errorf("Expected %d scripts, got %d", len(complexScripts), len(scripts))
	}

	for name, expectedCommand := range complexScripts {
		if actualCommand, exists := scripts[name]; !exists {
			t.Errorf("Expected script %s to exist", name)
		} else if actualCommand != expectedCommand {
			t.Errorf("Expected script %s command %s, got %s", name, expectedCommand, actualCommand)
		}
	}

	// 测试脚本存在性检查
	if !pkg.HasScript("build") {
		t.Error("Expected to have build script")
	}

	if pkg.HasScript("nonexistent-script") {
		t.Error("Expected not to have nonexistent script")
	}
}

func TestPackageJSONAdvancedValidation(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试各种无效的包名
	invalidNames := []string{
		"",                        // 空名称
		"Invalid Name",            // 包含空格
		".invalid",                // 以点开头
		"_invalid",                // 以下划线开头
		"UPPERCASE",               // 全大写
		"name with spaces",        // 包含空格
		"name\twith\ttabs",        // 包含制表符
		"name\nwith\nnewlines",    // 包含换行符
		"name/with/slashes",       // 包含斜杠
		"name\\with\\backslashes", // 包含反斜杠
		strings.Repeat("a", 215),  // 超长名称
	}

	for _, name := range invalidNames {
		pkg.SetName(name)
		err := pkg.Validate()
		if err == nil {
			t.Errorf("Expected validation error for invalid name '%s'", name)
		}
	}

	// 测试各种无效的版本
	invalidVersions := []string{
		"",          // 空版本
		"invalid",   // 无效格式
		"1",         // 不完整
		"1.2",       // 不完整
		"v1.2.3",    // 带v前缀
		"1.2.3.4.5", // 过多部分
		"1.2.3-",    // 无效预发布
		"1.2.3+",    // 无效构建元数据
	}

	pkg.SetName("valid-name") // 设置有效名称
	for _, version := range invalidVersions {
		pkg.SetVersion(version)
		err := pkg.Validate()
		if err == nil {
			t.Errorf("Expected validation error for invalid version '%s'", version)
		}
	}

	// 测试有效的版本格式
	validVersions := []string{
		"1.0.0",
		"0.0.1",
		"10.20.30",
		"1.2.3-alpha",
		"1.2.3-alpha.1",
		"1.2.3-0.3.7",
		"1.2.3-x.7.z.92",
		"1.2.3+20130313144700",
		"1.2.3-beta+exp.sha.5114f85",
	}

	for _, version := range validVersions {
		pkg.SetVersion(version)
		err := pkg.Validate()
		if err != nil {
			t.Errorf("Expected no validation error for valid version '%s', got: %v", version, err)
		}
	}
}

func TestPackageJSONFileOperationsEdgeCases(t *testing.T) {
	// 测试不存在的文件
	nonexistentPath := "/nonexistent/path/package.json"
	pkg := NewPackageJSON(nonexistentPath)

	err := pkg.Load()
	if err == nil {
		t.Error("Expected error when loading nonexistent file")
	}

	// 测试无效的JSON文件
	tempDir := t.TempDir()
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")

	err = os.WriteFile(invalidJSONPath, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON file: %v", err)
	}

	invalidPkg := NewPackageJSON(invalidJSONPath)
	err = invalidPkg.Load()
	if err == nil {
		t.Error("Expected error when loading invalid JSON file")
	}

	// 测试权限问题（在支持的系统上）
	if runtime.GOOS != "windows" {
		readOnlyDir := filepath.Join(tempDir, "readonly")
		err = os.MkdirAll(readOnlyDir, 0444) // 只读目录
		if err != nil {
			t.Fatalf("Failed to create readonly directory: %v", err)
		}

		readOnlyPath := filepath.Join(readOnlyDir, "package.json")
		readOnlyPkg := NewPackageJSON(readOnlyPath)
		readOnlyPkg.SetName("test")
		readOnlyPkg.SetVersion("1.0.0")

		err = readOnlyPkg.Save()
		if err == nil {
			t.Error("Expected error when saving to readonly directory")
		}
	}
}

func TestPackageJSONComplexRepository(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试复杂的仓库配置
	repo := &Repository{
		Type: "git",
		URL:  "https://github.com/user/repo.git",
	}
	pkg.SetRepository(repo)

	retrievedRepo := pkg.GetRepository()
	if retrievedRepo == nil {
		t.Fatal("Expected repository to be set")
	}

	if retrievedRepo.Type != "git" {
		t.Errorf("Expected repository type 'git', got '%s'", retrievedRepo.Type)
	}

	if retrievedRepo.URL != "https://github.com/user/repo.git" {
		t.Errorf("Expected repository URL 'https://github.com/user/repo.git', got '%s'", retrievedRepo.URL)
	}

	// 测试通过URL设置仓库
	pkg.SetRepositoryURL("https://github.com/another/repo.git")
	retrievedRepo = pkg.GetRepository()

	if retrievedRepo.URL != "https://github.com/another/repo.git" {
		t.Errorf("Expected repository URL to be updated")
	}

	if retrievedRepo.Type != "git" {
		t.Errorf("Expected repository type to be inferred as 'git'")
	}

	// 测试非git URL
	pkg.SetRepositoryURL("https://example.com/repo.zip")
	retrievedRepo = pkg.GetRepository()

	if retrievedRepo.Type != "git" {
		t.Logf("Repository type for non-git URL: %s", retrievedRepo.Type)
	}
}

func TestPackageJSONComplexBugs(t *testing.T) {
	pkg := NewPackageJSON("/tmp/test.json")

	// 测试复杂的bugs配置
	bugs := &Bugs{
		URL:   "https://github.com/user/repo/issues",
		Email: "bugs@example.com",
	}
	pkg.SetBugs(bugs)

	retrievedBugs := pkg.GetBugs()
	if retrievedBugs == nil {
		t.Fatal("Expected bugs to be set")
	}

	if retrievedBugs.URL != "https://github.com/user/repo/issues" {
		t.Errorf("Expected bugs URL 'https://github.com/user/repo/issues', got '%s'", retrievedBugs.URL)
	}

	if retrievedBugs.Email != "bugs@example.com" {
		t.Errorf("Expected bugs email 'bugs@example.com', got '%s'", retrievedBugs.Email)
	}

	// 测试通过URL设置bugs
	pkg.SetBugsURL("https://gitlab.com/user/repo/-/issues")
	retrievedBugs = pkg.GetBugs()

	if retrievedBugs.URL != "https://gitlab.com/user/repo/-/issues" {
		t.Errorf("Expected bugs URL to be updated")
	}
}
