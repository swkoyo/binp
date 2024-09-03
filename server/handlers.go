package server

import (
	"binp/storage"
	"binp/util"
	"binp/views"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type PostSnippetReq struct {
	Text          string `form:"text" json:"text" validate:"required,min=1,max=10000"`
	BurnAfterRead bool   `form:"burn_after_read" json:"burn_after_read"`
	Language      string `form:"language" json:"language" validate:"required,oneof=plaintext bash css docker go html javascript json markdown python ruby typescript"`
	Expiry        string `form:"expiry" json:"expiry" validate:"required,oneof=one_hour one_day one_week one_month"`
}

func (s *Server) HandleGetIndex(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	logger.Info().Msg("Index page")
	return Render(c, http.StatusOK, views.Index())
}

func (s *Server) HandleGetSnippet(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	id := c.Param("id")
	contentType := c.Request().Header.Get("Content-Type")

	snippet, err := s.store.GetSnippetByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("GetSnippetByID")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		} else {
			return err
		}
	}

	if snippet == nil {
		logger.Warn().Str("id", id).Msg("Snippet not found")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Snippet not found"})
		} else {
			return Render(c, http.StatusNotFound, views.NotFoundPage())
		}
	}

	logger.Debug().Interface("snippet", snippet).Msg("Snippet found")
	if snippet.ExpiresAt.Before(time.Now().UTC()) {
		logger.Warn().Str("id", id).Msg("Snippet expired")
		err = s.store.DeleteSnippet(snippet.ID)
		if err != nil {
			logger.Error().Err(err).Msg("DeleteSnippet")
		}
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Snippet not found"})
		} else {
			return Render(c, http.StatusNotFound, views.NotFoundPage())
		}
	}

	if !snippet.IsRead {
		snippet.IsRead = true
		err = s.store.UpdateSnippet(snippet)
		if err != nil {
			return err
		}
	}

	if snippet.BurnAfterRead {
		err = s.store.DeleteSnippet(snippet.ID)
		if err != nil {
			return err
		}
	}

	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return c.JSON(http.StatusOK, snippet)
	} else {
		highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
		if err != nil {
			logger.Error().Err(err).Msg("HighlightCode")
			highlightedCode = snippet.Text
		}

		return Render(c, http.StatusOK, views.SnippetPage(snippet, highlightedCode))
	}
}

func (s *Server) HandlePostSnippet(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	data := new(PostSnippetReq)
	contentType := c.Request().Header.Get("Content-Type")

	if err := c.Bind(data); err != nil {
		logger.Error().Err(err).Msg("Bind")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		} else {
			return err
		}
	}

	logger.Debug().Interface("data", data).Msg("PostSnippet")
	if err := c.Validate(data); err != nil {
		logger.Error().Err(err).Msg("Validate")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			return err
		}
	}

	snippet, err := s.store.CreateSnippet(data.Text, data.BurnAfterRead, storage.GetSnippetExpiration(data.Expiry), data.Language)
	if err != nil {
		logger.Error().Err(err).Msg("CreateSnippet")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			return err
		}
	}

	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return c.JSON(http.StatusCreated, snippet)
	} else {
		highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
		if err != nil {
			logger.Error().Err(err).Msg("HighlightCode")
			highlightedCode = snippet.Text
		}

		c.Response().Header().Set("Hx-Push-Url", fmt.Sprintf("/%s", snippet.ID))
		return Render(c, http.StatusCreated, views.SnippetDetails(snippet, highlightedCode))
	}
}
