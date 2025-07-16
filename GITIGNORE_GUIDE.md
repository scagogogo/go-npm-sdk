# .gitignore 配置说明

本项目的 `.gitignore` 文件经过精心配置，确保只有必要的源代码和文档被纳入版本控制，而忽略所有生成的文件、临时文件和敏感信息。

## 忽略的文件类型

### Go 语言相关

#### 编译产物
- `*.exe`, `*.exe~` - Windows 可执行文件
- `*.dll` - Windows 动态链接库
- `*.so` - Linux 共享对象文件
- `*.dylib` - macOS 动态库
- `*.test` - Go 测试二进制文件

#### 测试和性能分析
- `*.out`, `coverage.out`, `coverage.html`, `*.cover` - 测试覆盖率文件
- `*.prof`, `*.pprof` - 性能分析文件
- `*.trace` - 执行跟踪文件
- `*.bench` - 基准测试结果
- `*.mprof` - 内存分析文件
- `*.cprof` - CPU 分析文件

#### Go 工作区和依赖
- `go.work`, `go.work.sum` - Go 工作区文件
- `vendor/` - 依赖目录

### 开发环境相关

#### IDE 和编辑器
- `.vscode/` - Visual Studio Code 配置
- `.idea/` - IntelliJ IDEA 配置
- `*.swp`, `*.swo`, `*~` - Vim 临时文件
- `.project`, `.settings/` - Eclipse 配置

#### 操作系统生成的文件
- `.DS_Store`, `.DS_Store?`, `._*` - macOS 系统文件
- `.Spotlight-V100`, `.Trashes` - macOS 系统目录
- `ehthumbs.db`, `Thumbs.db` - Windows 缩略图缓存
- `desktop.ini` - Windows 桌面配置

### 构建和部署

#### 构建产物
- `/bin/`, `/build/`, `/dist/` - 构建输出目录

#### 临时文件
- `*.tmp`, `*.temp` - 临时文件
- `*.bak`, `*.backup` - 备份文件
- `*.log` - 日志文件
- `logs/` - 日志目录

### 配置和环境

#### 环境配置
- `.env`, `.env.local`, `.env.*.local` - 环境变量文件
- `config.local.*` - 本地配置文件
- `*.local.json`, `*.local.yaml`, `*.local.yml` - 本地配置文件

#### 运行时数据
- `*.pid`, `*.seed`, `*.pid.lock` - 进程ID文件

### Node.js 相关（用于测试npm功能）

#### Node.js 文件
- `node_modules/` - Node.js 依赖目录
- `package-lock.json`, `yarn.lock` - 锁定文件
- `npm-debug.log*`, `yarn-debug.log*`, `yarn-error.log*` - 调试日志
- `.npm`, `.yarn-integrity` - npm/yarn 缓存

### 测试和文档

#### 测试目录
- `/testdata/`, `/tmp/`, `/temp/` - 测试数据目录
- `test-output/` - 测试输出目录

#### 文档构建
- `/docs/_build/`, `/docs/site/` - 文档构建输出

### 其他文件类型

#### 缓存文件
- `.cache/`, `*.cache` - 缓存目录和文件

#### 压缩文件
- `*.tar`, `*.tar.gz`, `*.zip`, `*.rar` - 压缩包

#### 证书和密钥
- `*.pem`, `*.key`, `*.crt`, `*.p12` - 证书和密钥文件

#### 数据库文件
- `*.db`, `*.sqlite`, `*.sqlite3` - 数据库文件

## 项目特定的忽略规则

### 便携版 Node.js 安装
- `.go-npm-sdk/` - SDK 便携版安装目录
- `.go-npm-sdk-example/` - 示例便携版安装目录

### 测试项目目录
- `npm-sdk-test/` - 基本示例创建的测试项目
- `portable-npm-test/` - 便携版示例创建的测试项目
- `test-project/` - 其他测试项目

### 本地开发文件
- `*.local` - 本地文件
- `local.*` - 本地配置文件

## 使用建议

### 添加新的忽略规则
如果需要添加新的忽略规则，请按照以下原则：

1. **按类型分组** - 将相关的文件类型放在一起
2. **添加注释** - 为每个部分添加清晰的注释
3. **使用通配符** - 合理使用 `*` 和 `?` 通配符
4. **避免过度忽略** - 不要忽略可能需要的源代码文件

### 检查忽略效果
使用以下命令检查文件是否被正确忽略：

```bash
# 查看当前状态
git status

# 查看所有文件（包括被忽略的）
git status --ignored

# 检查特定文件是否被忽略
git check-ignore filename
```

### 强制添加被忽略的文件
如果确实需要添加某个被忽略的文件：

```bash
git add -f filename
```

## 维护说明

这个 `.gitignore` 文件应该随着项目的发展而更新：

1. **新增文件类型** - 当项目引入新的工具或框架时
2. **清理无用规则** - 移除不再需要的忽略规则
3. **优化性能** - 合并相似的规则，提高匹配效率

## 相关资源

- [Git 官方文档 - gitignore](https://git-scm.com/docs/gitignore)
- [GitHub gitignore 模板](https://github.com/github/gitignore)
- [gitignore.io](https://www.toptal.com/developers/gitignore) - 在线生成工具
