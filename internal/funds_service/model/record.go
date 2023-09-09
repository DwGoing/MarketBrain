package model

import "time"

type Record struct {
	Id        string    `gorm:"column:ID"`
	CreatedAt time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT"`
}
