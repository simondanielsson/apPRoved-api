package repositories

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/gorm"
)

type ReviewsRepository struct {
	db *gorm.DB
}

func NewReviewsRepository(db *gorm.DB) *ReviewsRepository {
	return &ReviewsRepository{db: db}
}

func (r *ReviewsRepository) GetRepositories(userID uint) ([]string, error) {
	var repositoryNames []string

	err := r.db.Model(&models.Repository{}).Where(&models.Repository{UserID: userID}).Pluck("name", &repositoryNames).Error
	if err != nil {
		return nil, err
	}

	return repositoryNames, nil
}

func (r *ReviewsRepository) CreateRepository(repo *models.Repository) (uint, error) {
	if err := r.db.Create(repo).Error; err != nil {
		return 0, err
	}

	return repo.ID, nil
}
