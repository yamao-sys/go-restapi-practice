package repositories

import (
	"app/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(user *models.User) {
	ur.db.Create(&user)
}
