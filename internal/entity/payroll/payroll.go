package payroll

import (
	"context"
)

type (
	ReportRow struct {
		BiWeekID      int     `db:"biweekly_id"`
		PayrollAmount float64 `db:"payroll_amount"`
		Year          int     `db:"year"`
		EmployeeID    string  `db:"employee_id"`
	}
)

type Repository interface {
	GetReport(ctx context.Context) ([]ReportRow, error)
}
