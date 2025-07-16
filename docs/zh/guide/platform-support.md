# 平台支持

Go NPM SDK支持多个平台和架构，具有平台特定的优化。

## 支持的平台

### Windows
- **版本**: Windows 10, Windows 11, Windows Server 2019, Windows Server 2022
- **架构**: x86_64 (amd64), x86 (386)
- **安装方法**: Chocolatey, winget, 官方安装程序

### macOS
- **版本**: macOS 10.15+ (Catalina及更高版本)
- **架构**: Intel (x86_64), Apple Silicon (ARM64)
- **安装方法**: Homebrew, MacPorts, 官方安装程序

### Linux
- **发行版**: Ubuntu, Debian, CentOS, RHEL, Fedora, SUSE, Arch, Alpine
- **架构**: x86_64 (amd64), ARM64, ARM, x86 (386)
- **安装方法**: 包管理器 (apt, yum, pacman等), 官方安装程序

## 平台检测

SDK自动检测您的平台：

```go
import "github.com/scagogogo/go-npm-sdk/pkg/platform"

detector := platform.NewDetector()
info, err := detector.Detect()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("平台: %s\n", info.Platform)
fmt.Printf("架构: %s\n", info.Architecture)

if info.IsLinux() {
    fmt.Printf("发行版: %s\n", info.Distribution)
    fmt.Printf("版本: %s\n", info.Version)
}
```

## 按平台的安装方法

### Windows

1. **Chocolatey** (推荐)
   ```bash
   choco install nodejs
   ```

2. **winget**
   ```bash
   winget install OpenJS.NodeJS
   ```

3. **官方安装程序**
   - 从nodejs.org下载
   - 自动安装

### macOS

1. **Homebrew** (推荐)
   ```bash
   brew install node
   ```

2. **MacPorts**
   ```bash
   sudo port install nodejs18
   ```

3. **官方安装程序**
   - 从nodejs.org下载
   - 自动安装

### Linux

1. **包管理器**
   
   **Ubuntu/Debian:**
   ```bash
   sudo apt update
   sudo apt install nodejs npm
   ```
   
   **CentOS/RHEL/Fedora:**
   ```bash
   sudo yum install nodejs npm
   # 或
   sudo dnf install nodejs npm
   ```
   
   **Arch Linux:**
   ```bash
   sudo pacman -S nodejs npm
   ```

2. **官方安装程序**
   - 从nodejs.org下载
   - 自动安装

## 平台特定功能

### Windows
- 支持Windows风格的路径
- 与Windows包管理器集成
- PowerShell和命令提示符兼容性

### macOS
- 支持Intel和Apple Silicon
- 与Homebrew和MacPorts集成
- Xcode命令行工具兼容性

### Linux
- 自动发行版检测
- 支持多个包管理器
- 容器友好安装

## 故障排除

### 常见问题

1. **权限被拒绝**
   - 在Linux/macOS上使用sudo进行全局安装
   - 在Windows上以管理员身份运行

2. **路径问题**
   - 确保npm在您的PATH中
   - 如需要使用绝对路径

3. **架构不匹配**
   - 使用SDK验证您的架构
   - 使用平台特定的下载

### 平台特定解决方案

**Windows:**
- 启用开发者模式以支持符号链接
- 如需要使用Windows子系统Linux (WSL)

**macOS:**
- 安装Xcode命令行工具
- 使用Homebrew以便于包管理

**Linux:**
- 更新包仓库
- 如需要安装build-essential工具

## 下一步

- [快速开始](./getting-started.md) - 基本用法
- [配置](./configuration.md) - 高级配置
- [API参考](/zh/api/overview.md) - 完整API文档
