package repository

import (
	"context"
	"se-challenge-payroll/internal/entity/timetracking"
)

type TimeTrackerRepoMock struct {
	db map[int][]timetracking.TimeTracking
}

func NewTimeTrackerRepoMock() *TimeTrackerRepoMock {
	return &TimeTrackerRepoMock{
		db: make(map[int][]timetracking.TimeTracking),
	}
}

func (t *TimeTrackerRepoMock) BulkInsert(ctx context.Context, tracks []timetracking.TimeTracking) error {
	if len(tracks) > 0 {
		t.db[tracks[0].TimeReportID] = tracks
	}

	return nil
}

func (t *TimeTrackerRepoMock) GetByFilter(ctx context.Context, filter timetracking.TimeTrackingFilter) ([]timetracking.TimeTracking, error) {
	if filter.TimeReportID != nil {
		tracks, ok := t.db[*filter.TimeReportID]
		if !ok {
			return make([]timetracking.TimeTracking, 0), nil
		}

		return tracks, nil
	}

	return make([]timetracking.TimeTracking, 0), nil
}
