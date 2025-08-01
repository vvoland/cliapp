# cliapp

I always forget how to setup the slog to output nice colored output so I made this.

And also handle context cancellation via SIGINT.

## Usage

```go
package main

import (
    "log/slog"
    "grono.dev/cliapp"
)

func main() {
    ctx := cliapp.Init()
    
    slog.InfoContext(ctx, "Application started")
    
    // Your application logic here
    <-ctx.Done()
    
    slog.InfoContext(ctx, "Application stopped")
}
```

## Features

- Structured logging with slog
- Terminal-aware colored output
- Graceful shutdown on SIGINT
