package services

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
	"github.com/simondanielsson/apPRoved/pkg/utils"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

func (s *AuthService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := utils.ValidatePassword(password, user.Password); err != nil {
		return "", err
	}

	token, err := utils.CreateJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
