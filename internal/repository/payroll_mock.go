package repository

import (
	"context"
	"math"
	"se-challenge-payroll/internal/entity/payroll"
	"se-challenge-payroll/internal/entity/timetracking"
)

type payrollRepoMock struct {
	rateDBMap        map[string]float64
	timeTrackingRepo timetracking.Repository
}

func NewPayrollRepoMock(ttr timetracking.Repository) payroll.Repository {
	return &payrollRepoMock{
		rateDBMap:        map[string]float64{"A": 20.0, "B": 30.0},
		timeTrackingRepo: ttr,
	}
}

func (p *payrollRepoMock) GetReport(ctx context.Context) ([]payroll.ReportRow, error) {
	var r = make([]payroll.ReportRow, 0)
	tt, err := p.timeTrackingRepo.GetByFilter(ctx, timetracking.TimeTrackingFilter{})
	if err != nil {
		return nil, err
	}

	for i := range tt {
		rate, ok := p.rateDBMap[tt[i].JobGroupID]
		if !ok {
			continue
		}

		rr := payroll.ReportRow{
			Year:          tt[i].Date.Year(),
			PayrollAmount: math.Round((tt[i].WorkedHours*rate)*100) / 100,
			EmployeeID:    tt[i].EmployeeID,
			BiWeekID:      tt[i].BiWeekID,
		}

		r = append(r, rr)
	}

	return r, nil
}
