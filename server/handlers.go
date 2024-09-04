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
	Language      string `form:"language" json:"language" validate:"required"`
	Expiry        string `form:"expiry" json:"expiry" validate:"required"`
}

func (s *Server) HandleGetIndex(c echo.Context) error {
	return Render(c, http.StatusOK, views.Index())
}

func (s *Server) HandleGetSnippet(c echo.Context) error {
	logger := util.GetLoggerWithRequestID(c)
	id := c.Param("id")
	contentType := c.Request().Header.Get("Content-Type")

	snippet, err := s.store.GetSnippetByID(id)
	if err != nil {
		logger.Error().Str("ID", id).Err(err).Msg("Error while getting snippet")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		} else {
			return Render(c, http.StatusInternalServerError, views.ErrorPage())
		}
	}

	if snippet == nil {
		logger.Warn().Str("ID", id).Msg("Snippet not found")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Snippet not found"})
		} else {
			return Render(c, http.StatusNotFound, views.NotFoundPage())
		}
	}

	logger.Debug().Str("ID", id).Interface("snippet", snippet).Msg("Snippet found")
	if snippet.ExpiresAt.Before(time.Now().UTC()) {
		logger.Warn().Str("ID", id).Msg("Snippet expired")
		err = s.store.DeleteSnippet(snippet.ID)
		if err != nil {
			logger.Error().Str("ID", id).Err(err).Msg("Error while deleting expired snippet")
			if strings.HasPrefix(contentType, "application/json") {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			} else {
				return Render(c, http.StatusInternalServerError, views.ErrorPage())
			}
		}
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Snippet not found"})
		} else {
			return Render(c, http.StatusNotFound, views.NotFoundPage())
		}
	}

	if snippet.BurnAfterRead {
		logger.Info().Str("ID", id).Msg("Burned snippet")
		err = s.store.DeleteSnippet(snippet.ID)
		if err != nil {
			logger.Error().Str("ID", id).Err(err).Msg("Error while deleting burned snippet")
			if strings.HasPrefix(contentType, "application/json") {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			} else {
				return Render(c, http.StatusInternalServerError, views.ErrorPage())
			}
		}
	}

	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return c.JSON(http.StatusOK, snippet)
	} else {
		highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
		if err != nil {
			logger.Error().Err(err).Msg("Error while highlighting code")
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
		logger.Error().Err(err).Msg("Error while binding data")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		} else {
			return Render(c, http.StatusBadRequest, views.ErrorAlert("Invalid request data"))
		}
	}

	logger.Debug().Interface("data", data).Msg("Creating snippet")
	if err := c.Validate(data); err != nil {
		logger.Error().Err(err).Msg("Error while validating data")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			return Render(c, http.StatusBadRequest, views.ErrorAlert("Invalid request data"))
		}
	}

	if !storage.IsValidLanguage(data.Language) {
		logger.Warn().Str("language", data.Language).Msg("Invalid language")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid language. Options: %v", storage.GetValidLanguages())})
		} else {
			return Render(c, http.StatusBadRequest, views.ErrorAlert("Invalid language"))
		}
	}

	if !storage.IsValidExpiration(data.Expiry) {
		logger.Warn().Str("expiry", data.Expiry).Msg("Invalid expiry")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid expiry. Options: %v", storage.GetValidExpirations())})
		} else {
			return Render(c, http.StatusBadRequest, views.ErrorAlert("Invalid expiry"))
		}
	}

	snippet, err := s.store.CreateSnippet(data.Text, data.BurnAfterRead, storage.GetSnippetExpiration(data.Expiry), data.Language)
	if err != nil {
		logger.Error().Err(err).Msg("Error while creating snippet")
		if strings.HasPrefix(contentType, "application/json") {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		} else {
			return Render(c, http.StatusInternalServerError, views.ErrorAlert("Failed to create snippet"))
		}
	}

	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return c.JSON(http.StatusCreated, snippet)
	} else {
		highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
		if err != nil {
			logger.Error().Err(err).Msg("Error while highlighting code")
			highlightedCode = snippet.Text
		}

		c.Response().Header().Set("Hx-Push-Url", fmt.Sprintf("/%s", snippet.ID))
		return Render(c, http.StatusCreated, views.PostSnippetResponse(snippet, highlightedCode))
	}
}
