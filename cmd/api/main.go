package main

import (
	"binp/models"
	"binp/storage"
	"binp/views"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbStore, err := storage.GetDatabaseStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := dbStore.Init(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/css", "views/css")

	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, views.Index(nil))
	})

	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		snippet, err := models.GetSnippetByID(id)
		if err != nil {
			return err
		}
		return Render(c, http.StatusOK, views.Index(snippet))
	})

	e.POST("/snippet", func(c echo.Context) error {
		text := c.FormValue("text")
		snippet, err := models.CreateSnippet(text)
		if err != nil {
			return err
		}
		c.Response().Header().Set("HX-Redirect", "/"+snippet.Id)
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
