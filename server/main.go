package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/KrishnaCD93/contractcanvas-v2/db"

	"html/template"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func GetDevelopers(conn *pgx.Conn, ctx context.Context) ([]db.Developer, error) {
	queries := db.New(conn)

	developers, err := queries.GetDevelopers(ctx)
	if err != nil {
		return nil, err
	}

	return developers, nil
}

func InsertDeveloper(devInfo db.Developer, conn *pgx.Conn, ctx context.Context) (db.Developer, error) {

	queries := db.New(conn)

	insertedDeveloper, err := queries.CreateDeveloper(ctx, db.CreateDeveloperParams{
		Username:  devInfo.Username,
		Firstname: devInfo.Firstname,
		Lastname:  devInfo.Lastname,
		Role:      devInfo.Role,
		Email:     devInfo.Email,
		Bio:       devInfo.Bio,
	})
	if err != nil {
		return db.Developer{}, err
	}

	return insertedDeveloper, nil
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

	ctx := context.Background()
	postgresURL := os.Getenv("POSTGRES_URL")
	conn, err := pgx.Connect(ctx, postgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	e.GET("/api/developers", func(c echo.Context) error {
		developers, err := GetDevelopers(conn, ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, developers)
	})

	e.POST("/api/createDeveloper", func(c echo.Context) error {
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
		insertedDeveloper, err := InsertDeveloper(developer, conn, ctx)
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
