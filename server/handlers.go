package server

import (
	"binp/storage"
	"binp/util"
	"binp/views"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type PostSnippetReq struct {
	Text          string `form:"text" validate:"required,min=1,max=10000"`
	BurnAfterRead bool   `form:"burn_after_read"`
	Language      string `form:"language" validate:"required,oneof=plaintext bash css docker go html javascript json markdown python ruby typescript"`
	Expiry        string `form:"expiry" validate:"required,oneof=one_hour one_day one_week one_month"`
}

func (s *Server) HandleGetIndex(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	logger.Info().Msg("Index page")
	return Render(c, http.StatusOK, views.Index())
}

func (s *Server) HandleGetSnippet(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	id := c.Param("id")
	snippet, err := s.store.GetSnippetByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("GetSnippetByID")
		return err
	}
	if snippet == nil {
		logger.Warn().Str("id", id).Msg("Snippet not found")
		return Render(c, http.StatusNotFound, views.NotFoundPage())
	}
	logger.Info().Interface("snippet", snippet).Msg("Snippet found")
	if snippet.ExpiresAt.Before(time.Now().UTC()) {
		logger.Warn().Str("id", id).Msg("Snippet expired")
		err = s.store.DeleteSnippet(snippet.ID)
		if err != nil {
			logger.Error().Err(err).Msg("DeleteSnippet")
		}
		return Render(c, http.StatusNotFound, views.NotFoundPage())
	}
	if !snippet.IsRead {
		snippet.IsRead = true
		err = s.store.UpdateSnippet(snippet)
		if err != nil {
			return err
		}
	} else if snippet.BurnAfterRead {
		err = s.store.DeleteSnippet(snippet.ID)
		if err != nil {
			return err
		}
	}

	highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
	if err != nil {
		logger.Error().Err(err).Msg("HighlightCode")
		highlightedCode = snippet.Text
	}

	return Render(c, http.StatusOK, views.SnippetPage(snippet, highlightedCode))
}

func (s *Server) HandlePostSnippet(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	data := new(PostSnippetReq)
	if err := c.Bind(data); err != nil {
		logger.Error().Err(err).Msg("Bind")
		return err
	}
	logger.Info().Interface("data", data).Msg("PostSnippet")
	if err := c.Validate(data); err != nil {
		logger.Error().Err(err).Msg("Validate")
		return err
	}
	snippet, err := s.store.CreateSnippet(data.Text, data.BurnAfterRead, storage.GetSnippetExpiration(data.Expiry), data.Language)
	if err != nil {
		logger.Error().Err(err).Msg("CreateSnippet")
		return err
	}
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/%s", snippet.ID))
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
