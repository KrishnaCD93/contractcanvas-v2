package main

import (
	"encoding/json"
	"log"

	"html/template"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/KrishnaCD93/contractcanvas-v2/db"
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

func CreateDeveloper(c echo.Context, developer db.Developer) error {
	developer, err := RunDBTest(developer)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, developer)
}

func CreateDeveloperHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	role := r.FormValue("role")
	email := r.FormValue("email")
	bio := r.FormValue("bio")

	developer := db.Developer{
		Username: pgtype.Text(username),
		Firstname: pgtype.Text{
			String: firstname,
			Status: pgtype.Present,
		},
		Lastname: pgtype.Text{
			String: lastname,
			Status: pgtype.Present,
		},
		Role: pgtype.Text{
			String: role,
			Status: pgtype.Present,
		},
		Email: pgtype.Text{
			String: email,
			Status: pgtype.Present,
		},
		Bio: pgtype.Text{
			String: bio,
			Status: pgtype.Present,
		},
	}

	insertedDeveloper, err := RunDBTest(developer)

	json.NewEncoder(w).Encode(insertedDeveloper)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Renderer = newTemplate()

	developer, err := RunDBTest()
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", developer)
	})

	e.GET("/api/hello", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h2>Hello, World!<h2>")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
