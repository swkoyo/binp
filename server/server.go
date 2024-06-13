package server

import (
	"binp/storage"
	"fmt"

	"github.com/a-h/templ"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	store *storage.Store
	echo  *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewServer(s *storage.Store) Server {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Static("/css", "views/css")

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
