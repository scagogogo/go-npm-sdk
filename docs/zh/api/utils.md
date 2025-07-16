# 工具包API

`pkg/utils`包提供命令执行的工具函数，具有超时控制、流式输出和批量操作等高级功能。

## 执行器

执行器组件提供灵活的命令执行功能，支持各种执行模式和高级功能。

### NewExecutor

```go
func NewExecutor() *Executor
```

使用默认配置创建新的命令执行器。

**示例:**
```go
executor := utils.NewExecutor()
```

## 基本执行

### ExecuteSimple

```go
func (e *Executor) ExecuteSimple(ctx context.Context, command string, args ...string) (*ExecuteResult, error)
```

使用基本配置执行命令并返回结果。

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `command` (string): 要执行的命令
- `args` (...string): 命令参数

**返回:**
- `*ExecuteResult`: 执行结果
- `error`: 如果执行失败返回错误

**示例:**
```go
ctx := context.Background()
result, err := executor.ExecuteSimple(ctx, "npm", "--version")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("npm版本: %s\n", strings.TrimSpace(result.Stdout))
fmt.Printf("退出代码: %d\n", result.ExitCode)
fmt.Printf("持续时间: %v\n", result.Duration)
```

### ExecuteWithInput

```go
func (e *Executor) ExecuteWithInput(ctx context.Context, input string, command string, args ...string) (*ExecuteResult, error)
```

执行命令并向stdin提供输入数据。

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `input` (string): 发送到命令stdin的输入数据
- `command` (string): 要执行的命令
- `args` (...string): 命令参数

**示例:**
```go
ctx := context.Background()
input := "console.log('Hello, World!');"
result, err := executor.ExecuteWithInput(ctx, input, "node", "-e", "-")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("输出: %s\n", result.Stdout)
```

## 高级执行

### Execute

```go
func (e *Executor) Execute(ctx context.Context, options ExecuteOptions) (*ExecuteResult, error)
```

使用高级配置选项执行命令。

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `options` (ExecuteOptions): 执行配置

**返回:**
- `*ExecuteResult`: 执行结果
- `error`: 如果执行失败返回错误

### ExecuteOptions

```go
type ExecuteOptions struct {
    Command       string            `json:"command"`
    Args          []string          `json:"args,omitempty"`
    WorkingDir    string            `json:"working_dir,omitempty"`
    Env           map[string]string `json:"env,omitempty"`
    Timeout       time.Duration     `json:"timeout,omitempty"`
    Input         string            `json:"input,omitempty"`
    CaptureOutput bool              `json:"capture_output"`
    StreamOutput  bool              `json:"stream_output"`
    OutputCallback func(string)     `json:"-"`
}
```

**示例:**
```go
ctx := context.Background()
options := utils.ExecuteOptions{
    Command:       "npm",
    Args:          []string{"install", "lodash"},
    WorkingDir:    "/path/to/project",
    Timeout:       5 * time.Minute,
    CaptureOutput: true,
    StreamOutput:  true,
    Env: map[string]string{
        "NODE_ENV": "development",
    },
    OutputCallback: func(output string) {
        fmt.Printf("npm: %s", output)
    },
}

result, err := executor.Execute(ctx, options)
if err != nil {
    log.Fatal(err)
}

if result.Success {
    fmt.Println("包安装成功")
} else {
    fmt.Printf("安装失败: %s\n", result.Stderr)
}
```

## 流式执行

### ExecuteStream

```go
func (e *Executor) ExecuteStream(ctx context.Context, callback func(string), command string, args ...string) (*ExecuteResult, error)
```

执行命令并实时流式输出。

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `callback` (func(string)): 为每行输出调用的函数
- `command` (string): 要执行的命令
- `args` (...string): 命令参数

**示例:**
```go
ctx := context.Background()

callback := func(output string) {
    fmt.Printf("实时输出: %s", output)
}

result, err := executor.ExecuteStream(ctx, callback, "npm", "install", "--verbose")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("命令完成，退出代码: %d\n", result.ExitCode)
```

## 执行结果

### ExecuteResult

```go
type ExecuteResult struct {
    Success   bool          `json:"success"`
    ExitCode  int           `json:"exit_code"`
    Stdout    string        `json:"stdout,omitempty"`
    Stderr    string        `json:"stderr,omitempty"`
    Duration  time.Duration `json:"duration"`
    Cancelled bool          `json:"cancelled"`
    Error     error         `json:"-"`
}
```

**字段:**
- `Success`: 命令是否成功执行（退出代码0）
- `ExitCode`: 命令退出代码
- `Stdout`: 标准输出内容
- `Stderr`: 标准错误内容
- `Duration`: 执行持续时间
- `Cancelled`: 执行是否被取消
- `Error`: 执行错误（如果有）

## 批量执行

### NewBatchExecutor

```go
func NewBatchExecutor(maxConcurrency int) *BatchExecutor
```

使用指定的最大并发数创建新的批量执行器。

**参数:**
- `maxConcurrency` (int): 最大并发执行数

### ExecuteBatch

```go
func (be *BatchExecutor) ExecuteBatch(ctx context.Context, options BatchOptions) (*BatchResult, error)
```

使用指定选项并发执行多个命令。

**参数:**
- `ctx` (context.Context): 用于取消和超时的上下文
- `options` (BatchOptions): 批量执行配置

### BatchOptions

```go
type BatchOptions struct {
    Commands       []ExecuteOptions `json:"commands"`
    StopOnError    bool             `json:"stop_on_error"`
    MaxConcurrency int              `json:"max_concurrency,omitempty"`
}
```

### BatchResult

```go
type BatchResult struct {
    Results     []*ExecuteResult `json:"results"`
    Success     bool             `json:"success"`
    TotalTime   time.Duration    `json:"total_time"`
    FailedCount int              `json:"failed_count"`
}
```

**示例:**
```go
batchExecutor := utils.NewBatchExecutor(3)
ctx := context.Background()

commands := []utils.ExecuteOptions{
    {
        Command:       "npm",
        Args:          []string{"install", "lodash"},
        WorkingDir:    "/project1",
        CaptureOutput: true,
    },
    {
        Command:       "npm",
        Args:          []string{"install", "axios"},
        WorkingDir:    "/project2",
        CaptureOutput: true,
    },
    {
        Command:       "npm",
        Args:          []string{"test"},
        WorkingDir:    "/project3",
        CaptureOutput: true,
    },
}

options := utils.BatchOptions{
    Commands:       commands,
    StopOnError:    false,
    MaxConcurrency: 2,
}

result, err := batchExecutor.ExecuteBatch(ctx, options)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("批量执行在%v内完成\n", result.TotalTime)
fmt.Printf("成功: %v, 失败: %d/%d\n", result.Success, result.FailedCount, len(result.Results))

for i, res := range result.Results {
    if res.Success {
        fmt.Printf("命令%d: 成功\n", i+1)
    } else {
        fmt.Printf("命令%d: 失败 (退出代码: %d)\n", i+1, res.ExitCode)
    }
}
```

## 配置

### SetDefaultTimeout

```go
func (e *Executor) SetDefaultTimeout(timeout time.Duration)
```

设置命令执行的默认超时时间。

### SetDefaultWorkingDir

```go
func (e *Executor) SetDefaultWorkingDir(dir string)
```

设置命令执行的默认工作目录。

### SetDefaultEnv

```go
func (e *Executor) SetDefaultEnv(env map[string]string)
```

设置命令执行的默认环境变量。

**示例:**
```go
executor := utils.NewExecutor()

// 配置默认值
executor.SetDefaultTimeout(30 * time.Second)
executor.SetDefaultWorkingDir("/home/user/projects")
executor.SetDefaultEnv(map[string]string{
    "NODE_ENV": "production",
    "PATH":     "/usr/local/bin:/usr/bin:/bin",
})

// 当在ExecuteOptions中未指定时，将使用这些默认值
```

## 进程管理

### KillProcess

```go
func (e *Executor) KillProcess(cmd *exec.Cmd) error
```

终止正在运行的进程。

**参数:**
- `cmd` (*exec.Cmd): 要终止的命令进程

**示例:**
```go
// 这通常在内部使用，但如果需要可以手动调用
err := executor.KillProcess(cmd)
if err != nil {
    log.Printf("终止进程失败: %v", err)
}
```

## 错误处理

工具包提供特定的错误常量：

```go
var (
    ErrCommandTimeout  = errors.New("command execution timeout")
    ErrCommandFailed   = errors.New("command execution failed")
    ErrInvalidCommand  = errors.New("invalid command")
)
```

**示例:**
```go
result, err := executor.ExecuteSimple(ctx, "invalid-command")
if err != nil {
    if errors.Is(err, utils.ErrCommandTimeout) {
        fmt.Println("命令超时")
    } else if errors.Is(err, utils.ErrCommandFailed) {
        fmt.Println("命令失败")
    } else if errors.Is(err, utils.ErrInvalidCommand) {
        fmt.Println("无效命令")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 最佳实践

1. **使用context**: 始终传递适当的context以进行超时和取消控制
2. **处理超时**: 为长时间运行的命令设置合理的超时
3. **流式输出**: 对有大量输出的命令使用流式处理
4. **批量操作**: 对多个相关命令使用批量执行器
5. **错误处理**: 检查特定错误类型以便更好地处理错误
6. **工作目录**: 为项目特定命令设置适当的工作目录
7. **环境变量**: 根据需要配置环境变量
8. **资源清理**: 确保正确清理资源和进程

## 集成示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/scagogogo/go-npm-sdk/pkg/utils"
)

func main() {
    executor := utils.NewExecutor()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    // 配置执行器
    executor.SetDefaultTimeout(30 * time.Second)
    executor.SetDefaultWorkingDir("/path/to/project")
    
    // 使用流式输出执行npm install
    options := utils.ExecuteOptions{
        Command:       "npm",
        Args:          []string{"install"},
        CaptureOutput: true,
        StreamOutput:  true,
        OutputCallback: func(output string) {
            fmt.Printf("npm: %s", output)
        },
    }
    
    result, err := executor.Execute(ctx, options)
    if err != nil {
        log.Fatal(err)
    }
    
    if result.Success {
        fmt.Printf("安装在%v内完成\n", result.Duration)
    } else {
        fmt.Printf("安装失败: %s\n", result.Stderr)
    }
}
```
