package models

type Todo struct {
	ID      int    `gorm:"primary_key" json:"id"`
	Title   string `gorm:"size:255;not null" validate:"required"`
	Content string `gorm:"type:text"`
	UserID  int    `gorm:"not null" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" validate:"omitempty"`
}
