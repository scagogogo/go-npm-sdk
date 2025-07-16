package npm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PackageJSON package.json文件管理器
type PackageJSON struct {
	filePath string
	data     *Package
}

// NewPackageJSON 创建新的package.json管理器
func NewPackageJSON(filePath string) *PackageJSON {
	return &PackageJSON{
		filePath: filePath,
		data:     &Package{},
	}
}

// Load 加载package.json文件
func (p *PackageJSON) Load() error {
	if _, err := os.Stat(p.filePath); os.IsNotExist(err) {
		return fmt.Errorf("package.json file not found: %s", p.filePath)
	}

	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	if err := json.Unmarshal(data, p.data); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	return nil
}

// Save 保存package.json文件
func (p *PackageJSON) Save() error {
	// 确保目录存在
	dir := filepath.Dir(p.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	if err := os.WriteFile(p.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	return nil
}

// GetData 获取package数据
func (p *PackageJSON) GetData() *Package {
	return p.data
}

// SetData 设置package数据
func (p *PackageJSON) SetData(data *Package) {
	p.data = data
}

// GetName 获取包名
func (p *PackageJSON) GetName() string {
	return p.data.Name
}

// SetName 设置包名
func (p *PackageJSON) SetName(name string) {
	p.data.Name = name
}

// GetVersion 获取版本
func (p *PackageJSON) GetVersion() string {
	return p.data.Version
}

// SetVersion 设置版本
func (p *PackageJSON) SetVersion(version string) {
	p.data.Version = version
}

// GetDescription 获取描述
func (p *PackageJSON) GetDescription() string {
	return p.data.Description
}

// SetDescription 设置描述
func (p *PackageJSON) SetDescription(description string) {
	p.data.Description = description
}

// GetAuthor 获取作者
func (p *PackageJSON) GetAuthor() string {
	return p.data.Author
}

// SetAuthor 设置作者
func (p *PackageJSON) SetAuthor(author string) {
	p.data.Author = author
}

// GetLicense 获取许可证
func (p *PackageJSON) GetLicense() string {
	return p.data.License
}

// SetLicense 设置许可证
func (p *PackageJSON) SetLicense(license string) {
	p.data.License = license
}

// IsPrivate 检查是否为私有包
func (p *PackageJSON) IsPrivate() bool {
	return p.data.Private
}

// SetPrivate 设置私有标志
func (p *PackageJSON) SetPrivate(private bool) {
	p.data.Private = private
}

// GetMain 获取主入口文件
func (p *PackageJSON) GetMain() string {
	return p.data.Main
}

// SetMain 设置主入口文件
func (p *PackageJSON) SetMain(main string) {
	p.data.Main = main
}

// GetKeywords 获取关键词
func (p *PackageJSON) GetKeywords() []string {
	return p.data.Keywords
}

// SetKeywords 设置关键词
func (p *PackageJSON) SetKeywords(keywords []string) {
	p.data.Keywords = keywords
}

// AddKeyword 添加关键词
func (p *PackageJSON) AddKeyword(keyword string) {
	if p.data.Keywords == nil {
		p.data.Keywords = make([]string, 0)
	}
	
	// 检查是否已存在
	for _, k := range p.data.Keywords {
		if k == keyword {
			return
		}
	}
	
	p.data.Keywords = append(p.data.Keywords, keyword)
}

// RemoveKeyword 移除关键词
func (p *PackageJSON) RemoveKeyword(keyword string) {
	if p.data.Keywords == nil {
		return
	}
	
	for i, k := range p.data.Keywords {
		if k == keyword {
			p.data.Keywords = append(p.data.Keywords[:i], p.data.Keywords[i+1:]...)
			return
		}
	}
}

// GetDependencies 获取依赖
func (p *PackageJSON) GetDependencies() map[string]string {
	if p.data.Dependencies == nil {
		p.data.Dependencies = make(map[string]string)
	}
	return p.data.Dependencies
}

// GetDevDependencies 获取开发依赖
func (p *PackageJSON) GetDevDependencies() map[string]string {
	if p.data.DevDeps == nil {
		p.data.DevDeps = make(map[string]string)
	}
	return p.data.DevDeps
}

// GetOptionalDependencies 获取可选依赖
func (p *PackageJSON) GetOptionalDependencies() map[string]string {
	if p.data.OptionalDeps == nil {
		p.data.OptionalDeps = make(map[string]string)
	}
	return p.data.OptionalDeps
}

// GetPeerDependencies 获取同级依赖
func (p *PackageJSON) GetPeerDependencies() map[string]string {
	if p.data.PeerDeps == nil {
		p.data.PeerDeps = make(map[string]string)
	}
	return p.data.PeerDeps
}

// AddDependency 添加依赖
func (p *PackageJSON) AddDependency(name, version string) {
	if p.data.Dependencies == nil {
		p.data.Dependencies = make(map[string]string)
	}
	p.data.Dependencies[name] = version
}

// AddDevDependency 添加开发依赖
func (p *PackageJSON) AddDevDependency(name, version string) {
	if p.data.DevDeps == nil {
		p.data.DevDeps = make(map[string]string)
	}
	p.data.DevDeps[name] = version
}

// AddOptionalDependency 添加可选依赖
func (p *PackageJSON) AddOptionalDependency(name, version string) {
	if p.data.OptionalDeps == nil {
		p.data.OptionalDeps = make(map[string]string)
	}
	p.data.OptionalDeps[name] = version
}

// AddPeerDependency 添加同级依赖
func (p *PackageJSON) AddPeerDependency(name, version string) {
	if p.data.PeerDeps == nil {
		p.data.PeerDeps = make(map[string]string)
	}
	p.data.PeerDeps[name] = version
}

// RemoveDependency 移除依赖
func (p *PackageJSON) RemoveDependency(name string) {
	if p.data.Dependencies != nil {
		delete(p.data.Dependencies, name)
	}
}

// RemoveDevDependency 移除开发依赖
func (p *PackageJSON) RemoveDevDependency(name string) {
	if p.data.DevDeps != nil {
		delete(p.data.DevDeps, name)
	}
}

// RemoveOptionalDependency 移除可选依赖
func (p *PackageJSON) RemoveOptionalDependency(name string) {
	if p.data.OptionalDeps != nil {
		delete(p.data.OptionalDeps, name)
	}
}

// RemovePeerDependency 移除同级依赖
func (p *PackageJSON) RemovePeerDependency(name string) {
	if p.data.PeerDeps != nil {
		delete(p.data.PeerDeps, name)
	}
}

// HasDependency 检查是否有指定依赖
func (p *PackageJSON) HasDependency(name string) bool {
	if p.data.Dependencies == nil {
		return false
	}
	_, exists := p.data.Dependencies[name]
	return exists
}

// HasDevDependency 检查是否有指定开发依赖
func (p *PackageJSON) HasDevDependency(name string) bool {
	if p.data.DevDeps == nil {
		return false
	}
	_, exists := p.data.DevDeps[name]
	return exists
}

// GetScripts 获取脚本
func (p *PackageJSON) GetScripts() map[string]string {
	if p.data.Scripts == nil {
		p.data.Scripts = make(map[string]string)
	}
	return p.data.Scripts
}

// AddScript 添加脚本
func (p *PackageJSON) AddScript(name, command string) {
	if p.data.Scripts == nil {
		p.data.Scripts = make(map[string]string)
	}
	p.data.Scripts[name] = command
}

// RemoveScript 移除脚本
func (p *PackageJSON) RemoveScript(name string) {
	if p.data.Scripts != nil {
		delete(p.data.Scripts, name)
	}
}

// HasScript 检查是否有指定脚本
func (p *PackageJSON) HasScript(name string) bool {
	if p.data.Scripts == nil {
		return false
	}
	_, exists := p.data.Scripts[name]
	return exists
}

// GetRepository 获取仓库信息
func (p *PackageJSON) GetRepository() *Repository {
	return p.data.Repository
}

// SetRepository 设置仓库信息
func (p *PackageJSON) SetRepository(repo *Repository) {
	p.data.Repository = repo
}

// SetRepositoryURL 设置仓库URL
func (p *PackageJSON) SetRepositoryURL(url string) {
	if p.data.Repository == nil {
		p.data.Repository = &Repository{}
	}
	p.data.Repository.URL = url
	
	// 根据URL推断类型
	if strings.Contains(url, "git") {
		p.data.Repository.Type = "git"
	}
}

// GetBugs 获取bug报告信息
func (p *PackageJSON) GetBugs() *Bugs {
	return p.data.Bugs
}

// SetBugs 设置bug报告信息
func (p *PackageJSON) SetBugs(bugs *Bugs) {
	p.data.Bugs = bugs
}

// SetBugsURL 设置bug报告URL
func (p *PackageJSON) SetBugsURL(url string) {
	if p.data.Bugs == nil {
		p.data.Bugs = &Bugs{}
	}
	p.data.Bugs.URL = url
}

// GetHomepage 获取主页
func (p *PackageJSON) GetHomepage() string {
	return p.data.Homepage
}

// SetHomepage 设置主页
func (p *PackageJSON) SetHomepage(homepage string) {
	p.data.Homepage = homepage
}

// Validate 验证package.json数据
func (p *PackageJSON) Validate() error {
	if p.data.Name == "" {
		return NewValidationError("name", "", "package name is required")
	}
	
	if p.data.Version == "" {
		return NewValidationError("version", "", "package version is required")
	}
	
	// 验证包名格式
	if !isValidPackageName(p.data.Name) {
		return NewValidationError("name", p.data.Name, "invalid package name format")
	}
	
	// 验证版本格式
	if !isValidVersion(p.data.Version) {
		return NewValidationError("version", p.data.Version, "invalid version format")
	}
	
	return nil
}

// isValidPackageName 验证包名格式
func isValidPackageName(name string) bool {
	if name == "" {
		return false
	}
	
	// 简单验证：不能包含空格，不能以.或_开头
	if strings.Contains(name, " ") {
		return false
	}
	
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return false
	}
	
	return true
}

// isValidVersion 验证版本格式
func isValidVersion(version string) bool {
	if version == "" {
		return false
	}
	
	// 简单验证：应该包含数字和点
	return strings.Contains(version, ".")
}
