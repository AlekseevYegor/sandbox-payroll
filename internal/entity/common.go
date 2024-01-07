package entity

import "time"

type (
	CreatedUpdated struct {
		UpdatedAt *time.Time `json:"-" db:"updated_at"`
		CreatedAt *time.Time `json:"-" db:"created_at"`
		CreatedBy string     `json:"-" db:"created_by"`
		UpdatedBy string     `json:"-" db:"updated_by"`
	}
)
