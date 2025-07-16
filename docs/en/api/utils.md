# Utils Package API

The `pkg/utils` package provides utility functions for command execution with advanced features like timeout control, streaming output, and batch operations.

## Executor

The executor component provides flexible command execution capabilities with support for various execution modes and advanced features.

### NewExecutor

```go
func NewExecutor() *Executor
```

Creates a new command executor with default configuration.

**Example:**
```go
executor := utils.NewExecutor()
```

## Basic Execution

### ExecuteSimple

```go
func (e *Executor) ExecuteSimple(ctx context.Context, command string, args ...string) (*ExecuteResult, error)
```

Executes a command with basic configuration and returns the result.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `command` (string): Command to execute
- `args` (...string): Command arguments

**Returns:**
- `*ExecuteResult`: Execution result
- `error`: Error if execution fails

**Example:**
```go
ctx := context.Background()
result, err := executor.ExecuteSimple(ctx, "npm", "--version")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("npm version: %s\n", strings.TrimSpace(result.Stdout))
fmt.Printf("Exit code: %d\n", result.ExitCode)
fmt.Printf("Duration: %v\n", result.Duration)
```

### ExecuteWithInput

```go
func (e *Executor) ExecuteWithInput(ctx context.Context, input string, command string, args ...string) (*ExecuteResult, error)
```

Executes a command with input data provided to stdin.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `input` (string): Input data to send to command's stdin
- `command` (string): Command to execute
- `args` (...string): Command arguments

**Example:**
```go
ctx := context.Background()
input := "console.log('Hello, World!');"
result, err := executor.ExecuteWithInput(ctx, input, "node", "-e", "-")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Output: %s\n", result.Stdout)
```

## Advanced Execution

### Execute

```go
func (e *Executor) Execute(ctx context.Context, options ExecuteOptions) (*ExecuteResult, error)
```

Executes a command with advanced configuration options.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `options` (ExecuteOptions): Execution configuration

**Returns:**
- `*ExecuteResult`: Execution result
- `error`: Error if execution fails

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

**Example:**
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
    fmt.Println("Package installed successfully")
} else {
    fmt.Printf("Installation failed: %s\n", result.Stderr)
}
```

## Streaming Execution

### ExecuteStream

```go
func (e *Executor) ExecuteStream(ctx context.Context, callback func(string), command string, args ...string) (*ExecuteResult, error)
```

Executes a command with real-time output streaming.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `callback` (func(string)): Function called for each line of output
- `command` (string): Command to execute
- `args` (...string): Command arguments

**Example:**
```go
ctx := context.Background()

callback := func(output string) {
    fmt.Printf("Real-time output: %s", output)
}

result, err := executor.ExecuteStream(ctx, callback, "npm", "install", "--verbose")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Command completed with exit code: %d\n", result.ExitCode)
```

## Execution Results

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

**Fields:**
- `Success`: Whether the command executed successfully (exit code 0)
- `ExitCode`: Command exit code
- `Stdout`: Standard output content
- `Stderr`: Standard error content
- `Duration`: Execution duration
- `Cancelled`: Whether execution was cancelled
- `Error`: Execution error (if any)

## Batch Execution

### NewBatchExecutor

```go
func NewBatchExecutor(maxConcurrency int) *BatchExecutor
```

Creates a new batch executor with specified maximum concurrency.

**Parameters:**
- `maxConcurrency` (int): Maximum number of concurrent executions

### ExecuteBatch

```go
func (be *BatchExecutor) ExecuteBatch(ctx context.Context, options BatchOptions) (*BatchResult, error)
```

Executes multiple commands concurrently with specified options.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeout
- `options` (BatchOptions): Batch execution configuration

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

**Example:**
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

fmt.Printf("Batch execution completed in %v\n", result.TotalTime)
fmt.Printf("Success: %v, Failed: %d/%d\n", result.Success, result.FailedCount, len(result.Results))

for i, res := range result.Results {
    if res.Success {
        fmt.Printf("Command %d: SUCCESS\n", i+1)
    } else {
        fmt.Printf("Command %d: FAILED (exit code: %d)\n", i+1, res.ExitCode)
    }
}
```

## Configuration

### SetDefaultTimeout

```go
func (e *Executor) SetDefaultTimeout(timeout time.Duration)
```

Sets the default timeout for command execution.

### SetDefaultWorkingDir

```go
func (e *Executor) SetDefaultWorkingDir(dir string)
```

Sets the default working directory for command execution.

### SetDefaultEnv

```go
func (e *Executor) SetDefaultEnv(env map[string]string)
```

Sets the default environment variables for command execution.

**Example:**
```go
executor := utils.NewExecutor()

// Configure defaults
executor.SetDefaultTimeout(30 * time.Second)
executor.SetDefaultWorkingDir("/home/user/projects")
executor.SetDefaultEnv(map[string]string{
    "NODE_ENV": "production",
    "PATH":     "/usr/local/bin:/usr/bin:/bin",
})

// These defaults will be used when not specified in ExecuteOptions
```

## Process Management

### KillProcess

```go
func (e *Executor) KillProcess(cmd *exec.Cmd) error
```

Kills a running process.

**Parameters:**
- `cmd` (*exec.Cmd): Command process to kill

**Example:**
```go
// This is typically used internally, but can be called manually if needed
err := executor.KillProcess(cmd)
if err != nil {
    log.Printf("Failed to kill process: %v", err)
}
```

## Error Handling

The utils package provides specific error constants:

```go
var (
    ErrCommandTimeout  = errors.New("command execution timeout")
    ErrCommandFailed   = errors.New("command execution failed")
    ErrInvalidCommand  = errors.New("invalid command")
)
```

**Example:**
```go
result, err := executor.ExecuteSimple(ctx, "invalid-command")
if err != nil {
    if errors.Is(err, utils.ErrCommandTimeout) {
        fmt.Println("Command timed out")
    } else if errors.Is(err, utils.ErrCommandFailed) {
        fmt.Println("Command failed")
    } else if errors.Is(err, utils.ErrInvalidCommand) {
        fmt.Println("Invalid command")
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Best Practices

1. **Use contexts**: Always pass appropriate contexts for timeout and cancellation control
2. **Handle timeouts**: Set reasonable timeouts for long-running commands
3. **Stream output**: Use streaming for commands with significant output
4. **Batch operations**: Use batch executor for multiple related commands
5. **Error handling**: Check specific error types for better error handling
6. **Working directories**: Set appropriate working directories for project-specific commands
7. **Environment variables**: Configure environment variables as needed
8. **Resource cleanup**: Ensure proper cleanup of resources and processes

## Integration Example

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
    
    // Configure executor
    executor.SetDefaultTimeout(30 * time.Second)
    executor.SetDefaultWorkingDir("/path/to/project")
    
    // Execute npm install with streaming output
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
        fmt.Printf("Installation completed in %v\n", result.Duration)
    } else {
        fmt.Printf("Installation failed: %s\n", result.Stderr)
    }
}
```
