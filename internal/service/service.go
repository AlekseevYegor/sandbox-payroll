package service

import (
	"se-challenge-payroll/internal/entity/payroll"
	"se-challenge-payroll/internal/entity/timetracking"
)

type (
	// Service ...
	Service struct {
		TimeTrackingService timetracking.Service
		PayrollRepository   payroll.Repository
	}
)

func NewService(ts timetracking.Service, pr payroll.Repository) Service {
	return Service{
		TimeTrackingService: ts,
		PayrollRepository:   pr,
	}
}
