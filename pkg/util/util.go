package util

import (
	"fmt"
	"time"
)

var biWeeklyMap = make(map[string]struct{ dateFrom, dateTo time.Time })

// ContainsString - return true if `in` contains `find`. False - otherwise
func ContainsString(find string, in ...string) bool {
	for _, ct := range in {
		if ct == find {
			return true
		}
	}

	return false
}

func StringToDate(dateString string) (time.Time, error) {
	layout := "2/1/2006"
	return time.Parse(layout, dateString)
}

func DateToBiweeklyID(date time.Time) int {
	biweeklyID := int(date.Month()*2) - 1

	if date.Day() > 15 {
		biweeklyID++
	}

	return biweeklyID
}

func BiweeklyPaymentDate(periodNum, year int) (time.Time, time.Time) {
	var dateFrom, dateTo time.Time

	if timePeriod, ok := biWeeklyMap[fmt.Sprintf("%d%d", year, periodNum)]; ok {
		return timePeriod.dateFrom, timePeriod.dateTo
	}

	month := time.Month(periodNum / 2)

	dateFrom = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	dateTo = time.Date(year, month, 15, 0, 0, 0, 0, time.UTC)

	if periodNum%2 > 0 {
		dateFrom = time.Date(year, month, 16, 0, 0, 0, 0, time.UTC)
		dateTo = getLastDayOfMonth(year, month)
	}

	biWeeklyMap[fmt.Sprintf("%d%d", year, periodNum)] = struct{ dateFrom, dateTo time.Time }{dateFrom: dateFrom, dateTo: dateTo}

	return dateFrom, dateTo
}

func getLastDayOfMonth(year int, month time.Month) time.Time {
	// Get the first day of the next month
	firstDayOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	// Subtract one day to get the last day of the current month
	lastDayOfMonth := firstDayOfNextMonth.Add(-time.Hour * 24)

	return lastDayOfMonth
}
