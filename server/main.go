package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Renderer = newTemplate()

	// Routes
	e.GET("/", func(c echo.Context) error {
		count := Count{Count: 0}
		return c.Render(http.StatusOK, "index.html", count)
	})

	e.GET("/api/hello", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h2>Hello, World!<h2>")
	})
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
