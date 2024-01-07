package timetracking

import (
	"context"
	"fmt"
	"se-challenge-payroll/internal/entity"
	"se-challenge-payroll/pkg/util"
	"strconv"
	"strings"
	"time"
)

type (
	TimeTracking struct {
		ID           int64     `db:"id"`
		TimeReportID int       `db:"time_report_id"`
		Date         time.Time `db:"date"`
		BiWeekID     int       `db:"biweekly_id"`
		WorkedHours  float64   `db:"worked_hours"`
		JobGroupID   string    `db:"job_group_id"`
		EmployeeID   string    `db:"employee_id"`

		entity.CreatedUpdated
	}

	TimeTrackingFilter struct {
		TimeReportID *int
		EmployeeID   *string
	}
)

type Service interface {
	BulkCreate(ctx context.Context, tracking []TimeTracking) error
	GetByFilter(ctx context.Context, filter TimeTrackingFilter) ([]TimeTracking, error)
}

type Repository interface {
	BulkInsert(ctx context.Context, tracking []TimeTracking) error
	GetByFilter(ctx context.Context, filter TimeTrackingFilter) ([]TimeTracking, error)
}

// CreateTimeTrackingList - create tTimeTracking list from scv data
func CreateTimeTrackingList(data [][]string, reportID int) ([]TimeTracking, []error) {
	var (
		tts     = make([]TimeTracking, 0)
		errs    = make([]error, 0)
		now     = time.Now().UTC()
		created = entity.CreatedUpdated{
			CreatedAt: &now,
			UpdatedAt: &now,
			CreatedBy: "csv download",
			UpdatedBy: "csv download",
		}
	)

	for i, line := range data {
		if i > 0 { //omit header line
			rec := TimeTracking{TimeReportID: reportID, CreatedUpdated: created}

			for j, field := range line {
				switch j {
				case 0:
					date, err := util.StringToDate(field)
					if err != nil {
						errs = append(errs, fmt.Errorf("row %d contain incorrect data format", i+1))
						continue
					}
					rec.Date = date
				case 1:
					workHours, err := strconv.ParseFloat(strings.TrimSpace(field), 64)
					if err != nil {
						errs = append(errs, fmt.Errorf("row %d contain incorrect working hours format", i+1))
						continue
					}
					rec.WorkedHours = workHours
				case 2:
					rec.EmployeeID = strings.TrimSpace(field)
				case 3:
					rec.JobGroupID = strings.TrimSpace(field)
				}
			}

			if !rec.Date.IsZero() {
				rec.BiWeekID = util.DateToBiweeklyID(rec.Date)
			}

			tts = append(tts, rec)
		}
	}

	return tts, errs
}
