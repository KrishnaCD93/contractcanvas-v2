package main

import (
	"context"
	"log"
	"reflect"

	"github.com/KrishnaCD93/contractcanvas-v2/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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

func runDBTest() ([]*string, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=postgres password=postgres dbname=postgres sslmode=verify-full")
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// get all developers
	developers, err := queries.GetDevelopers(ctx)
	if err != nil {
		return nil, err
	}

	log.Println(developers)

	// insert a developer
	insertDeveloper, err := queries.CreateDeveloper(ctx, db.CreateDeveloperParams{
		Username:  "Krishna",
		Firstname: "Krishna",
		Lastname:  "Duvvuri",
		Role:      "Software Engineer",
		Email:     "krishna.c.duvvuri",
		Bio: pgtype.Text{
			String: "I am a software engineer and I want to build cool things!",
		},
	})
	if err != nil {
		return nil, err
	}

	// get a developer
	developer, err := queries.GetDeveloper(ctx, insertDeveloper.ID)
	if err != nil {
		return nil, err
	}

	log.Println(reflect.DeepEqual(developer, insertDeveloper))

	return developer, nil
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Renderer = newTemplate()

	developers, err := runDBTest()
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	e.GET("/", func(c echo.Context) error {
		developer := developers[0]
		return c.Render(http.StatusOK, "index.html", developer)
	})

	e.GET("/api/hello", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h2>Hello, World!<h2>")
	})
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
