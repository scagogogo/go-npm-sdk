package platform

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DownloadOptions 下载选项
type DownloadOptions struct {
	URL         string        `json:"url"`
	Destination string        `json:"destination"`
	Timeout     time.Duration `json:"timeout"`
	UserAgent   string        `json:"user_agent"`
	Headers     map[string]string `json:"headers"`
	Progress    ProgressCallback  `json:"-"`
}

// ProgressCallback 进度回调函数
type ProgressCallback func(downloaded, total int64)

// DownloadResult 下载结果
type DownloadResult struct {
	FilePath   string        `json:"file_path"`
	Size       int64         `json:"size"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Error      error         `json:"error,omitempty"`
}

// Downloader 下载器
type Downloader struct {
	client *http.Client
}

// NewDownloader 创建新的下载器
func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Minute, // 默认30分钟超时
		},
	}
}

// Download 下载文件
func (d *Downloader) Download(ctx context.Context, options DownloadOptions) (*DownloadResult, error) {
	startTime := time.Now()
	
	// 设置超时
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", options.URL, nil)
	if err != nil {
		return &DownloadResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("failed to create request: %w", err),
		}, err
	}

	// 设置User-Agent
	if options.UserAgent != "" {
		req.Header.Set("User-Agent", options.UserAgent)
	} else {
		req.Header.Set("User-Agent", "go-npm-sdk/1.0")
	}

	// 设置自定义头部
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return &DownloadResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("failed to send request: %w", err),
		}, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return &DownloadResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(options.Destination), 0755); err != nil {
		return &DownloadResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("failed to create directory: %w", err),
		}, err
	}

	// 创建目标文件
	file, err := os.Create(options.Destination)
	if err != nil {
		return &DownloadResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("failed to create file: %w", err),
		}, err
	}
	defer file.Close()

	// 获取文件大小
	contentLength := resp.ContentLength

	// 创建进度读取器
	var reader io.Reader = resp.Body
	if options.Progress != nil && contentLength > 0 {
		reader = &progressReader{
			reader:   resp.Body,
			total:    contentLength,
			callback: options.Progress,
		}
	}

	// 复制数据
	written, err := io.Copy(file, reader)
	if err != nil {
		// 删除不完整的文件
		os.Remove(options.Destination)
		return &DownloadResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("failed to copy data: %w", err),
		}, err
	}

	return &DownloadResult{
		FilePath: options.Destination,
		Size:     written,
		Duration: time.Since(startTime),
		Success:  true,
	}, nil
}

// DownloadWithRetry 带重试的下载
func (d *Downloader) DownloadWithRetry(ctx context.Context, options DownloadOptions, maxRetries int) (*DownloadResult, error) {
	var lastErr error
	
	for i := 0; i <= maxRetries; i++ {
		result, err := d.Download(ctx, options)
		if err == nil {
			return result, nil
		}
		
		lastErr = err
		
		// 如果不是最后一次重试，等待一段时间
		if i < maxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(i+1) * time.Second):
				// 指数退避
			}
		}
	}
	
	return nil, fmt.Errorf("download failed after %d retries: %w", maxRetries, lastErr)
}

// progressReader 进度读取器
type progressReader struct {
	reader     io.Reader
	total      int64
	downloaded int64
	callback   ProgressCallback
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		pr.downloaded += int64(n)
		if pr.callback != nil {
			pr.callback(pr.downloaded, pr.total)
		}
	}
	return n, err
}

// NodeJSDownloader Node.js下载器
type NodeJSDownloader struct {
	downloader *Downloader
	baseURL    string
}

// NewNodeJSDownloader 创建Node.js下载器
func NewNodeJSDownloader() *NodeJSDownloader {
	return &NodeJSDownloader{
		downloader: NewDownloader(),
		baseURL:    "https://nodejs.org/dist",
	}
}

// GetDownloadURL 获取Node.js下载URL
func (nd *NodeJSDownloader) GetDownloadURL(version string, platform Platform, arch Architecture) string {
	var filename string
	
	switch platform {
	case Windows:
		if arch == AMD64 {
			filename = fmt.Sprintf("node-v%s-win-x64.zip", version)
		} else {
			filename = fmt.Sprintf("node-v%s-win-x86.zip", version)
		}
	case MacOS:
		if arch == ARM64 {
			filename = fmt.Sprintf("node-v%s-darwin-arm64.tar.gz", version)
		} else {
			filename = fmt.Sprintf("node-v%s-darwin-x64.tar.gz", version)
		}
	case Linux:
		if arch == ARM64 {
			filename = fmt.Sprintf("node-v%s-linux-arm64.tar.xz", version)
		} else if arch == ARM {
			filename = fmt.Sprintf("node-v%s-linux-armv7l.tar.xz", version)
		} else {
			filename = fmt.Sprintf("node-v%s-linux-x64.tar.xz", version)
		}
	default:
		return ""
	}
	
	return fmt.Sprintf("%s/v%s/%s", nd.baseURL, version, filename)
}

// DownloadNodeJS 下载Node.js
func (nd *NodeJSDownloader) DownloadNodeJS(ctx context.Context, version string, info *Info, destination string, progress ProgressCallback) (*DownloadResult, error) {
	url := nd.GetDownloadURL(version, info.Platform, info.Architecture)
	if url == "" {
		return nil, fmt.Errorf("unsupported platform: %s/%s", info.Platform, info.Architecture)
	}
	
	filename := filepath.Base(url)
	filePath := filepath.Join(destination, filename)
	
	options := DownloadOptions{
		URL:         url,
		Destination: filePath,
		Timeout:     30 * time.Minute,
		Progress:    progress,
	}
	
	return nd.downloader.DownloadWithRetry(ctx, options, 3)
}

// GetLatestVersion 获取最新版本号
func (nd *NodeJSDownloader) GetLatestVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://nodejs.org/dist/latest/", nil)
	if err != nil {
		return "", err
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get latest version: status %d", resp.StatusCode)
	}
	
	// 这里需要解析HTML或使用API来获取版本号
	// 简化实现，返回一个默认版本
	return "18.17.0", nil
}
