package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string

		// ContextGenerator uses the context for generating an ID in the case of AWS lambda as an example
		ContextGenerator func(c echo.Context) string
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:   DefaultSkipper,
		Generator: generator,
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = config.Generator()
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfigContextGenerator(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}

	config.Generator = nil

	if config.ContextGenerator == nil {
		// handle err?
	}

	//if config.ContextGenerator == nil {
	//	config.Generator = contextGenerator(c echo.Context())
	//}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = config.ContextGenerator(c)
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

func generator() string {
	return random.String(32)
}

func contextGenerator(c echo.Context) string {
	return random.String(32)
}
