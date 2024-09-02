package server

import (
	"binp/storage"
	"binp/util"
	"fmt"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type Server struct {
	store *storage.Store
	echo  *echo.Echo
	log   *zerolog.Logger
}

type CustomValidator struct {
	validator *validator.Validate
}

func setCorrectMIMETypeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if strings.HasSuffix(path, ".css") {
			c.Response().Header().Set(echo.HeaderContentType, "text/css")
		}
		return next(c)
	}
}

func NewServer(s *storage.Store) Server {
	util.InitLogger()
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	env := os.Getenv("GO_ENV")

	e.Use(middleware.RequestID())
	e.Use(util.CustomLoggerMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	allowedOrigins := []string{fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))}
	if env == "production" {
		allowedOrigins = []string{"https://binp.io"}
	}

	corsConfig := middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.POST},
		MaxAge:       300,
	}

	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Use(middleware.Secure())
	e.Use(setCorrectMIMETypeMiddleware)

	e.Static("/css", "static/css")
	e.Static("/assets", "static/assets")

	server := Server{
		store: s,
		echo:  e,
	}

	e.GET("/", server.HandleGetIndex)
	e.GET("/:id", server.HandleGetSnippet)
	e.POST("/snippet", server.HandlePostSnippet)

	return server
}

func (s *Server) Start(port string) error {
	return s.echo.Start(fmt.Sprintf(":%s", port))
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
