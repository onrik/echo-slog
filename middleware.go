package echoslog

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// Skipper defines a function to skip middleware.
type Skipper func(c echo.Context) bool

// Config defines the config for Logger middleware.
type Config struct {
	// Logger - slog instance
	Logger *slog.Logger

	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Fields - list of fields
	Fields []Field

	// MinStatus - minimum http status value for logging
	MinStatus int

	Attrs func(config Config, c echo.Context, start time.Time) []any
}

var (
	// DefaultConfig is the default Logger middleware config.
	DefaultConfig = Config{
		Logger:    slog.Default(),
		Skipper:   func(c echo.Context) bool { return false },
		Fields:    FieldsDefault(),
		Attrs:     AttrsDefault,
		MinStatus: 0,
	}
)

func MiddlewareDefault() echo.MiddlewareFunc {
	return Middleware(DefaultConfig)
}

// Middleware returns a Logger middleware with config.
func Middleware(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	if config.Logger == nil {
		config.Logger = DefaultConfig.Logger
	}
	if config.Attrs == nil {
		config.Attrs = DefaultConfig.Attrs
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			err = next(c)
			if err != nil {
				c.Error(err)
			}
			if res.Status < config.MinStatus {
				return err
			}

			attrs := config.Attrs(config, c, start)
			if err != nil {
				if _, ok := err.(*echo.HTTPError); !ok {
					attrs = append(attrs, "error", err)
				}
			}

			msg := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
			switch {
			case res.Status >= 500:
				config.Logger.Error(msg, attrs...)
			case res.Status >= 400:
				config.Logger.Warn(msg, attrs...)
			default:
				config.Logger.Info(msg, attrs...)
			}

			return err
		}
	}
}
