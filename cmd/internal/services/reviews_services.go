package services

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
)

type ReviewsService struct {
	reviewsRepository *repositories.ReviewsRepository
}

func NewReviewsService(reviewsRepository *repositories.ReviewsRepository) *ReviewsService {
	return &ReviewsService{reviewsRepository: reviewsRepository}
}

func (rs *ReviewsService) GetRepositories(userID uint) ([]string, error) {
	return rs.reviewsRepository.GetRepositories(userID)
}

func (rs *ReviewsService) CreateRepository(userID uint, name, url string) (uint, error) {
	repo := &models.Repository{
		UserID: userID,
		Name:   name,
		URL:    url,
	}
	return rs.reviewsRepository.CreateRepository(repo)
}
