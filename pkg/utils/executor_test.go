package utils

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}

	if executor.defaultTimeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", executor.defaultTimeout)
	}
}

func TestExecutorSetters(t *testing.T) {
	executor := NewExecutor()

	// Test SetDefaultTimeout
	timeout := 60 * time.Second
	executor.SetDefaultTimeout(timeout)
	if executor.defaultTimeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, executor.defaultTimeout)
	}

	// Test SetDefaultWorkingDir
	workDir := "/tmp"
	executor.SetDefaultWorkingDir(workDir)
	if executor.defaultWorkDir != workDir {
		t.Errorf("Expected working dir %s, got %s", workDir, executor.defaultWorkDir)
	}

	// Test SetDefaultEnv
	env := map[string]string{"TEST": "value"}
	executor.SetDefaultEnv(env)
	if executor.defaultEnv["TEST"] != "value" {
		t.Errorf("Expected env TEST=value, got %s", executor.defaultEnv["TEST"])
	}
}

func TestExecuteSimple(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test simple command that should work on all platforms
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "hello"}
	} else {
		cmd = "echo"
		args = []string{"hello"}
	}

	result, err := executor.ExecuteSimple(ctx, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteSimple() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("Expected stdout to contain 'hello', got '%s'", result.Stdout)
	}
}

func TestExecuteWithTimeout(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command with short timeout
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "timeout", "2"}
	} else {
		cmd = "sleep"
		args = []string{"2"}
	}

	result, err := executor.ExecuteWithTimeout(ctx, 100*time.Millisecond, cmd, args...)

	// Should timeout
	if err == nil {
		t.Error("Expected timeout error")
	}

	if result != nil && !result.Cancelled {
		t.Error("Expected command to be cancelled")
	}
}

func TestExecuteInDir(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command in specific directory
	var cmd string
	var args []string
	workDir := "/tmp"

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "cd"}
		workDir = "C:\\Windows"
	} else {
		cmd = "pwd"
		args = []string{}
		workDir = "/tmp"
	}

	result, err := executor.ExecuteInDir(ctx, workDir, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteInDir() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// Check if output contains the working directory
	if !strings.Contains(result.Stdout, workDir) {
		t.Errorf("Expected stdout to contain working dir '%s', got '%s'", workDir, result.Stdout)
	}
}

func TestExecuteWithInput(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command with input
	var cmd string
	var args []string
	input := "test input"

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "findstr", ".*"}
	} else {
		cmd = "cat"
		args = []string{}
	}

	result, err := executor.ExecuteWithInput(ctx, input, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteWithInput() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	if !strings.Contains(result.Stdout, input) {
		t.Errorf("Expected stdout to contain input '%s', got '%s'", input, result.Stdout)
	}
}

func TestExecuteStream(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test streaming output
	var outputs []string
	callback := func(output string) {
		outputs = append(outputs, output)
	}

	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "line1", "&", "echo", "line2"}
	} else {
		cmd = "sh"
		args = []string{"-c", "echo line1; echo line2"}
	}

	result, err := executor.ExecuteStream(ctx, callback, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteStream() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	if len(outputs) == 0 {
		t.Error("Expected to receive streaming output")
	}
}

func TestExecuteFailure(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command that should fail
	result, err := executor.ExecuteSimple(ctx, "nonexistent-command")

	if err == nil {
		t.Error("Expected error for nonexistent command")
	}

	if result != nil && result.Success {
		t.Error("Expected command to fail")
	}
}

func TestExecuteOptions(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test with custom options
	options := ExecuteOptions{
		Command:       "echo",
		Args:          []string{"test"},
		CaptureOutput: true,
		Timeout:       5 * time.Second,
		Env:           map[string]string{"TEST_VAR": "test_value"},
	}

	if runtime.GOOS == "windows" {
		options.Command = "cmd"
		options.Args = []string{"/c", "echo", "test"}
	}

	result, err := executor.Execute(ctx, options)
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	if result.Duration <= 0 {
		t.Error("Expected positive duration")
	}
}

func TestIsCommandAvailable(t *testing.T) {
	executor := NewExecutor()

	// Test with command that should exist
	var existingCmd string
	if runtime.GOOS == "windows" {
		existingCmd = "cmd"
	} else {
		existingCmd = "echo"
	}

	if !executor.IsCommandAvailable(existingCmd) {
		t.Errorf("Expected %s to be available", existingCmd)
	}

	// Test with command that should not exist
	if executor.IsCommandAvailable("definitely-nonexistent-command-12345") {
		t.Error("Expected nonexistent command to be unavailable")
	}
}

func TestGetCommandPath(t *testing.T) {
	executor := NewExecutor()

	// Test with command that should exist
	var existingCmd string
	if runtime.GOOS == "windows" {
		existingCmd = "cmd"
	} else {
		existingCmd = "echo"
	}

	path, err := executor.GetCommandPath(existingCmd)
	if err != nil {
		t.Errorf("GetCommandPath(%s) failed: %v", existingCmd, err)
	}

	if path == "" {
		t.Errorf("Expected non-empty path for %s", existingCmd)
	}

	// Test with command that should not exist
	_, err = executor.GetCommandPath("definitely-nonexistent-command-12345")
	if err == nil {
		t.Error("Expected error for nonexistent command")
	}
}

func TestExecuteEnvironmentVariables(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Set default environment
	executor.SetDefaultEnv(map[string]string{"DEFAULT_VAR": "default_value"})

	// Test command that prints environment variable
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "%TEST_VAR%"}
	} else {
		cmd = "sh"
		args = []string{"-c", "echo $TEST_VAR"}
	}

	options := ExecuteOptions{
		Command:       cmd,
		Args:          args,
		CaptureOutput: true,
		Env:           map[string]string{"TEST_VAR": "custom_value"},
	}

	result, err := executor.Execute(ctx, options)
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// On Unix systems, should contain the custom value
	if runtime.GOOS != "windows" && !strings.Contains(result.Stdout, "custom_value") {
		t.Errorf("Expected stdout to contain 'custom_value', got '%s'", result.Stdout)
	}
}

func TestExecuteCancellation(t *testing.T) {
	executor := NewExecutor()
	ctx, cancel := context.WithCancel(context.Background())

	// Start a long-running command
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "timeout", "10"}
	} else {
		cmd = "sleep"
		args = []string{"10"}
	}

	// Cancel after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	result, err := executor.ExecuteSimple(ctx, cmd, args...)

	// Should be cancelled
	if err == nil {
		t.Error("Expected cancellation error")
	}

	if result != nil && !result.Cancelled {
		t.Error("Expected command to be cancelled")
	}
}

func TestExecuteOptionsValidation(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test with empty command
	options := ExecuteOptions{
		Command: "",
		Args:    []string{"test"},
	}

	result, err := executor.Execute(ctx, options)
	if err == nil {
		t.Error("Expected error for empty command")
	}

	if result != nil && result.Success {
		t.Error("Expected command to fail")
	}

	// Test with invalid working directory
	options = ExecuteOptions{
		Command:    "echo",
		Args:       []string{"test"},
		WorkingDir: "/nonexistent/directory/path",
	}

	if runtime.GOOS == "windows" {
		options.Command = "cmd"
		options.Args = []string{"/c", "echo", "test"}
		options.WorkingDir = "Z:\\nonexistent\\directory\\path"
	}

	result, err = executor.Execute(ctx, options)
	if err == nil {
		t.Error("Expected error for invalid working directory")
	}

	if result != nil && result.Success {
		t.Error("Expected command to fail")
	}
}

func TestExecuteResultFields(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test successful command
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "success"}
	} else {
		cmd = "echo"
		args = []string{"success"}
	}

	result, err := executor.ExecuteSimple(ctx, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteSimple() failed: %v", err)
	}

	// Validate all result fields
	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected ExitCode 0, got %d", result.ExitCode)
	}

	if result.Duration <= 0 {
		t.Error("Expected positive Duration")
	}

	if result.Cancelled {
		t.Error("Expected Cancelled to be false")
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}

	if !strings.Contains(result.Stdout, "success") {
		t.Errorf("Expected stdout to contain 'success', got '%s'", result.Stdout)
	}

	// Test failed command
	result, err = executor.ExecuteSimple(ctx, "nonexistent-command-12345")
	if err == nil {
		t.Error("Expected error for nonexistent command")
	}

	if result == nil {
		t.Fatal("Expected result even for failed command")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}

	if result.Error == nil {
		t.Error("Expected Error to be set")
	}

	if result.Duration <= 0 {
		t.Error("Expected positive Duration even for failed command")
	}
}

func TestExecuteWithComplexInput(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test with multiline input
	input := "line1\nline2\nline3"
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "findstr", ".*"}
	} else {
		cmd = "cat"
		args = []string{}
	}

	result, err := executor.ExecuteWithInput(ctx, input, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteWithInput() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// Check that all lines are in output
	for _, line := range []string{"line1", "line2", "line3"} {
		if !strings.Contains(result.Stdout, line) {
			t.Errorf("Expected stdout to contain '%s', got '%s'", line, result.Stdout)
		}
	}

	// Test with empty input
	result, err = executor.ExecuteWithInput(ctx, "", cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteWithInput() with empty input failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed with empty input")
	}
}

func TestExecuteStreamingOutput(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test streaming with multiple lines
	var outputs []string
	var mu sync.Mutex
	callback := func(output string) {
		mu.Lock()
		outputs = append(outputs, output)
		mu.Unlock()
	}

	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "line1", "&", "echo", "line2", "&", "echo", "line3"}
	} else {
		cmd = "sh"
		args = []string{"-c", "echo line1; echo line2; echo line3"}
	}

	result, err := executor.ExecuteStream(ctx, callback, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteStream() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// Wait a bit for streaming to complete
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	outputCount := len(outputs)
	mu.Unlock()

	if outputCount == 0 {
		t.Error("Expected to receive streaming output")
	}

	// Test streaming with nil callback
	result, err = executor.ExecuteStream(ctx, nil, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteStream() with nil callback failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed with nil callback")
	}
}

func TestExecuteWithStderr(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command that outputs to stderr
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "error", "1>&2"}
	} else {
		cmd = "sh"
		args = []string{"-c", "echo error >&2"}
	}

	options := ExecuteOptions{
		Command:       cmd,
		Args:          args,
		CaptureOutput: true,
	}

	result, err := executor.Execute(ctx, options)
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// On Unix systems, stderr should contain the error message
	if runtime.GOOS != "windows" && !strings.Contains(result.Stderr, "error") {
		t.Errorf("Expected stderr to contain 'error', got '%s'", result.Stderr)
	}
}

func TestExecuteWithLargeOutput(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command that generates large output
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "for /L %i in (1,1,100) do echo Line %i"}
	} else {
		cmd = "sh"
		args = []string{"-c", "for i in $(seq 1 100); do echo Line $i; done"}
	}

	result, err := executor.ExecuteSimple(ctx, cmd, args...)
	if err != nil {
		t.Fatalf("ExecuteSimple() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// Check that output contains multiple lines
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	if len(lines) < 50 { // Should have many lines
		t.Errorf("Expected many lines of output, got %d", len(lines))
	}
}

func TestExecuteDefaultValues(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Set custom defaults
	executor.SetDefaultTimeout(10 * time.Second)
	executor.SetDefaultWorkingDir("/tmp")
	executor.SetDefaultEnv(map[string]string{"DEFAULT_TEST": "value"})

	// Test that defaults are used when not specified
	options := ExecuteOptions{
		Command:       "echo",
		Args:          []string{"test"},
		CaptureOutput: true,
		// Timeout, WorkingDir, and Env not specified - should use defaults
	}

	if runtime.GOOS == "windows" {
		options.Command = "cmd"
		options.Args = []string{"/c", "echo", "test"}
		executor.SetDefaultWorkingDir("C:\\Windows")
	}

	result, err := executor.Execute(ctx, options)
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// Test zero timeout uses default
	options.Timeout = 0
	result, err = executor.Execute(ctx, options)
	if err != nil {
		t.Fatalf("Execute() with zero timeout failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed with default timeout")
	}
}

func TestExecuteExitCodes(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command with specific exit code
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "exit", "42"}
	} else {
		cmd = "sh"
		args = []string{"-c", "exit 42"}
	}

	result, err := executor.ExecuteSimple(ctx, cmd, args...)
	if err == nil {
		t.Error("Expected error for non-zero exit code")
	}

	if result == nil {
		t.Fatal("Expected result even for failed command")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}

	if result.ExitCode != 42 {
		t.Errorf("Expected ExitCode 42, got %d", result.ExitCode)
	}
}

func TestExecuteWithoutCapture(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	// Test command without capturing output
	options := ExecuteOptions{
		Command:       "echo",
		Args:          []string{"test"},
		CaptureOutput: false,
		StreamOutput:  false,
	}

	if runtime.GOOS == "windows" {
		options.Command = "cmd"
		options.Args = []string{"/c", "echo", "test"}
	}

	result, err := executor.Execute(ctx, options)
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected command to succeed")
	}

	// Output should be empty since we didn't capture it
	if result.Stdout != "" {
		t.Errorf("Expected empty stdout, got '%s'", result.Stdout)
	}

	if result.Stderr != "" {
		t.Errorf("Expected empty stderr, got '%s'", result.Stderr)
	}
}

func TestKillProcess(t *testing.T) {
	executor := NewExecutor()

	// Test with nil command
	err := executor.KillProcess(nil)
	if err == nil {
		t.Error("Expected error for nil command")
	}

	// Test with command without process
	cmd := &exec.Cmd{}
	err = executor.KillProcess(cmd)
	if err == nil {
		t.Error("Expected error for command without process")
	}
}

func TestErrorConstants(t *testing.T) {
	// Test predefined errors
	if ErrCommandTimeout == nil {
		t.Error("Expected ErrCommandTimeout to be defined")
	}

	if ErrCommandFailed == nil {
		t.Error("Expected ErrCommandFailed to be defined")
	}

	if ErrInvalidCommand == nil {
		t.Error("Expected ErrInvalidCommand to be defined")
	}

	// Test error messages
	if ErrCommandTimeout.Error() != "command execution timeout" {
		t.Errorf("Expected timeout error message, got '%s'", ErrCommandTimeout.Error())
	}
}

func TestNewBatchExecutor(t *testing.T) {
	// Test with positive concurrency
	batchExecutor := NewBatchExecutor(5)
	if batchExecutor == nil {
		t.Fatal("NewBatchExecutor() returned nil")
	}

	if batchExecutor.maxConcurrency != 5 {
		t.Errorf("Expected maxConcurrency 5, got %d", batchExecutor.maxConcurrency)
	}

	if batchExecutor.executor == nil {
		t.Error("Expected executor to be initialized")
	}

	// Test with zero concurrency (should default to 1)
	batchExecutor = NewBatchExecutor(0)
	if batchExecutor.maxConcurrency != 1 {
		t.Errorf("Expected maxConcurrency 1 for zero input, got %d", batchExecutor.maxConcurrency)
	}

	// Test with negative concurrency (should default to 1)
	batchExecutor = NewBatchExecutor(-5)
	if batchExecutor.maxConcurrency != 1 {
		t.Errorf("Expected maxConcurrency 1 for negative input, got %d", batchExecutor.maxConcurrency)
	}
}

func TestBatchExecutorSuccess(t *testing.T) {
	batchExecutor := NewBatchExecutor(2)
	ctx := context.Background()

	// Create multiple successful commands
	var commands []ExecuteOptions
	for i := 0; i < 3; i++ {
		cmd := ExecuteOptions{
			Command:       "echo",
			Args:          []string{fmt.Sprintf("test%d", i)},
			CaptureOutput: true,
		}

		if runtime.GOOS == "windows" {
			cmd.Command = "cmd"
			cmd.Args = []string{"/c", "echo", fmt.Sprintf("test%d", i)}
		}

		commands = append(commands, cmd)
	}

	options := BatchOptions{
		Commands:       commands,
		StopOnError:    false,
		MaxConcurrency: 2,
	}

	result, err := batchExecutor.ExecuteBatch(ctx, options)
	if err != nil {
		t.Fatalf("ExecuteBatch() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected batch execution to succeed")
	}

	if len(result.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Results))
	}

	if result.FailedCount != 0 {
		t.Errorf("Expected 0 failed commands, got %d", result.FailedCount)
	}

	if result.TotalTime <= 0 {
		t.Error("Expected positive total time")
	}

	// Check individual results
	for i, res := range result.Results {
		if res == nil {
			t.Errorf("Expected result %d to be non-nil", i)
			continue
		}

		if !res.Success {
			t.Errorf("Expected result %d to succeed", i)
		}

		expectedOutput := fmt.Sprintf("test%d", i)
		if !strings.Contains(res.Stdout, expectedOutput) {
			t.Errorf("Expected result %d stdout to contain '%s', got '%s'", i, expectedOutput, res.Stdout)
		}
	}
}

func TestBatchExecutorWithFailures(t *testing.T) {
	batchExecutor := NewBatchExecutor(2)
	ctx := context.Background()

	// Create commands with some failures
	commands := []ExecuteOptions{
		{
			Command:       "echo",
			Args:          []string{"success1"},
			CaptureOutput: true,
		},
		{
			Command:       "nonexistent-command-12345",
			Args:          []string{},
			CaptureOutput: true,
		},
		{
			Command:       "echo",
			Args:          []string{"success2"},
			CaptureOutput: true,
		},
	}

	if runtime.GOOS == "windows" {
		commands[0].Command = "cmd"
		commands[0].Args = []string{"/c", "echo", "success1"}
		commands[2].Command = "cmd"
		commands[2].Args = []string{"/c", "echo", "success2"}
	}

	options := BatchOptions{
		Commands:       commands,
		StopOnError:    false,
		MaxConcurrency: 2,
	}

	result, err := batchExecutor.ExecuteBatch(ctx, options)
	if err != nil {
		t.Fatalf("ExecuteBatch() failed: %v", err)
	}

	if result.Success {
		t.Error("Expected batch execution to fail due to failed command")
	}

	if len(result.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Results))
	}

	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed command, got %d", result.FailedCount)
	}

	// Check that successful commands still succeeded
	if result.Results[0] != nil && !result.Results[0].Success {
		t.Error("Expected first command to succeed")
	}

	if result.Results[2] != nil && !result.Results[2].Success {
		t.Error("Expected third command to succeed")
	}

	// Check that failed command failed
	if result.Results[1] != nil && result.Results[1].Success {
		t.Error("Expected second command to fail")
	}
}

func TestBatchExecutorStopOnError(t *testing.T) {
	batchExecutor := NewBatchExecutor(1) // Use concurrency 1 for predictable order
	ctx := context.Background()

	// Create commands with early failure
	commands := []ExecuteOptions{
		{
			Command:       "nonexistent-command-12345",
			Args:          []string{},
			CaptureOutput: true,
		},
		{
			Command:       "echo",
			Args:          []string{"should-not-run"},
			CaptureOutput: true,
		},
	}

	if runtime.GOOS == "windows" {
		commands[1].Command = "cmd"
		commands[1].Args = []string{"/c", "echo", "should-not-run"}
	}

	options := BatchOptions{
		Commands:       commands,
		StopOnError:    true,
		MaxConcurrency: 1,
	}

	result, err := batchExecutor.ExecuteBatch(ctx, options)
	if err != nil {
		t.Fatalf("ExecuteBatch() failed: %v", err)
	}

	if result.Success {
		t.Error("Expected batch execution to fail")
	}

	if result.FailedCount == 0 {
		t.Error("Expected at least one failed command")
	}

	// First command should fail
	if result.Results[0] != nil && result.Results[0].Success {
		t.Error("Expected first command to fail")
	}
}

func TestBatchExecutorEmptyCommands(t *testing.T) {
	batchExecutor := NewBatchExecutor(2)
	ctx := context.Background()

	// Test with empty command list
	options := BatchOptions{
		Commands:       []ExecuteOptions{},
		StopOnError:    false,
		MaxConcurrency: 2,
	}

	result, err := batchExecutor.ExecuteBatch(ctx, options)
	if err != nil {
		t.Fatalf("ExecuteBatch() with empty commands failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected empty batch to succeed")
	}

	if len(result.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(result.Results))
	}

	if result.FailedCount != 0 {
		t.Errorf("Expected 0 failed commands, got %d", result.FailedCount)
	}
}

func TestBatchExecutorConcurrencyLimits(t *testing.T) {
	batchExecutor := NewBatchExecutor(1) // Max concurrency 1
	ctx := context.Background()

	// Create multiple commands that take some time
	var commands []ExecuteOptions
	for i := 0; i < 3; i++ {
		cmd := ExecuteOptions{
			Command:       "echo",
			Args:          []string{fmt.Sprintf("test%d", i)},
			CaptureOutput: true,
		}

		if runtime.GOOS == "windows" {
			cmd.Command = "cmd"
			cmd.Args = []string{"/c", "echo", fmt.Sprintf("test%d", i)}
		}

		commands = append(commands, cmd)
	}

	options := BatchOptions{
		Commands:       commands,
		StopOnError:    false,
		MaxConcurrency: 0, // Should use executor's default (1)
	}

	startTime := time.Now()
	result, err := batchExecutor.ExecuteBatch(ctx, options)
	duration := time.Since(startTime)

	if err != nil {
		t.Fatalf("ExecuteBatch() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected batch execution to succeed")
	}

	// With concurrency 1, commands should run sequentially
	// This is hard to test precisely, but we can check basic functionality
	if len(result.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Results))
	}

	t.Logf("Batch execution took %v", duration)
}

func TestBatchExecutorCancellation(t *testing.T) {
	batchExecutor := NewBatchExecutor(2)
	ctx, cancel := context.WithCancel(context.Background())

	// Create long-running commands
	var commands []ExecuteOptions
	for i := 0; i < 3; i++ {
		cmd := ExecuteOptions{
			Command:       "sleep",
			Args:          []string{"5"},
			CaptureOutput: true,
		}

		if runtime.GOOS == "windows" {
			cmd.Command = "cmd"
			cmd.Args = []string{"/c", "timeout", "5"}
		}

		commands = append(commands, cmd)
	}

	options := BatchOptions{
		Commands:       commands,
		StopOnError:    false,
		MaxConcurrency: 2,
	}

	// Cancel after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	result, err := batchExecutor.ExecuteBatch(ctx, options)
	if err != nil {
		t.Fatalf("ExecuteBatch() failed: %v", err)
	}

	// Some commands might be cancelled
	if result.Success {
		t.Log("Batch completed before cancellation (possible)")
	} else {
		t.Log("Batch was cancelled as expected")
	}

	if len(result.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Results))
	}
}

func TestBatchOptionsValidation(t *testing.T) {
	// Test BatchOptions struct
	options := BatchOptions{
		Commands: []ExecuteOptions{
			{
				Command:       "echo",
				Args:          []string{"test"},
				CaptureOutput: true,
			},
		},
		StopOnError:    true,
		MaxConcurrency: 5,
	}

	if len(options.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(options.Commands))
	}

	if !options.StopOnError {
		t.Error("Expected StopOnError to be true")
	}

	if options.MaxConcurrency != 5 {
		t.Errorf("Expected MaxConcurrency 5, got %d", options.MaxConcurrency)
	}
}

func TestBatchResultValidation(t *testing.T) {
	// Test BatchResult struct
	results := []*ExecuteResult{
		{
			Success:  true,
			ExitCode: 0,
			Duration: time.Second,
		},
		{
			Success:  false,
			ExitCode: 1,
			Duration: 2 * time.Second,
		},
	}

	batchResult := BatchResult{
		Results:     results,
		Success:     false,
		TotalTime:   5 * time.Second,
		FailedCount: 1,
	}

	if len(batchResult.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(batchResult.Results))
	}

	if batchResult.Success {
		t.Error("Expected Success to be false")
	}

	if batchResult.FailedCount != 1 {
		t.Errorf("Expected FailedCount 1, got %d", batchResult.FailedCount)
	}

	if batchResult.TotalTime != 5*time.Second {
		t.Errorf("Expected TotalTime 5s, got %v", batchResult.TotalTime)
	}
}
