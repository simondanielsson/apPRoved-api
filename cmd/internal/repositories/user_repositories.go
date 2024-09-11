package repositories

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(tx *gorm.DB, user *models.User) (uint, error) {
	if err := tx.Create(user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepository) GetUsers(tx *gorm.DB) ([]models.User, error) {
	var users []models.User

	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUser(tx *gorm.DB, userID uint) (*models.User, error) {
	var user models.User

	if err := tx.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(tx *gorm.DB, username string) (*models.User, error) {
	var user models.User

	if err := tx.First(&user, &models.User{Username: username}).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
