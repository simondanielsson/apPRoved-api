package services

import (
	"context"

	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
	"github.com/simondanielsson/apPRoved/pkg/utils"
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

func (rs *ReviewsService) RegisterRepository(ctx context.Context, userID uint, name, owner, url string) (uint, error) {
	repo := &models.Repository{
		UserID: userID,
		Name:   name,
		Owner:  owner,
		URL:    url,
	}
	repo, err := rs.reviewsRepository.CreateRepository(repo)
	if err != nil {
		return 0, err
	}

	prs, err := rs.findPullRequests(ctx, repo, userID)
	if err != nil {
		return 0, err
	}

	if err := rs.reviewsRepository.CreatePullRequests(prs); err != nil {
		return 0, err
	}
	return repo.ID, nil
}

func (rs *ReviewsService) findPullRequests(ctx context.Context, repo *models.Repository, userID uint) ([]*models.PullRequest, error) {
	fetched_prs, err := utils.ListPullRequests(ctx, repo.Name, repo.Owner, userID)
	if err != nil {
		return nil, err
	}

	prs := make([]*models.PullRequest, len(fetched_prs))
	for _, pr := range fetched_prs {
		prs = append(prs, &models.PullRequest{
			RepositoryID: repo.ID,
			Number:       pr.Number,
			Title:        pr.Title,
			URL:          pr.URL,
		})
	}

	return prs, nil
}
