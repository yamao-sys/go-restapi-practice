package repositories

import (
	"app/models"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(user *models.User, email string) error
	FindUserByID(id int) models.User
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(user *models.User) error {
	if err := ur.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FindUserByEmail(user *models.User, email string) error {
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FindUserByID(id int) models.User {
	user := models.User{}
	if err := ur.db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Fatalln(err)
	}
	return user
}
