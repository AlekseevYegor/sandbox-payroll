CREATE TABLE IF NOT EXISTS payment_rate (
    id BIGSERIAL PRIMARY KEY,
    created_at                      TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME zone 'UTC') NOT NULL,
    updated_at                      TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME zone 'UTC') NOT NULL,
    created_by                      VARCHAR(64)                                                    NOT NULL,
    updated_by                      VARCHAR(64) NOT NULL,
    job_group_id varchar(64) NOT NULL UNIQUE,
    rate numeric(10,2) NOT NULL
);


INSERT INTO payment_rate (created_at, updated_at, created_by, updated_by, job_group_id, rate)
VALUES (now() at time zone 'UTC', now() at time zone 'UTC', 'migration_1', 'migration_1', 'A', 20.0),
        (now() at time zone 'UTC', now() at time zone 'UTC', 'migration_1', 'migration_1', 'B', 30.0);


CREATE TABLE IF NOT EXISTS time_tracking(
    id BIGSERIAL PRIMARY KEY,
    created_at                      TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME zone 'UTC') NOT NULL,
    updated_at                      TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME zone 'UTC') NOT NULL,
    created_by                      VARCHAR(64)                                                    NOT NULL,
    updated_by                      VARCHAR(64) NOT NULL,
    date DATE NOT NULL,
    worked_hours NUMERIC(10,2) NOT NULL,
    time_report_id INTEGER NOT NULL,
    biweekly_id INTEGER NOT NULL,
    employee_id VARCHAR(64) NOT NULL, -- should have references to employee table
    job_group_id VARCHAR(64) REFERENCES payment_rate(job_group_id) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_time_tracking_time_report_id on time_tracking(time_report_id);
CREATE INDEX IF NOT EXISTS idx_time_tracking_employee_id on time_tracking(employee_id);

