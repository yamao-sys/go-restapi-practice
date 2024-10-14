package models

import "time"

type User struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
