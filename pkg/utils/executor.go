package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ExecuteOptions 执行选项
type ExecuteOptions struct {
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	Timeout     time.Duration     `json:"timeout"`
	Input       string            `json:"input"`
	CaptureOutput bool            `json:"capture_output"`
	StreamOutput  bool            `json:"stream_output"`
	OutputCallback func(string)   `json:"-"`
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	Success    bool          `json:"success"`
	ExitCode   int           `json:"exit_code"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	Duration   time.Duration `json:"duration"`
	Error      error         `json:"error,omitempty"`
	Cancelled  bool          `json:"cancelled"`
}

// Executor 命令执行器
type Executor struct {
	defaultTimeout time.Duration
	defaultWorkDir string
	defaultEnv     map[string]string
}

// NewExecutor 创建新的执行器
func NewExecutor() *Executor {
	return &Executor{
		defaultTimeout: 30 * time.Second,
		defaultWorkDir: "",
		defaultEnv:     make(map[string]string),
	}
}

// SetDefaultTimeout 设置默认超时时间
func (e *Executor) SetDefaultTimeout(timeout time.Duration) {
	e.defaultTimeout = timeout
}

// SetDefaultWorkingDir 设置默认工作目录
func (e *Executor) SetDefaultWorkingDir(dir string) {
	e.defaultWorkDir = dir
}

// SetDefaultEnv 设置默认环境变量
func (e *Executor) SetDefaultEnv(env map[string]string) {
	e.defaultEnv = env
}

// Execute 执行命令
func (e *Executor) Execute(ctx context.Context, options ExecuteOptions) (*ExecuteResult, error) {
	startTime := time.Now()
	
	// 设置默认值
	if options.Timeout == 0 {
		options.Timeout = e.defaultTimeout
	}
	if options.WorkingDir == "" {
		options.WorkingDir = e.defaultWorkDir
	}

	// 创建带超时的上下文
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	// 创建命令
	cmd := exec.CommandContext(ctx, options.Command, options.Args...)
	
	// 设置工作目录
	if options.WorkingDir != "" {
		cmd.Dir = options.WorkingDir
	}

	// 设置环境变量
	cmd.Env = os.Environ()
	for key, value := range e.defaultEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	for key, value := range options.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// 设置输入
	if options.Input != "" {
		cmd.Stdin = strings.NewReader(options.Input)
	}

	var stdout, stderr strings.Builder
	var wg sync.WaitGroup

	// 处理输出
	if options.CaptureOutput || options.StreamOutput {
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return &ExecuteResult{
				Success:  false,
				Duration: time.Since(startTime),
				Error:    fmt.Errorf("failed to create stdout pipe: %w", err),
			}, err
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return &ExecuteResult{
				Success:  false,
				Duration: time.Since(startTime),
				Error:    fmt.Errorf("failed to create stderr pipe: %w", err),
			}, err
		}

		// 处理stdout
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.handleOutput(stdoutPipe, &stdout, "stdout", options.StreamOutput, options.OutputCallback)
		}()

		// 处理stderr
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.handleOutput(stderrPipe, &stderr, "stderr", options.StreamOutput, options.OutputCallback)
		}()
	}

	// 启动命令
	err := cmd.Start()
	if err != nil {
		return &ExecuteResult{
			Success:  false,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("failed to start command: %w", err),
		}, err
	}

	// 等待命令完成
	err = cmd.Wait()
	
	// 等待输出处理完成
	wg.Wait()

	// 构建结果
	result := &ExecuteResult{
		Duration: time.Since(startTime),
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}

	// 检查是否被取消
	if ctx.Err() == context.DeadlineExceeded {
		result.Cancelled = true
		result.Error = ErrCommandTimeout
		return result, ErrCommandTimeout
	} else if ctx.Err() == context.Canceled {
		result.Cancelled = true
		result.Error = context.Canceled
		return result, context.Canceled
	}

	// 处理命令执行结果
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
			result.Error = fmt.Errorf("command failed with exit code %d", result.ExitCode)
		} else {
			result.Error = fmt.Errorf("command execution failed: %w", err)
		}
		return result, result.Error
	}

	result.Success = true
	result.ExitCode = 0
	return result, nil
}

// ExecuteSimple 简单执行命令
func (e *Executor) ExecuteSimple(ctx context.Context, command string, args ...string) (*ExecuteResult, error) {
	return e.Execute(ctx, ExecuteOptions{
		Command:       command,
		Args:          args,
		CaptureOutput: true,
	})
}

// ExecuteWithTimeout 带超时执行命令
func (e *Executor) ExecuteWithTimeout(ctx context.Context, timeout time.Duration, command string, args ...string) (*ExecuteResult, error) {
	return e.Execute(ctx, ExecuteOptions{
		Command:       command,
		Args:          args,
		Timeout:       timeout,
		CaptureOutput: true,
	})
}

// ExecuteInDir 在指定目录执行命令
func (e *Executor) ExecuteInDir(ctx context.Context, workingDir, command string, args ...string) (*ExecuteResult, error) {
	return e.Execute(ctx, ExecuteOptions{
		Command:       command,
		Args:          args,
		WorkingDir:    workingDir,
		CaptureOutput: true,
	})
}

// ExecuteWithInput 带输入执行命令
func (e *Executor) ExecuteWithInput(ctx context.Context, input, command string, args ...string) (*ExecuteResult, error) {
	return e.Execute(ctx, ExecuteOptions{
		Command:       command,
		Args:          args,
		Input:         input,
		CaptureOutput: true,
	})
}

// ExecuteStream 流式执行命令
func (e *Executor) ExecuteStream(ctx context.Context, outputCallback func(string), command string, args ...string) (*ExecuteResult, error) {
	return e.Execute(ctx, ExecuteOptions{
		Command:        command,
		Args:           args,
		StreamOutput:   true,
		OutputCallback: outputCallback,
	})
}

// handleOutput 处理输出流
func (e *Executor) handleOutput(pipe io.ReadCloser, builder *strings.Builder, streamType string, streamOutput bool, callback func(string)) {
	defer pipe.Close()
	
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		
		// 写入builder用于捕获
		if builder != nil {
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		
		// 流式输出
		if streamOutput && callback != nil {
			callback(fmt.Sprintf("[%s] %s", streamType, line))
		}
	}
}

// IsCommandAvailable 检查命令是否可用
func (e *Executor) IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GetCommandPath 获取命令路径
func (e *Executor) GetCommandPath(command string) (string, error) {
	return exec.LookPath(command)
}

// KillProcess 终止进程
func (e *Executor) KillProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("invalid process")
	}
	
	return cmd.Process.Kill()
}

// 预定义错误
var (
	ErrCommandTimeout = fmt.Errorf("command execution timeout")
	ErrCommandFailed  = fmt.Errorf("command execution failed")
	ErrInvalidCommand = fmt.Errorf("invalid command")
)

// BatchExecutor 批量执行器
type BatchExecutor struct {
	executor *Executor
	maxConcurrency int
}

// NewBatchExecutor 创建批量执行器
func NewBatchExecutor(maxConcurrency int) *BatchExecutor {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}
	
	return &BatchExecutor{
		executor:       NewExecutor(),
		maxConcurrency: maxConcurrency,
	}
}

// BatchOptions 批量执行选项
type BatchOptions struct {
	Commands       []ExecuteOptions `json:"commands"`
	StopOnError    bool             `json:"stop_on_error"`
	MaxConcurrency int              `json:"max_concurrency"`
}

// BatchResult 批量执行结果
type BatchResult struct {
	Results    []*ExecuteResult `json:"results"`
	Success    bool             `json:"success"`
	TotalTime  time.Duration    `json:"total_time"`
	FailedCount int             `json:"failed_count"`
}

// ExecuteBatch 批量执行命令
func (be *BatchExecutor) ExecuteBatch(ctx context.Context, options BatchOptions) (*BatchResult, error) {
	startTime := time.Now()
	
	concurrency := options.MaxConcurrency
	if concurrency <= 0 {
		concurrency = be.maxConcurrency
	}
	
	results := make([]*ExecuteResult, len(options.Commands))
	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failedCount int
	var shouldStop bool

	for i, cmd := range options.Commands {
		wg.Add(1)
		go func(index int, command ExecuteOptions) {
			defer wg.Done()
			
			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			// 检查是否应该停止
			mu.Lock()
			if shouldStop {
				mu.Unlock()
				results[index] = &ExecuteResult{
					Success:   false,
					Cancelled: true,
					Error:     fmt.Errorf("execution stopped due to previous error"),
				}
				return
			}
			mu.Unlock()
			
			// 执行命令
			result, _ := be.executor.Execute(ctx, command)
			results[index] = result
			
			// 检查是否失败
			if !result.Success {
				mu.Lock()
				failedCount++
				if options.StopOnError {
					shouldStop = true
				}
				mu.Unlock()
			}
		}(i, cmd)
	}
	
	wg.Wait()
	
	return &BatchResult{
		Results:     results,
		Success:     failedCount == 0,
		TotalTime:   time.Since(startTime),
		FailedCount: failedCount,
	}, nil
}
