package echoslog

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Field string

const (
	FieldID        = Field("id")
	FieldIP        = Field("ip")
	FieldLatency   = Field("latency")
	FieldStatus    = Field("status")
	FieldReferer   = Field("referer")
	FieldUserAgent = Field("user_agent")
	FieldHeaders   = Field("headers")
)

func FieldsDefault() []Field {
	return []Field{
		FieldLatency,
		FieldStatus,
	}
}

func FieldsAll() []Field {
	return []Field{
		FieldID,
		FieldIP,
		FieldLatency,
		FieldStatus,
		FieldHeaders,
	}
}

func AttrsDefault(config Config, c echo.Context, start time.Time) []any {
	attrs := make([]any, 0, len(config.Fields)*2)

	for _, field := range config.Fields {
		switch field {
		case FieldID:
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}
			attrs = append(attrs, string(field), id)
		case FieldIP:
			attrs = append(attrs, string(field), c.RealIP())
		case FieldReferer:
			attrs = append(attrs, string(field), c.Request().Referer())
		case FieldUserAgent:
			attrs = append(attrs, string(field), c.Request().UserAgent())
		case FieldStatus:
			attrs = append(attrs, string(field), c.Response().Status)
		case FieldLatency:
			attrs = append(attrs, string(field), time.Since(start).String())
		case FieldHeaders:
			attrs = append(attrs, string(field), c.Request().Header)
		}
	}

	return attrs
}
