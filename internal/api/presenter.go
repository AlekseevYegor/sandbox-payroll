package api

import (
	"fmt"
	"se-challenge-payroll/internal/entity/payroll"
	"se-challenge-payroll/pkg/util"
	"sort"
	"time"
)

type (

	// PayrollReportResponse represents the payroll report structure
	//
	// swagger:model
	PayrollReportResponse struct {
		*PayrollReport `json:"payrollReport"`
	}

	// PayrollReport represents inner payroll report structure
	//
	// swagger:model
	PayrollReport struct {
		// list if employee report
		//
		// required: true
		EmployeeReports []EmployeeReport `json:"employeeReports"`
	}

	// EmployeeReport represents payroll report row structure
	//
	// swagger:model
	EmployeeReport struct {
		// EmployeeID...
		//
		// required: true
		EmployeeID string `json:"employeeId"`
		// AmountPaid - amount to pay in this period
		//
		// required: true
		AmountPaid string `json:"amountPaid"`
		// PayPeriod brackets of payment period
		//
		// required: true
		*PayPeriod `json:"payPeriod"`
	}

	// PayPeriod represents brackets of payment period
	//
	// swagger:model
	PayPeriod struct {
		// Started date
		//
		// required: true
		StartDate string `json:"startDate"`
		// End date
		//
		// required: true
		EndDate string `json:"endDate"`
	}
)

func NewEmployeeReport(row *payroll.ReportRow) EmployeeReport {
	dateFrom, dateTo := util.BiweeklyPaymentDate(row.BiWeekID, row.Year)

	return EmployeeReport{
		EmployeeID: row.EmployeeID,
		AmountPaid: fmt.Sprintf("$%.2f", row.PayrollAmount),
		PayPeriod: &PayPeriod{
			StartDate: dateFrom.Format(time.DateOnly),
			EndDate:   dateTo.Format(time.DateOnly),
		},
	}
}

func NewPayrollReportResponse(report []payroll.ReportRow) *PayrollReportResponse {
	var ers = make([]EmployeeReport, 0, len(report))
	for _, row := range report {
		ers = append(ers, NewEmployeeReport(&row))
	}

	// sort in ascending order by employee ID and period
	sort.Slice(ers, func(i, j int) bool {
		if ers[i].EmployeeID != ers[j].EmployeeID {
			return ers[i].EmployeeID < ers[j].EmployeeID
		}
		return ers[i].StartDate < ers[j].StartDate
	})

	return &PayrollReportResponse{PayrollReport: &PayrollReport{EmployeeReports: ers}}
}
