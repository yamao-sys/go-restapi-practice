package models

import "time"

type User struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	Email     string `gorm:"size:255;not null" validate:"required"`
	Password  string `gorm:"size:255;not null" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
