package server

import (
	"binp/views"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostSnippetReq struct {
	Text string `form:"text" validate:"required,min=1,max=1000"`
}

func (s *Server) HandleGetIndex(c echo.Context) error {
	return Render(c, http.StatusOK, views.Index(nil))
}

func (s *Server) HandleGetSnippet(c echo.Context) error {
	id := c.Param("id")
	snippet, err := s.store.GetSnippetByID(id)
	if err != nil {
		return err
	}
	return Render(c, http.StatusOK, views.Index(snippet))
}

func (s *Server) HandlePostSnippet(c echo.Context) error {
	data := new(PostSnippetReq)
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	snippet, err := s.store.CreateSnippet(data.Text)
	if err != nil {
		return err
	}
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/%s", snippet.ID))
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
