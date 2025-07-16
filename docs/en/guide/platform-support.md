# Platform Support

The Go NPM SDK supports multiple platforms and architectures with platform-specific optimizations.

## Supported Platforms

### Windows
- **Versions**: Windows 10, Windows 11, Windows Server 2019, Windows Server 2022
- **Architectures**: x86_64 (amd64), x86 (386)
- **Installation Methods**: Chocolatey, winget, official installer

### macOS
- **Versions**: macOS 10.15+ (Catalina and later)
- **Architectures**: Intel (x86_64), Apple Silicon (ARM64)
- **Installation Methods**: Homebrew, MacPorts, official installer

### Linux
- **Distributions**: Ubuntu, Debian, CentOS, RHEL, Fedora, SUSE, Arch, Alpine
- **Architectures**: x86_64 (amd64), ARM64, ARM, x86 (386)
- **Installation Methods**: Package managers (apt, yum, pacman, etc.), official installer

## Platform Detection

The SDK automatically detects your platform:

```go
import "github.com/scagogogo/go-npm-sdk/pkg/platform"

detector := platform.NewDetector()
info, err := detector.Detect()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Platform: %s\n", info.Platform)
fmt.Printf("Architecture: %s\n", info.Architecture)

if info.IsLinux() {
    fmt.Printf("Distribution: %s\n", info.Distribution)
    fmt.Printf("Version: %s\n", info.Version)
}
```

## Installation Methods by Platform

### Windows

1. **Chocolatey** (Recommended)
   ```bash
   choco install nodejs
   ```

2. **winget**
   ```bash
   winget install OpenJS.NodeJS
   ```

3. **Official Installer**
   - Downloads from nodejs.org
   - Automatic installation

### macOS

1. **Homebrew** (Recommended)
   ```bash
   brew install node
   ```

2. **MacPorts**
   ```bash
   sudo port install nodejs18
   ```

3. **Official Installer**
   - Downloads from nodejs.org
   - Automatic installation

### Linux

1. **Package Managers**
   
   **Ubuntu/Debian:**
   ```bash
   sudo apt update
   sudo apt install nodejs npm
   ```
   
   **CentOS/RHEL/Fedora:**
   ```bash
   sudo yum install nodejs npm
   # or
   sudo dnf install nodejs npm
   ```
   
   **Arch Linux:**
   ```bash
   sudo pacman -S nodejs npm
   ```

2. **Official Installer**
   - Downloads from nodejs.org
   - Automatic installation

## Platform-Specific Features

### Windows
- Support for Windows-style paths
- Integration with Windows package managers
- PowerShell and Command Prompt compatibility

### macOS
- Support for both Intel and Apple Silicon
- Integration with Homebrew and MacPorts
- Xcode command line tools compatibility

### Linux
- Automatic distribution detection
- Support for multiple package managers
- Container-friendly installation

## Troubleshooting

### Common Issues

1. **Permission Denied**
   - Use sudo on Linux/macOS for global installations
   - Run as Administrator on Windows

2. **Path Issues**
   - Ensure npm is in your PATH
   - Use absolute paths if needed

3. **Architecture Mismatch**
   - Verify your architecture with the SDK
   - Use platform-specific downloads

### Platform-Specific Solutions

**Windows:**
- Enable Developer Mode for symlinks
- Use Windows Subsystem for Linux (WSL) if needed

**macOS:**
- Install Xcode command line tools
- Use Homebrew for easier package management

**Linux:**
- Update package repositories
- Install build-essential tools if needed

## Next Steps

- [Getting Started](./getting-started.md) - Basic usage
- [Configuration](./configuration.md) - Advanced configuration
- [API Reference](/en/api/overview.md) - Complete API documentation
