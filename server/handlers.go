package server

import (
	"binp/storage"
	"binp/util"
	"binp/views"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostSnippetReq struct {
	Text          string `form:"text" validate:"required,min=1,max=1000"`
	BurnAfterRead bool   `form:"burn_after_read"`
	Expiry        string `form:"expiry" validate:"required,oneof=one_hour one_day one_week one_month"`
}

func (s *Server) HandleGetIndex(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	logger.Info().Msg("Index page")
	return Render(c, http.StatusOK, views.Index())
}

func (s *Server) HandleGetSnippet(c echo.Context) error {
	id := c.Param("id")
	snippet, err := s.store.GetSnippetByID(id)
	if err != nil {
		return err
	}
	if snippet == nil {
		return Render(c, http.StatusNotFound, views.NotFoundPage())
	}
	if !snippet.IsRead {
		err = s.store.SetSnippetIsRead(id)
		if err != nil {
			return err
		}
		snippet.IsRead = true
	} else if snippet.BurnAfterRead {
		err = s.store.DeleteSnippet(id)
		if err != nil {
			return err
		}
	}
	return Render(c, http.StatusOK, views.SnippetPage(snippet))
}

func (s *Server) HandlePostSnippet(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	data := new(PostSnippetReq)
	if err := c.Bind(data); err != nil {
		return err
	}
	logger.Info().Interface("data", data).Msg("PostSnippet")
	if err := c.Validate(data); err != nil {
		return err
	}
	snippet, err := s.store.CreateSnippet(data.Text, data.BurnAfterRead, storage.GetSnippetExpiration(data.Expiry))
	if err != nil {
		return err
	}
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/%s", snippet.ID))
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
