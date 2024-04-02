package main

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"

	"github.com/KrishnaCD93/contractcanvas-v2/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func RunDBTest(devInfo db.Developer) (db.Developer, error) {
	var developer struct {
		Username  pgtype.Text `json:"username"`
		Firstname pgtype.Text `json:"firstname"`
		Lastname  pgtype.Text `json:"lastname"`
		Role      pgtype.Text `json:"role"`
		Email     pgtype.Text `json:"email"`
		Bio       pgtype.Text `json:"bio"`
	}
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	postgresURL := os.Getenv("POSTGRES_URL")
	conn, err := pgx.Connect(ctx, postgresURL)
	if err != nil {
		return db.Developer{}, err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// get all developers
	developers, err := queries.GetDevelopers(ctx)
	if err != nil {
		return db.Developer{}, err
	}

	log.Println(developers)

	// insert a developer
	insertedDeveloper, err := queries.CreateDeveloper(ctx, db.CreateDeveloperParams{
		Username:  developer.Username,
		Firstname: developer.Firstname,
		Lastname:  developer.Lastname,
		Role:      developer.Role,
		Email:     developer.Email,
		Bio:       developer.Bio,
	})
	if err != nil {
		return db.Developer{}, err
	}

	// get a developer
	getDeveloper, err := queries.GetDeveloper(ctx, insertedDeveloper.ID)
	if err != nil {
		return db.Developer{}, err
	}

	log.Println(reflect.DeepEqual(getDeveloper, insertedDeveloper))

	return getDeveloper, nil
}
