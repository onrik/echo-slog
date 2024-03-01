# echo-slog
Logger for echo based on log/slog.

```golang
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/onrik/echo-slog"
)

func main() {
    server := echo.New()
    server.Use(echoslog.MiddlewareDefault())
}

```