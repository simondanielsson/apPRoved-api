package services

import (
	"github.com/asaskevich/govalidator"
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
	customerrors "github.com/simondanielsson/apPRoved/pkg/custom_errors"
	"github.com/simondanielsson/apPRoved/pkg/utils"
	"gorm.io/gorm"
)

// TODO: check if we can define an interface for the service.
type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(tx *gorm.DB, username, email, password string) (uint, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	if !govalidator.IsEmail(email) {
		return 0, customerrors.NewValidationError("email", "Invalid email")
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	userID, err := s.userRepository.CreateUser(tx, &user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *UserService) GetUsers(tx *gorm.DB) ([]models.User, error) {
	users, err := s.userRepository.GetUsers(tx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUser(tx *gorm.DB, userID uint) (*models.User, error) {
	user, err := s.userRepository.GetUser(tx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
