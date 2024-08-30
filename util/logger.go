package util

import (
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var (
	zlog zerolog.Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
		zlog = zerolog.New(output).With().Timestamp().Caller().Logger()
	})
}

func GetLogger() *zerolog.Logger {
	return &zlog
}

func GetLoggerWithRequestID(c echo.Context) zerolog.Logger {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	return zlog.With().Str("request_id", requestID).Logger()
}

func CustomLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			zlog.Info().
				Str("request_id", id).
				Str("remote_ip", c.RealIP()).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("user_agent", req.UserAgent()).
				Int("status", res.Status).
				Dur("latency", time.Since(start)).
				Msg("Request")

			return err
		}
	}
}
