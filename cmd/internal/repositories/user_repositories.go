package repositories

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) (uint, error) {
	if err := r.db.Create(user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUser(userID uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, &models.User{Username: username}).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
