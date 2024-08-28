package logger

import (
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var (
	log  zerolog.Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
		log = zerolog.New(output).With().Timestamp().Caller().Logger()
	})
}

func GetLogger() *zerolog.Logger {
	return &log
}

func HTTPLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			if err := next(c); err != nil {
				c.Error(err)
			}
			log.Info().
				Str("request_id", res.Header().Get(echo.HeaderXRequestID)).
				Str("remote_ip", c.RealIP()).
				Str("host", req.Host).
				Str("uri", req.RequestURI).
				Str("method", req.Method).
				Int("status", res.Status).
				Int64("bytes", res.Size).
				Str("user_agent", req.UserAgent()).
				Msg("request")
			return nil
		}
	}
}
