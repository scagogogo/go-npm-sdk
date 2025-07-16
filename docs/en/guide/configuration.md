# Configuration

Learn how to configure the Go NPM SDK for your specific needs.

## Basic Configuration

### Client Configuration

```go
// Create client with default configuration
client, err := npm.NewClient()

// Create client with specific npm path
client, err := npm.NewClientWithPath("/usr/local/bin/npm")
```

### Working Directory

Set the working directory for npm operations:

```go
options := npm.InstallOptions{
    WorkingDir: "/path/to/project",
}

err := client.InstallPackage(ctx, "express", options)
```

### Registry Configuration

Configure custom npm registry:

```go
options := npm.InstallOptions{
    Registry: "https://registry.npmjs.org/",
}

err := client.InstallPackage(ctx, "package", options)
```

## Advanced Configuration

### Timeout Configuration

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := client.InstallPackage(ctx, "large-package", npm.InstallOptions{})
```

### Environment Variables

```go
import "github.com/scagogogo/go-npm-sdk/pkg/utils"

executor := utils.NewExecutor()
executor.SetDefaultEnv(map[string]string{
    "NODE_ENV": "production",
    "NPM_CONFIG_REGISTRY": "https://registry.npmjs.org/",
})
```

## Platform-Specific Configuration

### Windows Configuration

```go
// Windows-specific npm path
client, err := npm.NewClientWithPath("C:\\Program Files\\nodejs\\npm.cmd")
```

### macOS Configuration

```go
// macOS with Homebrew
client, err := npm.NewClientWithPath("/opt/homebrew/bin/npm")
```

### Linux Configuration

```go
// Linux system npm
client, err := npm.NewClientWithPath("/usr/bin/npm")
```

## Next Steps

- [Platform Support](./platform-support.md) - Platform-specific information
- [API Reference](/en/api/overview.md) - Complete API documentation
