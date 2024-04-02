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
	log.Println("CreateDeveloper")

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
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
		Firstname: pgtype.Text{
			String: firstname,
			Valid:  true,
		},
		Lastname: pgtype.Text{
			String: lastname,
			Valid:  true,
		},
		Role: pgtype.Text{
			String: role,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
		Bio: pgtype.Text{
			String: bio,
			Valid:  true,
		},
	}

	insertedDeveloper, err := RunDBTest(developer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(insertedDeveloper)
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

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	// e.GET("/developers", func(c echo.Context) error {
	// 	developers, err := GetDevelopers()
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}

	// 	return c.JSON(http.StatusOK, developers)
	// })

	e.POST("/createDeveloper", func(c echo.Context) error {
		// Parse the form data from the request
		err := c.Request().ParseForm()
		if err != nil {
			log.Printf("Error parsing form data: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// Extract the form data
		username := c.FormValue("username")
		firstname := c.FormValue("firstname")
		lastname := c.FormValue("lastname")
		role := c.FormValue("role")
		email := c.FormValue("email")
		bio := c.FormValue("bio")

		// Create the developer object
		developer := db.Developer{
			Username:  pgtype.Text{String: username, Valid: true},
			Firstname: pgtype.Text{String: firstname, Valid: true},
			Lastname:  pgtype.Text{String: lastname, Valid: true},
			Role:      pgtype.Text{String: role, Valid: true},
			Email:     pgtype.Text{String: email, Valid: true},
			Bio:       pgtype.Text{String: bio, Valid: true},
		}

		// Insert the developer into the database
		insertedDeveloper, err := RunDBTest(developer)
		if err != nil {
			log.Printf("Error inserting developer: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		log.Println(insertedDeveloper)
		return c.JSON(http.StatusOK, insertedDeveloper)
	})
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
