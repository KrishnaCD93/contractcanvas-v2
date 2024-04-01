-- name: GetDeveloper :one
SELECT *
FROM developers
WHERE id = $1;

-- name: GetClient :one
SELECT *
FROM clients
WHERE id = $1;

-- projects: GetProjectsByDevId :many
SELECT *
FROM projects
WHERE dev_id = $1;

-- projects: GetProjectsByClientId :many
SELECT *
FROM projects
WHERE client_id = $1;

-- projects: GetProjectById :one
SELECT *
FROM projects
WHERE id = $1;

-- projects: GetDeliverablesByProjectId: many
SELECT *
FROM deliverables
WHERE project_id = $1;

-- projects: GetPercentCompleteByDelivId :many
SELECT *
FROM percent_complete
WHERE deliv_id = $1;

-- projects: GetTimesheetByDelivId :many
SELECT *
FROM timesheet
WHERE deliv_id = $1;

-- name: EarnedValue :many
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
    d.project_id = $1;

-- name: BurnedValue :many
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
    d.project_id = $1;

-- name: CostPerformanceIndex :many
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
    WHERE d.project_id = $1) BurnedValue ON EarnedValue.Deliverable = BurnedValue.Deliverable AND EarnedValue.ReportingPeriod = BurnedValue.ReportingPeriod;

-- name: CreateClient :one
INSERT INTO clients
    (username, firstname, lastname, email, bio)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateDeveloper :one
INSERT INTO developers
    (username, firstname, lastname, role, email, bio)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateWorkingRelationship :one
INSERT INTO working_relationships
    (client_id, developer_id)
VALUES
    ($1, $2)
RETURNING *;

-- name: CreateProject :one
INSERT INTO projects
    (title, description, client_id, dev_id, status)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateDeliverable :one
INSERT INTO deliverables
    (project_id, description, type, budget)
VALUES
    ($1, $2, $3, $4)
RETURNING *;

-- name: CreatePercentComplete :one
INSERT INTO percent_complete
    (percent_complete, deliv_id, report_period)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: CreateTimesheet :one
INSERT INTO timesheet
    (hours, deliv_id, date_of_work)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: DeleteDeveloper :exec
DELETE FROM developers
WHERE id = $1;

-- name: DeleteClient :exec
DELETE FROM clients
WHERE id = $1;

-- name: DeleteWorkingRelationship :exec
DELETE FROM working_relationships
WHERE client_id = $1 AND developer_id = $2;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: DeleteDeliverable :exec
DELETE FROM deliverables
WHERE id = $1;

-- name: DeletePercentComplete :exec
DELETE FROM percent_complete
WHERE id = $1;

-- name: DeleteTimesheet :exec
DELETE FROM timesheet
WHERE id = $1;

-- name: UpdateClient :one
UPDATE clients
SET
    username = $2,
    firstname = $3,
    lastname = $4,
    email = $5,
    bio = $6
WHERE id = $1
RETURNING *;

-- name: UpdateDeveloper :one
UPDATE developers
SET
    username = $2,
    firstname = $3,
    lastname = $4,
    role = $5,
    email = $6,
    bio = $7
WHERE id = $1
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
    title = $2,
    description = $3,
    client_id = $4,
    dev_id = $5,
    status = $6
WHERE id = $1
RETURNING *;

-- name: UpdateDeliverable :one
UPDATE deliverables
SET
    project_id = $2,
    description = $3,
    type = $4,
    budget = $5
WHERE id = $1
RETURNING *;

-- name: UpdatePercentComplete :one
UPDATE percent_complete
SET
    percent_complete = $2,
    deliv_id = $3,
    report_period = $4
WHERE id = $1
RETURNING *;

-- name: UpdateTimesheet :one
UPDATE timesheet
SET
    hours = $2,
    deliv_id = $3,
    date_of_work = $4
WHERE id = $1
RETURNING *;

-- name: GetClients :many
SELECT *
FROM clients;

-- name: GetDevelopers :many
SELECT *
FROM developers;

-- name: GetWorkingRelationships :many
SELECT *
FROM working_relationships;

-- name: GetProjects :many
SELECT *
FROM projects;
