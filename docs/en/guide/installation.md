# Installation

This guide covers different ways to install and use the Go NPM SDK.

## Prerequisites

- Go 1.19 or later
- Git (for cloning the repository)
- Internet connection (for downloading dependencies)

## Installing the SDK

### Using Go Modules (Recommended)

The easiest way to install the Go NPM SDK is using Go modules:

```bash
go get github.com/scagogogo/go-npm-sdk
```

This will download the SDK and all its dependencies.

### Using Git Clone

You can also clone the repository directly:

```bash
git clone https://github.com/scagogogo/go-npm-sdk.git
cd go-npm-sdk
go mod download
```

## Verifying Installation

Create a simple test file to verify the installation:

```go
// test.go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    if client.IsAvailable(ctx) {
        version, err := client.Version(ctx)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("npm version: %s\n", version)
    } else {
        fmt.Println("npm not available")
    }
}
```

Run the test:

```bash
go run test.go
```

## Next Steps

- [Getting Started](./getting-started.md) - Learn the basics
- [Configuration](./configuration.md) - Configure the SDK
- [Platform Support](./platform-support.md) - Platform-specific information
