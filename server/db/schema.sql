
CREATE TABLE "clients"
(
    "id" SERIAL PRIMARY KEY,
    "username" varchar,
    "firstname" varchar,
    "lastname" varchar,
    "created_at" timestamp default current_timestamp
);

CREATE TABLE "developers"
(
    "id" SERIAL PRIMARY KEY,
    "username" varchar,
    "firstname" varchar,
    "lastname" varchar,
    "role" varchar,
    "created_at" timestamp default current_timestamp
);

CREATE TABLE "working_relationships"
(
    "id" SERIAL PRIMARY KEY,
    "client_id" integer,
    "developer_id" integer
);

CREATE TABLE "projects"
(
    "id" SERIAL PRIMARY KEY,
    "title" varchar,
    "description" text,
    "client_id" integer,
    "dev_id" integer,
    "status" varchar,
    "created_at" timestamp default current_timestamp
);

CREATE TABLE "deliverables"
(
    "id" SERIAL PRIMARY KEY,
    "project_id" integer,
    "description" text,
    "type" nvarchar,
    "budget" float,
    "created_at" timestamp default current_timestamp
);

CREATE TABLE "percent_complete"
(
    "id" SERIAL PRIMARY KEY,
    "percent_complete" float,
    "deliv_id" integer,
    "report_period" date,
    "created_at" timestamp default current_timestamp
);

CREATE TABLE "timesheet"
(
    "id" SERIAL PRIMARY KEY,
    "hours" float,
    "deliv_id" integer,
    "date_of_work" date,
    "created_at" timestamp default current_timestamp
);

ALTER TABLE "projects" ADD FOREIGN KEY ("dev_id") REFERENCES "developers" ("id");

ALTER TABLE "working_relationships" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "working_relationships" ADD FOREIGN KEY ("developer_id") REFERENCES "developers" ("id");

ALTER TABLE "projects" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "deliverables" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "percent_complete" ADD FOREIGN KEY ("deliv_id") REFERENCES "deliverables" ("id");

ALTER TABLE "timesheet" ADD FOREIGN KEY ("deliv_id") REFERENCES "deliverables" ("id");

CREATE VIEW "EarnedValue"
AS
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
        deliverables d ON pc.deliv_id = d.id;

CREATE VIEW "BurnedValue"
AS
    SELECT
        d.project_id ProjectId,
        d.description Deliverable,
        d.type Type,
        to_char(date_of_work, 'YYYY-MM') ReportingPeriod,
        SUM(t.hours) BurnedValue
    FROM
        timesheet t
        JOIN
        deliverables d ON t.deliv_id = d.id;

CREATE VIEW "CPI"
AS
    SELECT
        e.ProjectId,
        e.Deliverable,
        e.Type,
        e.ReportingPeriod,
        e.EarnedValue,
        b.BurnedValue,
        e.EarnedValue / b.BurnedValue CPI
    FROM
        EarnedValue e
        JOIN
        BurnedValue b ON e.Deliverable = b.Deliverable AND e.ReportingPeriod = b.ReportingPeriod;