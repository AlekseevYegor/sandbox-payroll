package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"se-challenge-payroll/internal/entity/timetracking"
	"se-challenge-payroll/pkg/log"
)

type timeTrackerRepo struct {
	db *sqlx.DB
}

func NewTimeTrackerRepo(db *sqlx.DB) timetracking.Repository {
	return &timeTrackerRepo{
		db: db,
	}
}

func (t *timeTrackerRepo) BulkInsert(ctx context.Context, tracks []timetracking.TimeTracking) error {
	q := `INSERT INTO time_tracking (
				created_at,
				updated_at,
				created_by,
				updated_by,
                date,
			    biweekly_id,
			   	worked_hours,
			    time_report_id, 
			    employee_id,   
			    job_group_id
				)
			VALUES (:created_at, :updated_at, :created_by, :updated_by, :date, :biweekly_id, :worked_hours, :time_report_id, :employee_id, :job_group_id)`
	_, err := t.db.NamedExecContext(ctx, q, tracks)
	if err != nil {
		log.ZL.Error().
			Err(err).
			Msg("failed to create time tracking")
		return err
	}

	return nil
}

func (t *timeTrackerRepo) GetByFilter(ctx context.Context, filter timetracking.TimeTrackingFilter) ([]timetracking.TimeTracking, error) {
	var ts = make([]timetracking.TimeTracking, 0)

	q := `SELECT 
    			id,
                date,
                biweekly_id,
                worked_hours,
			    time_report_id, 
			    employee_id,   
			    job_group_id,
			    created_at,
				updated_at,
				created_by,
				updated_by
			FROM time_tracking
			WHERE (cast($1 as integer) IS NULL or time_report_id = $1)
			and (cast($2 as varchar) IS NULL or employee_id = $2)`

	rows, err := t.db.QueryContext(ctx, q,
		filter.TimeReportID,
		filter.EmployeeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ts, nil
		}

		log.ZL.Error().
			Err(err).
			Msg("failed to find time tracking")
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.ZL.Error().Err(err)
		}
	}()

	for rows.Next() {
		var t timetracking.TimeTracking
		err := rows.Scan(&t.ID,
			&t.Date,
			&t.BiWeekID,
			&t.WorkedHours,
			&t.TimeReportID,
			&t.EmployeeID,
			&t.JobGroupID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.CreatedBy,
			&t.UpdatedBy,
		)
		if err != nil {
			log.ZL.Error().
				Err(err).
				Msg("failed to Scan time tracking")
			return nil, err
		}

		ts = append(ts, t)
	}

	return ts, nil
}
