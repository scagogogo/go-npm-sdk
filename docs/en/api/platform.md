# Platform Package API

The `pkg/platform` package provides platform detection and file download capabilities with cross-platform support for Windows, macOS, and Linux.

## Detector

The detector component identifies the current operating system, architecture, and Linux distribution.

### NewDetector

```go
func NewDetector() *Detector
```

Creates a new platform detector.

### Detect

```go
func (d *Detector) Detect() (*Info, error)
```

Detects the current platform information.

**Returns:**
- `*Info`: Platform information structure
- `error`: Error if detection fails

**Example:**
```go
detector := platform.NewDetector()
info, err := detector.Detect()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Platform: %s\n", info.Platform)
fmt.Printf("Architecture: %s\n", info.Architecture)
fmt.Printf("Distribution: %s\n", info.Distribution)
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Kernel: %s\n", info.Kernel)
```

## Platform Information

### Info Structure

```go
type Info struct {
    Platform     Platform     `json:"platform"`
    Architecture Architecture `json:"architecture"`
    Distribution Distribution `json:"distribution,omitempty"`
    Version      string       `json:"version,omitempty"`
    Kernel       string       `json:"kernel,omitempty"`
}
```

### Platform Methods

```go
func (i *Info) IsWindows() bool
func (i *Info) IsMacOS() bool
func (i *Info) IsLinux() bool
func (i *Info) IsX86() bool
func (i *Info) IsARM() bool
func (i *Info) String() string
```

**Example:**
```go
info, _ := detector.Detect()

if info.IsLinux() {
    fmt.Println("Running on Linux")
    if info.Distribution == platform.Ubuntu {
        fmt.Println("Ubuntu distribution detected")
    }
}

if info.IsX86() {
    fmt.Println("x86 architecture")
} else if info.IsARM() {
    fmt.Println("ARM architecture")
}

fmt.Printf("Platform string: %s\n", info.String())
```

## Platform Constants

### Platform Types

```go
const (
    Windows Platform = "windows"
    MacOS   Platform = "darwin"
    Linux   Platform = "linux"
    Unknown Platform = "unknown"
)
```

### Architecture Types

```go
const (
    AMD64 Architecture = "amd64"
    I386  Architecture = "386"
    ARM64 Architecture = "arm64"
    ARM   Architecture = "arm"
)
```

### Linux Distributions

```go
const (
    Ubuntu        Distribution = "ubuntu"
    Debian        Distribution = "debian"
    CentOS        Distribution = "centos"
    RHEL          Distribution = "rhel"
    Fedora        Distribution = "fedora"
    SUSE          Distribution = "suse"
    Arch          Distribution = "arch"
    Alpine        Distribution = "alpine"
    UnknownDistro Distribution = "unknown"
)
```

## Downloader

The downloader component provides file download capabilities with progress tracking and retry mechanisms.

### NewDownloader

```go
func NewDownloader() *Downloader
```

Creates a new downloader with default configuration.

### Download

```go
func (d *Downloader) Download(ctx context.Context, options DownloadOptions) (*DownloadResult, error)
```

Downloads a file with the specified options.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `options` (DownloadOptions): Download configuration

**Returns:**
- `*DownloadResult`: Download result information
- `error`: Error if download fails

### DownloadWithRetry

```go
func (d *Downloader) DownloadWithRetry(ctx context.Context, options DownloadOptions, maxRetries int) (*DownloadResult, error)
```

Downloads a file with retry mechanism.

**Example:**
```go
downloader := platform.NewDownloader()
ctx := context.Background()

options := platform.DownloadOptions{
    URL:         "https://nodejs.org/dist/v18.17.0/node-v18.17.0-linux-x64.tar.xz",
    Destination: "/tmp/node-v18.17.0-linux-x64.tar.xz",
    Timeout:     30 * time.Minute,
    UserAgent:   "Go-NPM-SDK/1.0",
    Headers: map[string]string{
        "Accept": "application/octet-stream",
    },
    Progress: func(downloaded, total int64) {
        percentage := float64(downloaded) / float64(total) * 100
        fmt.Printf("Download progress: %.2f%%\n", percentage)
    },
}

result, err := downloader.DownloadWithRetry(ctx, options, 3)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Downloaded %d bytes in %v\n", result.Size, result.Duration)
```

### DownloadOptions

```go
type DownloadOptions struct {
    URL         string                           `json:"url"`
    Destination string                           `json:"destination"`
    Timeout     time.Duration                    `json:"timeout,omitempty"`
    UserAgent   string                           `json:"user_agent,omitempty"`
    Headers     map[string]string                `json:"headers,omitempty"`
    Progress    func(downloaded, total int64)    `json:"-"`
}
```

### DownloadResult

```go
type DownloadResult struct {
    FilePath string        `json:"file_path"`
    Size     int64         `json:"size"`
    Duration time.Duration `json:"duration"`
    Success  bool          `json:"success"`
    Error    string        `json:"error,omitempty"`
}
```

## Node.js Downloader

Specialized downloader for Node.js releases with platform-specific URL generation.

### NewNodeJSDownloader

```go
func NewNodeJSDownloader() *NodeJSDownloader
```

Creates a new Node.js downloader.

### GetDownloadURL

```go
func (n *NodeJSDownloader) GetDownloadURL(version string, platform Platform, arch Architecture) string
```

Generates the download URL for a specific Node.js version and platform.

**Parameters:**
- `version` (string): Node.js version (e.g., "18.17.0")
- `platform` (Platform): Target platform
- `arch` (Architecture): Target architecture

**Returns:**
- `string`: Download URL for the specified version and platform

### GetLatestVersion

```go
func (n *NodeJSDownloader) GetLatestVersion(ctx context.Context) (string, error)
```

Retrieves the latest Node.js version from the official API.

### Download

```go
func (n *NodeJSDownloader) Download(ctx context.Context, version string, platform Platform, arch Architecture, destination string) (*DownloadResult, error)
```

Downloads a specific Node.js version for the target platform.

**Example:**
```go
nodeDownloader := platform.NewNodeJSDownloader()
ctx := context.Background()

// Get latest version
latest, err := nodeDownloader.GetLatestVersion(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Latest Node.js version: %s\n", latest)

// Get download URL
detector := platform.NewDetector()
info, _ := detector.Detect()

url := nodeDownloader.GetDownloadURL("18.17.0", info.Platform, info.Architecture)
fmt.Printf("Download URL: %s\n", url)

// Download Node.js
result, err := nodeDownloader.Download(ctx, "18.17.0", info.Platform, info.Architecture, "/tmp/node.tar.xz")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Downloaded Node.js: %s (%d bytes)\n", result.FilePath, result.Size)
```

## Platform-Specific URLs

The Node.js downloader generates platform-specific URLs:

### Windows
- **AMD64**: `node-v{version}-win-x64.zip`
- **I386**: `node-v{version}-win-x86.zip`

### macOS
- **AMD64**: `node-v{version}-darwin-x64.tar.gz`
- **ARM64**: `node-v{version}-darwin-arm64.tar.gz`

### Linux
- **AMD64**: `node-v{version}-linux-x64.tar.xz`
- **ARM64**: `node-v{version}-linux-arm64.tar.xz`
- **ARM**: `node-v{version}-linux-armv7l.tar.xz`
- **I386**: `node-v{version}-linux-x86.tar.xz`

## Error Handling

The platform package provides specific error types:

```go
type PlatformError struct {
    Platform string
    Reason   string
    Err      error
}

type DownloadError struct {
    URL    string
    Reason string
    Err    error
}
```

**Example:**
```go
result, err := downloader.Download(ctx, options)
if err != nil {
    if platformErr, ok := err.(*platform.PlatformError); ok {
        fmt.Printf("Platform error on %s: %s\n", platformErr.Platform, platformErr.Reason)
    } else if downloadErr, ok := err.(*platform.DownloadError); ok {
        fmt.Printf("Download error for %s: %s\n", downloadErr.URL, downloadErr.Reason)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Best Practices

1. **Check platform compatibility**: Always verify platform support before operations
2. **Use appropriate timeouts**: Set reasonable timeouts for download operations
3. **Handle progress**: Implement progress callbacks for better user experience
4. **Retry on failure**: Use retry mechanisms for network operations
5. **Validate URLs**: Ensure download URLs are valid before attempting downloads
6. **Clean up on failure**: Remove incomplete downloads on failure
7. **Use contexts**: Always pass contexts for cancellation support

## Integration Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/platform"
)

func main() {
    // Detect platform
    detector := platform.NewDetector()
    info, err := detector.Detect()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Detected platform: %s\n", info.String())
    
    // Download Node.js for current platform
    nodeDownloader := platform.NewNodeJSDownloader()
    ctx := context.Background()
    
    result, err := nodeDownloader.Download(
        ctx,
        "18.17.0",
        info.Platform,
        info.Architecture,
        "/tmp/node.tar.xz",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Downloaded Node.js: %s (%d bytes)\n", result.FilePath, result.Size)
}
```
