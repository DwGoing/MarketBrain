package model

import "time"

type Record struct {
	Id        string    `gorm:"column:ID"`
	CreatedAt time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT"`
}

type DeleteOption struct {
	Conditions           string
	ConditionsParameters []any
}

type UpdateOption struct {
	Conditions           string
	ConditionsParameters []any
	Values               map[string]any
}

type GetOption struct {
	Conditions           string
	ConditionsParameters []any
	PageSize             int64
	PageIndex            int64
}
