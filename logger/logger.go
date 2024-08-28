package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

var (
	zlog zerolog.Logger
	once sync.Once
)

type EchoLogger struct {
	zl     zerolog.Logger
	header string
}

func InitLogger() {
	once.Do(func() {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
		zlog = zerolog.New(output).With().Timestamp().Caller().Logger()
	})
}

func GetLogger() *zerolog.Logger {
	return &zlog
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
				Str("id", id).
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

// GetEchoLogger returns an Echo-compatible logger
func GetEchoLogger() echo.Logger {
	return &EchoLogger{zl: zlog}
}

// Level returns the current log level
func (l *EchoLogger) Level() echolog.Lvl {
	return echolog.INFO
}

// SetLevel sets the log level (not implemented for simplicity)
func (l *EchoLogger) SetLevel(v echolog.Lvl) {
	// No-op for simplicity
}

// Output returns the logger's output (not implemented for simplicity)
func (l *EchoLogger) Output() io.Writer {
	return os.Stdout
}

// SetOutput sets the logger's output (not implemented for simplicity)
func (l *EchoLogger) SetOutput(w io.Writer) {
	// No-op for simplicity
}

// Prefix returns the logger's prefix (not implemented for simplicity)
func (l *EchoLogger) Prefix() string {
	return ""
}

// SetPrefix sets the logger's prefix (not implemented for simplicity)
func (l *EchoLogger) SetPrefix(p string) {
	// No-op for simplicity
}

// SetHeader sets the header for the logger
func (l *EchoLogger) SetHeader(h string) {
	l.header = h
}

// Print writes a log message
func (l *EchoLogger) Print(i ...interface{}) {
	l.zl.Info().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Printf writes a formatted log message
func (l *EchoLogger) Printf(format string, i ...interface{}) {
	l.zl.Info().Str("header", l.header).Msgf(format, i...)
}

// Printj writes a JSON log message
func (l *EchoLogger) Printj(j echolog.JSON) {
	l.zl.Info().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}

// Debug writes a debug log message
func (l *EchoLogger) Debug(i ...interface{}) {
	l.zl.Debug().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Debugf writes a formatted debug log message
func (l *EchoLogger) Debugf(format string, i ...interface{}) {
	l.zl.Debug().Str("header", l.header).Msgf(format, i...)
}

// Debugj writes a JSON debug log message
func (l *EchoLogger) Debugj(j echolog.JSON) {
	l.zl.Debug().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}

// Info writes an info log message
func (l *EchoLogger) Info(i ...interface{}) {
	l.zl.Info().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Infof writes a formatted info log message
func (l *EchoLogger) Infof(format string, i ...interface{}) {
	l.zl.Info().Str("header", l.header).Msgf(format, i...)
}

// Infoj writes a JSON info log message
func (l *EchoLogger) Infoj(j echolog.JSON) {
	l.zl.Info().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}

// Warn writes a warning log message
func (l *EchoLogger) Warn(i ...interface{}) {
	l.zl.Warn().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Warnf writes a formatted warning log message
func (l *EchoLogger) Warnf(format string, i ...interface{}) {
	l.zl.Warn().Str("header", l.header).Msgf(format, i...)
}

// Warnj writes a JSON warning log message
func (l *EchoLogger) Warnj(j echolog.JSON) {
	l.zl.Warn().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}

// Error writes an error log message
func (l *EchoLogger) Error(i ...interface{}) {
	l.zl.Error().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Errorf writes a formatted error log message
func (l *EchoLogger) Errorf(format string, i ...interface{}) {
	l.zl.Error().Str("header", l.header).Msgf(format, i...)
}

// Errorj writes a JSON error log message
func (l *EchoLogger) Errorj(j echolog.JSON) {
	l.zl.Error().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}

// Fatal writes a fatal log message
func (l *EchoLogger) Fatal(i ...interface{}) {
	l.zl.Fatal().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Fatalf writes a formatted fatal log message
func (l *EchoLogger) Fatalf(format string, i ...interface{}) {
	l.zl.Fatal().Str("header", l.header).Msgf(format, i...)
}

// Fatalj writes a JSON fatal log message
func (l *EchoLogger) Fatalj(j echolog.JSON) {
	l.zl.Fatal().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}

// Panic writes a panic log message
func (l *EchoLogger) Panic(i ...interface{}) {
	l.zl.Panic().Str("header", l.header).Msg(fmt.Sprint(i...))
}

// Panicf writes a formatted panic log message
func (l *EchoLogger) Panicf(format string, i ...interface{}) {
	l.zl.Panic().Str("header", l.header).Msgf(format, i...)
}

// Panicj writes a JSON panic log message
func (l *EchoLogger) Panicj(j echolog.JSON) {
	l.zl.Panic().Str("header", l.header).Fields(map[string]interface{}(j)).Msg("")
}
