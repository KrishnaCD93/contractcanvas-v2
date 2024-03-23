
CREATE TABLE "clients"
(
    "id" integer PRIMARY KEY,
    "works_with_user_id" integer,
    "submitted_project_id" integer,
    "created_at" timestamp DEFAULT current_timestamp
);

CREATE TABLE "developers"
(
    "id" integer PRIMARY KEY,
    "username" varchar,
    "firstname" varchar,
    "lastname" varchar,
    "role" varchar,
    "created_at" timestamp DEFAULT current_timestamp
);

CREATE TABLE "projects"
(
    "id" integer PRIMARY KEY,
    "title" varchar,
    "description" text,
    "dev_id" integer,
    "evm_id" integer,
    "procurement_id" integer,
    "staffing_id" integer,
    "status" varchar,
    "created_at" timestamp DEFAULT current_timestamp
);

CREATE TABLE "deliverables"
(
    "id" integer PRIMARY KEY,
    "project_id" integer,
    "description" text,
    "type" nvarchar,
    "budget" float,
    "created_at" timestamp DEFAULT current_timestamp
);

CREATE TABLE "percent_complete"
(
    "id" integer PRIMARY KEY,
    "deliv_id" integer,
    "percent_complete" float,
    "report_period" date,
    "created_at" timestamp DEFAULT current_timestamp
);

CREATE TABLE "burned_value"
(
    "id" integer PRIMARY KEY,
    "project_id" integer,
    "burned_hours" float,
    "deliv_id" integer,
    "date_of_work" date,
    "created_at" timestamp DEFAULT current_timestamp
);

ALTER TABLE "projects" ADD FOREIGN KEY ("dev_id") REFERENCES "developers" ("id");

ALTER TABLE "clients" ADD FOREIGN KEY ("works_with_user_id") REFERENCES "developers" ("id");

ALTER TABLE "projects" ADD FOREIGN KEY ("id") REFERENCES "developers" ("id");

ALTER TABLE "deliverables" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "percent_complete" ADD FOREIGN KEY ("deliv_id") REFERENCES "deliverables" ("id");

ALTER TABLE "burned_value" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "burned_value" ADD FOREIGN KEY ("deliv_id") REFERENCES "deliverables" ("id");
