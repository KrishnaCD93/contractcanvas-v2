// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Client struct {
	ID        int32
	Username  pgtype.Text
	Firstname pgtype.Text
	Lastname  pgtype.Text
	Email     pgtype.Text
	Bio       pgtype.Text
	CreatedAt pgtype.Timestamp
}

type Deliverable struct {
	ID          int32
	ProjectID   pgtype.Int4
	Description pgtype.Text
	Type        interface{}
	Budget      pgtype.Float8
	CreatedAt   pgtype.Timestamp
}

type Developer struct {
	ID        int32
	Username  pgtype.Text
	Firstname pgtype.Text
	Lastname  pgtype.Text
	Role      pgtype.Text
	Email     pgtype.Text
	Bio       pgtype.Text
	CreatedAt pgtype.Timestamp
}

type PercentComplete struct {
	ID              int32
	PercentComplete pgtype.Float8
	DelivID         pgtype.Int4
	ReportPeriod    pgtype.Date
	CreatedAt       pgtype.Timestamp
}

type Project struct {
	ID          int32
	Title       pgtype.Text
	Description pgtype.Text
	ClientID    pgtype.Int4
	DevID       pgtype.Int4
	Status      pgtype.Text
	CreatedAt   pgtype.Timestamp
}

type Timesheet struct {
	ID         int32
	Hours      pgtype.Float8
	DelivID    pgtype.Int4
	DateOfWork pgtype.Date
	CreatedAt  pgtype.Timestamp
}

type WorkingRelationship struct {
	ID          int32
	ClientID    pgtype.Int4
	DeveloperID pgtype.Int4
}
