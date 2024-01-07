package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"se-challenge-payroll/internal/entity/payroll"
	"se-challenge-payroll/pkg/log"
)

type payrollRepo struct {
	db *sqlx.DB
}

func NewPayrollRepo(db *sqlx.DB) payroll.Repository {
	return &payrollRepo{
		db: db,
	}
}

func (p *payrollRepo) GetReport(ctx context.Context) ([]payroll.ReportRow, error) {
	var r = make([]payroll.ReportRow, 0)

	q := `select biweekly_id, EXTRACT(YEAR FROM date), employee_id, sum(pr.rate*tt.worked_hours)::numeric(10,2)
		from time_tracking tt
		join payment_rate pr on tt.job_group_id = pr.job_group_id
		group by biweekly_id,EXTRACT(YEAR FROM date), employee_id;`

	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r, nil
		}

		log.ZL.Error().
			Err(err).
			Msg("failed to get payroll report")
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.ZL.Error().Err(err)
		}
	}()

	for rows.Next() {
		var rr payroll.ReportRow
		err := rows.Scan(&rr.BiWeekID,
			&rr.Year,
			&rr.EmployeeID,
			&rr.PayrollAmount,
		)
		if err != nil {
			log.ZL.Error().
				Err(err).
				Msg("failed to Scan time tracking")
			return nil, err
		}

		r = append(r, rr)
	}

	return r, nil
}
