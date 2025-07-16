package platform

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewDownloader(t *testing.T) {
	downloader := NewDownloader()
	if downloader == nil {
		t.Fatal("NewDownloader() returned nil")
	}

	if downloader.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	// 检查默认超时
	if downloader.client.Timeout != 30*time.Minute {
		t.Errorf("Expected default timeout 30m, got %v", downloader.client.Timeout)
	}
}

func TestDownloadOptions(t *testing.T) {
	options := DownloadOptions{
		URL:         "https://example.com/file.zip",
		Destination: "/tmp/file.zip",
		Timeout:     10 * time.Second,
		UserAgent:   "test-agent",
		Headers:     map[string]string{"X-Test": "value"},
		Progress:    func(downloaded, total int64) {},
	}

	if options.URL != "https://example.com/file.zip" {
		t.Errorf("Expected URL 'https://example.com/file.zip', got '%s'", options.URL)
	}

	if options.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", options.Timeout)
	}

	if options.UserAgent != "test-agent" {
		t.Errorf("Expected UserAgent 'test-agent', got '%s'", options.UserAgent)
	}

	if options.Headers["X-Test"] != "value" {
		t.Errorf("Expected header X-Test=value, got '%s'", options.Headers["X-Test"])
	}
}

func TestDownloadResult(t *testing.T) {
	result := &DownloadResult{
		FilePath: "/tmp/test.zip",
		Size:     1024,
		Duration: 5 * time.Second,
		Success:  true,
	}

	if result.FilePath != "/tmp/test.zip" {
		t.Errorf("Expected FilePath '/tmp/test.zip', got '%s'", result.FilePath)
	}

	if result.Size != 1024 {
		t.Errorf("Expected Size 1024, got %d", result.Size)
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}
}

func TestProgressReader(t *testing.T) {
	data := "test data for progress reader"
	reader := strings.NewReader(data)

	var progressCalls []struct {
		downloaded, total int64
	}

	callback := func(downloaded, total int64) {
		progressCalls = append(progressCalls, struct {
			downloaded, total int64
		}{downloaded, total})
	}

	progressReader := &progressReader{
		reader:   reader,
		total:    int64(len(data)),
		callback: callback,
	}

	// 读取数据
	buffer := make([]byte, 10)
	n, err := progressReader.Read(buffer)
	if err != nil {
		t.Fatalf("progressReader.Read() failed: %v", err)
	}

	if n != 10 {
		t.Errorf("Expected to read 10 bytes, got %d", n)
	}

	if len(progressCalls) != 1 {
		t.Errorf("Expected 1 progress callback, got %d", len(progressCalls))
	}

	if progressCalls[0].downloaded != 10 {
		t.Errorf("Expected downloaded 10, got %d", progressCalls[0].downloaded)
	}

	if progressCalls[0].total != int64(len(data)) {
		t.Errorf("Expected total %d, got %d", len(data), progressCalls[0].total)
	}
}

func TestDownloadWithMockServer(t *testing.T) {
	// 创建模拟HTTP服务器
	testData := "test file content for download"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查User-Agent
		if userAgent := r.Header.Get("User-Agent"); userAgent != "test-downloader" {
			t.Errorf("Expected User-Agent 'test-downloader', got '%s'", userAgent)
		}

		// 检查自定义头部
		if testHeader := r.Header.Get("X-Test-Header"); testHeader != "test-value" {
			t.Errorf("Expected X-Test-Header 'test-value', got '%s'", testHeader)
		}

		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(testData)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	}))
	defer server.Close()

	downloader := NewDownloader()
	ctx := context.Background()

	// 创建临时目录
	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "test-file.txt")

	var progressCalls int
	options := DownloadOptions{
		URL:         server.URL,
		Destination: destPath,
		UserAgent:   "test-downloader",
		Headers:     map[string]string{"X-Test-Header": "test-value"},
		Progress: func(downloaded, total int64) {
			progressCalls++
			t.Logf("Progress: %d/%d bytes", downloaded, total)
		},
	}

	result, err := downloader.Download(ctx, options)
	if err != nil {
		t.Fatalf("Download() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected download to succeed")
	}

	if result.FilePath != destPath {
		t.Errorf("Expected FilePath '%s', got '%s'", destPath, result.FilePath)
	}

	if result.Size != int64(len(testData)) {
		t.Errorf("Expected Size %d, got %d", len(testData), result.Size)
	}

	// 验证文件内容
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != testData {
		t.Errorf("Expected file content '%s', got '%s'", testData, string(content))
	}

	// 验证进度回调被调用
	if progressCalls == 0 {
		t.Error("Expected progress callback to be called")
	}
}

func TestDownloadWithTimeout(t *testing.T) {
	// 创建一个慢响应的服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 延迟2秒
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("slow response"))
	}))
	defer server.Close()

	downloader := NewDownloader()
	ctx := context.Background()

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "timeout-test.txt")

	options := DownloadOptions{
		URL:         server.URL,
		Destination: destPath,
		Timeout:     100 * time.Millisecond, // 很短的超时
	}

	result, err := downloader.Download(ctx, options)

	// 应该超时
	if err == nil {
		t.Error("Expected timeout error")
	}

	if result != nil && result.Success {
		t.Error("Expected download to fail due to timeout")
	}

	// 验证文件没有被创建或已被清理
	if _, err := os.Stat(destPath); err == nil {
		t.Error("Expected incomplete file to be cleaned up")
	}
}

func TestDownloadWithRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// 前两次请求失败
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 第三次请求成功
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success after retry"))
	}))
	defer server.Close()

	downloader := NewDownloader()
	ctx := context.Background()

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "retry-test.txt")

	options := DownloadOptions{
		URL:         server.URL,
		Destination: destPath,
	}

	result, err := downloader.DownloadWithRetry(ctx, options, 3)
	if err != nil {
		t.Fatalf("DownloadWithRetry() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected download to succeed after retry")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	// 验证文件内容
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != "success after retry" {
		t.Errorf("Expected file content 'success after retry', got '%s'", string(content))
	}
}

func TestDownloadWithRetryFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 总是返回错误
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	downloader := NewDownloader()
	ctx := context.Background()

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "retry-fail-test.txt")

	options := DownloadOptions{
		URL:         server.URL,
		Destination: destPath,
	}

	result, err := downloader.DownloadWithRetry(ctx, options, 2)

	// 应该失败
	if err == nil {
		t.Error("Expected download to fail after retries")
	}

	if result != nil {
		t.Error("Expected nil result for failed download")
	}
}

func TestNewNodeJSDownloader(t *testing.T) {
	downloader := NewNodeJSDownloader()
	if downloader == nil {
		t.Fatal("NewNodeJSDownloader() returned nil")
	}

	if downloader.downloader == nil {
		t.Error("Expected downloader to be initialized")
	}

	if downloader.baseURL != "https://nodejs.org/dist" {
		t.Errorf("Expected baseURL 'https://nodejs.org/dist', got '%s'", downloader.baseURL)
	}
}

func TestGetDownloadURL(t *testing.T) {
	downloader := NewNodeJSDownloader()

	testCases := []struct {
		version  string
		platform Platform
		arch     Architecture
		expected string
	}{
		{
			version:  "18.17.0",
			platform: Windows,
			arch:     AMD64,
			expected: "https://nodejs.org/dist/v18.17.0/node-v18.17.0-win-x64.zip",
		},
		{
			version:  "18.17.0",
			platform: MacOS,
			arch:     ARM64,
			expected: "https://nodejs.org/dist/v18.17.0/node-v18.17.0-darwin-arm64.tar.gz",
		},
		{
			version:  "18.17.0",
			platform: Linux,
			arch:     AMD64,
			expected: "https://nodejs.org/dist/v18.17.0/node-v18.17.0-linux-x64.tar.xz",
		},
		{
			version:  "16.20.0",
			platform: Linux,
			arch:     ARM64,
			expected: "https://nodejs.org/dist/v16.20.0/node-v16.20.0-linux-arm64.tar.xz",
		},
		{
			version:  "16.20.0",
			platform: Linux,
			arch:     ARM,
			expected: "https://nodejs.org/dist/v16.20.0/node-v16.20.0-linux-armv7l.tar.xz",
		},
	}

	for _, tc := range testCases {
		url := downloader.GetDownloadURL(tc.version, tc.platform, tc.arch)
		if url != tc.expected {
			t.Errorf("GetDownloadURL(%s, %s, %s) = '%s', expected '%s'",
				tc.version, tc.platform, tc.arch, url, tc.expected)
		}
	}
}

func TestGetDownloadURLUnsupportedPlatform(t *testing.T) {
	downloader := NewNodeJSDownloader()

	url := downloader.GetDownloadURL("18.17.0", "unsupported", AMD64)
	if url != "" {
		t.Errorf("Expected empty URL for unsupported platform, got '%s'", url)
	}
}

func TestGetLatestVersion(t *testing.T) {
	downloader := NewNodeJSDownloader()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 这个测试需要网络连接，可能会失败
	version, err := downloader.GetLatestVersion(ctx)
	if err != nil {
		t.Logf("GetLatestVersion() failed (expected without network): %v", err)
		return
	}

	if version == "" {
		t.Error("Expected non-empty version")
	}

	t.Logf("Latest Node.js version: %s", version)
}
