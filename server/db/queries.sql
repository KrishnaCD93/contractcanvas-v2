-- name: GetUser: Get user by id
SELECT *
FROM users
WHERE id = $1;

-- projects: GetProjectsByDevId: Get projects by developer id
SELECT *
FROM projects
WHERE dev_id = $1;

-- projects: GetProjectById: Get project by id
SELECT *
FROM projects
WHERE id = $1;

-- projects: GetDeliverablesByProjectId: Get deliverables by project id
SELECT *
FROM deliverables
WHERE project_id = $1;

-- projects: GetPercentCompleteByDelivId: Get percent complete by deliverable id
SELECT *
FROM percent_complete
WHERE deliv_id = $1;

-- projects: GetTimesheetByDelivId: Get timesheet by deliverable id
SELECT *
FROM timesheet
WHERE deliv_id = $1;

-- projects; GetEearnedValueByProjectId: Get earned value by project id
SELECT *
FROM EarnedValue
WHERE project_id = $1;

-- projects; GetBurnedValueByProjectId: Get burned value by project id
SELECT *
FROM BurnedValue
WHERE project_id = $1;

-- projects; GetCPIByProjectId: Get CPI by project id
SELECT *
FROM CPI
WHERE ProjectId = $1;