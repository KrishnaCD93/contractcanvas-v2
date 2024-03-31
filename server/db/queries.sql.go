// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const burnedValue = `-- name: BurnedValue :many
SELECT
    d.project_id ProjectId,
    d.description Deliverable,
    d.type Type,
    to_char(date_of_work, 'YYYY-MM') ReportingPeriod,
    SUM(t.hours) BurnedValue
FROM
    timesheet t
    JOIN
    deliverables d ON t.deliv_id = d.id
WHERE
    d.project_id = $1
`

type BurnedValueRow struct {
	Projectid       pgtype.Int4
	Deliverable     pgtype.Text
	Type            interface{}
	Reportingperiod string
	Burnedvalue     int64
}

func (q *Queries) BurnedValue(ctx context.Context, projectID pgtype.Int4) ([]BurnedValueRow, error) {
	rows, err := q.db.Query(ctx, burnedValue, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BurnedValueRow
	for rows.Next() {
		var i BurnedValueRow
		if err := rows.Scan(
			&i.Projectid,
			&i.Deliverable,
			&i.Type,
			&i.Reportingperiod,
			&i.Burnedvalue,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const costPerformanceIndex = `-- name: CostPerformanceIndex :many
SELECT
    EarnedValue.ProjectId,
    EarnedValue.Deliverable,
    EarnedValue.Type,
    EarnedValue.ReportingPeriod,
    EarnedValue.EarnedValue,
    BurnedValue.BurnedValue,
    EarnedValue.EarnedValue / BurnedValue.BurnedValue CPI
FROM
    (SELECT
        d.project_id ProjectId,
        d.description Deliverable,
        d.type Type,
        percent_complete PercentComplete,
        budget Budget,
        to_char(reporting_period, 'YYYY-MM') ReportingPeriod,
        percent_complete * budget EarnedValue
    FROM
        percent_complete pc
        JOIN
        deliverables d ON pc.deliv_id = d.id
    WHERE d.project_id = $1) EarnedValue
    JOIN
    (SELECT
        d.project_id ProjectId,
        d.description Deliverable,
        d.type Type,
        to_char(date_of_work, 'YYYY-MM') ReportingPeriod,
        SUM(t.hours) BurnedValue
    FROM
        timesheet t
        JOIN
        deliverables d ON t.deliv_id = d.id
    WHERE d.project_id = $1) BurnedValue ON EarnedValue.Deliverable = BurnedValue.Deliverable AND EarnedValue.ReportingPeriod = BurnedValue.ReportingPeriod
`

type CostPerformanceIndexRow struct {
	Projectid       pgtype.Int4
	Deliverable     pgtype.Text
	Type            interface{}
	Reportingperiod string
	Earnedvalue     int32
	Burnedvalue     int64
	Cpi             int32
}

func (q *Queries) CostPerformanceIndex(ctx context.Context, projectID pgtype.Int4) ([]CostPerformanceIndexRow, error) {
	rows, err := q.db.Query(ctx, costPerformanceIndex, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CostPerformanceIndexRow
	for rows.Next() {
		var i CostPerformanceIndexRow
		if err := rows.Scan(
			&i.Projectid,
			&i.Deliverable,
			&i.Type,
			&i.Reportingperiod,
			&i.Earnedvalue,
			&i.Burnedvalue,
			&i.Cpi,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createClient = `-- name: CreateClient :one
INSERT INTO clients
    (username, firstname, lastname)
VALUES
    ($1, $2, $3)
RETURNING id, username, firstname, lastname, created_at
`

type CreateClientParams struct {
	Username  pgtype.Text
	Firstname pgtype.Text
	Lastname  pgtype.Text
}

func (q *Queries) CreateClient(ctx context.Context, arg CreateClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, createClient, arg.Username, arg.Firstname, arg.Lastname)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.CreatedAt,
	)
	return i, err
}

const createDeliverable = `-- name: CreateDeliverable :one
INSERT INTO deliverables
    (project_id, description, type, budget)
VALUES
    ($1, $2, $3, $4)
RETURNING id, project_id, description, type, budget, created_at
`

type CreateDeliverableParams struct {
	ProjectID   pgtype.Int4
	Description pgtype.Text
	Type        interface{}
	Budget      pgtype.Float8
}

func (q *Queries) CreateDeliverable(ctx context.Context, arg CreateDeliverableParams) (Deliverable, error) {
	row := q.db.QueryRow(ctx, createDeliverable,
		arg.ProjectID,
		arg.Description,
		arg.Type,
		arg.Budget,
	)
	var i Deliverable
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Description,
		&i.Type,
		&i.Budget,
		&i.CreatedAt,
	)
	return i, err
}

const createDeveloper = `-- name: CreateDeveloper :one
INSERT INTO developers
    (username, firstname, lastname, role)
VALUES
    ($1, $2, $3, $4)
RETURNING id, username, firstname, lastname, role, created_at
`

type CreateDeveloperParams struct {
	Username  pgtype.Text
	Firstname pgtype.Text
	Lastname  pgtype.Text
	Role      pgtype.Text
}

func (q *Queries) CreateDeveloper(ctx context.Context, arg CreateDeveloperParams) (Developer, error) {
	row := q.db.QueryRow(ctx, createDeveloper,
		arg.Username,
		arg.Firstname,
		arg.Lastname,
		arg.Role,
	)
	var i Developer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const createPercentComplete = `-- name: CreatePercentComplete :one
INSERT INTO percent_complete
    (percent_complete, deliv_id, report_period)
VALUES
    ($1, $2, $3)
RETURNING id, percent_complete, deliv_id, report_period, created_at
`

type CreatePercentCompleteParams struct {
	PercentComplete pgtype.Float8
	DelivID         pgtype.Int4
	ReportPeriod    pgtype.Date
}

func (q *Queries) CreatePercentComplete(ctx context.Context, arg CreatePercentCompleteParams) (PercentComplete, error) {
	row := q.db.QueryRow(ctx, createPercentComplete, arg.PercentComplete, arg.DelivID, arg.ReportPeriod)
	var i PercentComplete
	err := row.Scan(
		&i.ID,
		&i.PercentComplete,
		&i.DelivID,
		&i.ReportPeriod,
		&i.CreatedAt,
	)
	return i, err
}

const createProject = `-- name: CreateProject :one
INSERT INTO projects
    (title, description, client_id, dev_id, status)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING id, title, description, client_id, dev_id, status, created_at
`

type CreateProjectParams struct {
	Title       pgtype.Text
	Description pgtype.Text
	ClientID    pgtype.Int4
	DevID       pgtype.Int4
	Status      pgtype.Text
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, createProject,
		arg.Title,
		arg.Description,
		arg.ClientID,
		arg.DevID,
		arg.Status,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.ClientID,
		&i.DevID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const createTimesheet = `-- name: CreateTimesheet :one
INSERT INTO timesheet
    (hours, deliv_id, date_of_work)
VALUES
    ($1, $2, $3)
RETURNING id, hours, deliv_id, date_of_work, created_at
`

type CreateTimesheetParams struct {
	Hours      pgtype.Float8
	DelivID    pgtype.Int4
	DateOfWork pgtype.Date
}

func (q *Queries) CreateTimesheet(ctx context.Context, arg CreateTimesheetParams) (Timesheet, error) {
	row := q.db.QueryRow(ctx, createTimesheet, arg.Hours, arg.DelivID, arg.DateOfWork)
	var i Timesheet
	err := row.Scan(
		&i.ID,
		&i.Hours,
		&i.DelivID,
		&i.DateOfWork,
		&i.CreatedAt,
	)
	return i, err
}

const createWorkingRelationship = `-- name: CreateWorkingRelationship :one
INSERT INTO working_relationships
    (client_id, developer_id)
VALUES
    ($1, $2)
RETURNING id, client_id, developer_id
`

type CreateWorkingRelationshipParams struct {
	ClientID    pgtype.Int4
	DeveloperID pgtype.Int4
}

func (q *Queries) CreateWorkingRelationship(ctx context.Context, arg CreateWorkingRelationshipParams) (WorkingRelationship, error) {
	row := q.db.QueryRow(ctx, createWorkingRelationship, arg.ClientID, arg.DeveloperID)
	var i WorkingRelationship
	err := row.Scan(&i.ID, &i.ClientID, &i.DeveloperID)
	return i, err
}

const deleteClient = `-- name: DeleteClient :exec
DELETE FROM clients
WHERE id = $1
`

func (q *Queries) DeleteClient(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteClient, id)
	return err
}

const deleteDeliverable = `-- name: DeleteDeliverable :exec
DELETE FROM deliverables
WHERE id = $1
`

func (q *Queries) DeleteDeliverable(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteDeliverable, id)
	return err
}

const deleteDeveloper = `-- name: DeleteDeveloper :exec
DELETE FROM developers
WHERE id = $1
`

func (q *Queries) DeleteDeveloper(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteDeveloper, id)
	return err
}

const deletePercentComplete = `-- name: DeletePercentComplete :exec
DELETE FROM percent_complete
WHERE id = $1
`

func (q *Queries) DeletePercentComplete(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deletePercentComplete, id)
	return err
}

const deleteProject = `-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1
`

func (q *Queries) DeleteProject(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteProject, id)
	return err
}

const deleteTimesheet = `-- name: DeleteTimesheet :exec
DELETE FROM timesheet
WHERE id = $1
`

func (q *Queries) DeleteTimesheet(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteTimesheet, id)
	return err
}

const deleteWorkingRelationship = `-- name: DeleteWorkingRelationship :exec
DELETE FROM working_relationships
WHERE client_id = $1 AND developer_id = $2
`

type DeleteWorkingRelationshipParams struct {
	ClientID    pgtype.Int4
	DeveloperID pgtype.Int4
}

func (q *Queries) DeleteWorkingRelationship(ctx context.Context, arg DeleteWorkingRelationshipParams) error {
	_, err := q.db.Exec(ctx, deleteWorkingRelationship, arg.ClientID, arg.DeveloperID)
	return err
}

const earnedValue = `-- name: EarnedValue :many
SELECT
    d.project_id ProjectId,
    d.description Deliverable,
    d.type Type,
    percent_complete PercentComplete,
    budget Budget,
    to_char(reporting_period, 'YYYY-MM') ReportingPeriod,
    percent_complete * budget EarnedValue
FROM
    percent_complete pc
    JOIN
    deliverables d ON pc.deliv_id = d.id
WHERE
    d.project_id = $1
`

type EarnedValueRow struct {
	Projectid       pgtype.Int4
	Deliverable     pgtype.Text
	Type            interface{}
	Percentcomplete pgtype.Float8
	Budget          pgtype.Float8
	Reportingperiod string
	Earnedvalue     int32
}

func (q *Queries) EarnedValue(ctx context.Context, projectID pgtype.Int4) ([]EarnedValueRow, error) {
	rows, err := q.db.Query(ctx, earnedValue, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EarnedValueRow
	for rows.Next() {
		var i EarnedValueRow
		if err := rows.Scan(
			&i.Projectid,
			&i.Deliverable,
			&i.Type,
			&i.Percentcomplete,
			&i.Budget,
			&i.Reportingperiod,
			&i.Earnedvalue,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getClient = `-- name: GetClient :one
SELECT id, username, firstname, lastname, created_at
FROM clients
WHERE id = $1
`

func (q *Queries) GetClient(ctx context.Context, id int32) (Client, error) {
	row := q.db.QueryRow(ctx, getClient, id)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.CreatedAt,
	)
	return i, err
}

const getClients = `-- name: GetClients :many
SELECT id, username, firstname, lastname, created_at
FROM clients
`

func (q *Queries) GetClients(ctx context.Context) ([]Client, error) {
	rows, err := q.db.Query(ctx, getClients)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Client
	for rows.Next() {
		var i Client
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Firstname,
			&i.Lastname,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeveloper = `-- name: GetDeveloper :one
SELECT id, username, firstname, lastname, role, created_at
FROM developers
WHERE id = $1
`

func (q *Queries) GetDeveloper(ctx context.Context, id int32) (Developer, error) {
	row := q.db.QueryRow(ctx, getDeveloper, id)
	var i Developer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const getDevelopers = `-- name: GetDevelopers :many
SELECT id, username, firstname, lastname, role, created_at
FROM developers
`

func (q *Queries) GetDevelopers(ctx context.Context) ([]Developer, error) {
	rows, err := q.db.Query(ctx, getDevelopers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Developer
	for rows.Next() {
		var i Developer
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Firstname,
			&i.Lastname,
			&i.Role,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProjects = `-- name: GetProjects :many
SELECT id, title, description, client_id, dev_id, status, created_at
FROM projects
`

func (q *Queries) GetProjects(ctx context.Context) ([]Project, error) {
	rows, err := q.db.Query(ctx, getProjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.ClientID,
			&i.DevID,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkingRelationships = `-- name: GetWorkingRelationships :many
SELECT id, client_id, developer_id
FROM working_relationships
`

func (q *Queries) GetWorkingRelationships(ctx context.Context) ([]WorkingRelationship, error) {
	rows, err := q.db.Query(ctx, getWorkingRelationships)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WorkingRelationship
	for rows.Next() {
		var i WorkingRelationship
		if err := rows.Scan(&i.ID, &i.ClientID, &i.DeveloperID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateClient = `-- name: UpdateClient :one
UPDATE clients
SET
    username = $2,
    firstname = $3,
    lastname = $4
WHERE id = $1
RETURNING id, username, firstname, lastname, created_at
`

type UpdateClientParams struct {
	ID        int32
	Username  pgtype.Text
	Firstname pgtype.Text
	Lastname  pgtype.Text
}

func (q *Queries) UpdateClient(ctx context.Context, arg UpdateClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, updateClient,
		arg.ID,
		arg.Username,
		arg.Firstname,
		arg.Lastname,
	)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.CreatedAt,
	)
	return i, err
}

const updateDeliverable = `-- name: UpdateDeliverable :one
UPDATE deliverables
SET
    project_id = $2,
    description = $3,
    type = $4,
    budget = $5
WHERE id = $1
RETURNING id, project_id, description, type, budget, created_at
`

type UpdateDeliverableParams struct {
	ID          int32
	ProjectID   pgtype.Int4
	Description pgtype.Text
	Type        interface{}
	Budget      pgtype.Float8
}

func (q *Queries) UpdateDeliverable(ctx context.Context, arg UpdateDeliverableParams) (Deliverable, error) {
	row := q.db.QueryRow(ctx, updateDeliverable,
		arg.ID,
		arg.ProjectID,
		arg.Description,
		arg.Type,
		arg.Budget,
	)
	var i Deliverable
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Description,
		&i.Type,
		&i.Budget,
		&i.CreatedAt,
	)
	return i, err
}

const updateDeveloper = `-- name: UpdateDeveloper :one
UPDATE developers
SET
    username = $2,
    firstname = $3,
    lastname = $4,
    role = $5
WHERE id = $1
RETURNING id, username, firstname, lastname, role, created_at
`

type UpdateDeveloperParams struct {
	ID        int32
	Username  pgtype.Text
	Firstname pgtype.Text
	Lastname  pgtype.Text
	Role      pgtype.Text
}

func (q *Queries) UpdateDeveloper(ctx context.Context, arg UpdateDeveloperParams) (Developer, error) {
	row := q.db.QueryRow(ctx, updateDeveloper,
		arg.ID,
		arg.Username,
		arg.Firstname,
		arg.Lastname,
		arg.Role,
	)
	var i Developer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const updatePercentComplete = `-- name: UpdatePercentComplete :one
UPDATE percent_complete
SET
    percent_complete = $2,
    deliv_id = $3,
    report_period = $4
WHERE id = $1
RETURNING id, percent_complete, deliv_id, report_period, created_at
`

type UpdatePercentCompleteParams struct {
	ID              int32
	PercentComplete pgtype.Float8
	DelivID         pgtype.Int4
	ReportPeriod    pgtype.Date
}

func (q *Queries) UpdatePercentComplete(ctx context.Context, arg UpdatePercentCompleteParams) (PercentComplete, error) {
	row := q.db.QueryRow(ctx, updatePercentComplete,
		arg.ID,
		arg.PercentComplete,
		arg.DelivID,
		arg.ReportPeriod,
	)
	var i PercentComplete
	err := row.Scan(
		&i.ID,
		&i.PercentComplete,
		&i.DelivID,
		&i.ReportPeriod,
		&i.CreatedAt,
	)
	return i, err
}

const updateProject = `-- name: UpdateProject :one
UPDATE projects
SET
    title = $2,
    description = $3,
    client_id = $4,
    dev_id = $5,
    status = $6
WHERE id = $1
RETURNING id, title, description, client_id, dev_id, status, created_at
`

type UpdateProjectParams struct {
	ID          int32
	Title       pgtype.Text
	Description pgtype.Text
	ClientID    pgtype.Int4
	DevID       pgtype.Int4
	Status      pgtype.Text
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, updateProject,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.ClientID,
		arg.DevID,
		arg.Status,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.ClientID,
		&i.DevID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const updateTimesheet = `-- name: UpdateTimesheet :one
UPDATE timesheet
SET
    hours = $2,
    deliv_id = $3,
    date_of_work = $4
WHERE id = $1
RETURNING id, hours, deliv_id, date_of_work, created_at
`

type UpdateTimesheetParams struct {
	ID         int32
	Hours      pgtype.Float8
	DelivID    pgtype.Int4
	DateOfWork pgtype.Date
}

func (q *Queries) UpdateTimesheet(ctx context.Context, arg UpdateTimesheetParams) (Timesheet, error) {
	row := q.db.QueryRow(ctx, updateTimesheet,
		arg.ID,
		arg.Hours,
		arg.DelivID,
		arg.DateOfWork,
	)
	var i Timesheet
	err := row.Scan(
		&i.ID,
		&i.Hours,
		&i.DelivID,
		&i.DateOfWork,
		&i.CreatedAt,
	)
	return i, err
}
