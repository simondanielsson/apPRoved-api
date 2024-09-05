package repositories

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/gorm"
)

type ReviewsRepository struct{}

func NewReviewsRepository(db *gorm.DB) *ReviewsRepository {
	return &ReviewsRepository{}
}

func (r *ReviewsRepository) GetRepositories(tx *gorm.DB, userID uint) ([]string, error) {
	var repositoryNames []string

	err := tx.Model(&models.Repository{}).Where(&models.Repository{UserID: userID}).Pluck("name", &repositoryNames).Error
	if err != nil {
		return nil, err
	}

	return repositoryNames, nil
}

func (r *ReviewsRepository) CreateRepository(tx *gorm.DB, repo *models.Repository) (*models.Repository, error) {
	if err := tx.Create(repo).Error; err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *ReviewsRepository) CreatePullRequests(tx *gorm.DB, prs []*models.PullRequest) error {
	return tx.CreateInBatches(prs, 30).Error
}
