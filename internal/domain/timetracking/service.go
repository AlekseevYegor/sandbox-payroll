package timetracking

import (
	"context"
	"se-challenge-payroll/internal/entity/timetracking"
)

type service struct {
	repo timetracking.Repository
}

func New(repo timetracking.Repository) timetracking.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) BulkCreate(ctx context.Context, tracking []timetracking.TimeTracking) error {
	return s.repo.BulkInsert(ctx, tracking)
}

func (s *service) GetByFilter(ctx context.Context, filter timetracking.TimeTrackingFilter) ([]timetracking.TimeTracking, error) {
	return s.repo.GetByFilter(ctx, filter)
}
