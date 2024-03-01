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
	Logger *slog.Logger

	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Fields available for logging
	// - id (Request ID)
	// - ip
	// - host
	// - referer
	// - user_agent
	// - status
	// - latency
	// - headers
	Fields []string
	Status int
}

var (
	// DefaultConfig is the default Logger middleware config.
	DefaultConfig = Config{
		Logger:  slog.Default(),
		Skipper: func(c echo.Context) bool { return false },
		Fields:  []string{"ip", "latency", "status"},
		Status:  0,
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

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}
			l := config.Logger.With()

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
				l = l.With("error", err)
			}
			stop := time.Now()

			if res.Status < config.Status {
				return
			}

			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			for _, field := range config.Fields {
				switch field {
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					l = l.With(field, id)
				case "ip":
					l = l.With(field, c.RealIP())
				case "host":
					l = l.With(field, req.Host)
				case "referer":
					l = l.With(field, req.Referer())
				case "user_agent":
					l = l.With(field, req.UserAgent())
				case "status":
					l = l.With(field, res.Status)
				case "latency":
					l = l.With(field, stop.Sub(start).String())
				case "headers":
					l = l.With(field, req.Header)
				}
			}

			msg := fmt.Sprintf("%s %s", req.Method, path)
			switch {
			case res.Status >= 500:
				l.Error(msg)
			case res.Status >= 400:
				l.Warn(msg)
			default:
				l.Debug(msg)
			}

			return
		}
	}
}
