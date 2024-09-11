package services

import (
	"context"
	"log"

	"github.com/simondanielsson/apPRoved/cmd/config"
	"github.com/simondanielsson/apPRoved/cmd/internal/constants"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/requests"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/responses"
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
	"github.com/simondanielsson/apPRoved/pkg/utils"
	"github.com/simondanielsson/apPRoved/pkg/utils/mq"
	"gorm.io/gorm"
)

type ReviewsService struct {
	reviewsRepository *repositories.ReviewsRepository
}

// NewReviewsService creates a new reviews service
func NewReviewsService(reviewsRepository *repositories.ReviewsRepository) *ReviewsService {
	return &ReviewsService{reviewsRepository: reviewsRepository}
}

// GetRepositories returns all repositories for a user
func (rs *ReviewsService) GetRepositories(tx *gorm.DB, userID uint) ([]*responses.GetRepositoriesResponse, error) {
	repos, err := rs.reviewsRepository.GetRepositories(tx, userID)
	if err != nil {
		return nil, err
	}

	var reposResponse []*responses.GetRepositoriesResponse
	for _, repo := range repos {
		reposResponse = append(reposResponse, &responses.GetRepositoriesResponse{
			ID:        repo.ID,
			Name:      repo.Name,
			Owner:     repo.Owner,
			URL:       repo.URL,
			CreatedAt: repo.CreatedAt,
			UpdatedAt: repo.UpdatedAt,
		})
	}

	return reposResponse, nil
}

// RegisterRepository registers a new repository and its pull requests
func (rs *ReviewsService) RegisterRepository(ctx context.Context, tx *gorm.DB, userID uint, name, owner, url string) (uint, error) {
	repo := &models.Repository{
		UserID: userID,
		Name:   name,
		Owner:  owner,
		URL:    url,
	}
	repo, err := rs.reviewsRepository.CreateRepository(tx, repo)
	if err != nil {
		return 0, err
	}

	prs, err := rs.findPullRequests(ctx, repo, userID)
	if err != nil {
		return 0, err
	}

	if err := rs.reviewsRepository.CreatePullRequests(tx, prs); err != nil {
		return 0, err
	}
	return repo.ID, nil
}

func (rs *ReviewsService) findPullRequests(ctx context.Context, repo *models.Repository, userID uint) ([]*models.PullRequest, error) {
	fetched_prs, err := utils.ListPullRequests(ctx, repo.Name, repo.Owner, userID)
	if err != nil {
		return nil, err
	}

	var prs []*models.PullRequest
	for _, pr := range fetched_prs {
		prs = append(prs, &models.PullRequest{
			RepositoryID: repo.ID,
			Number:       pr.Number,
			Title:        pr.Title,
			URL:          pr.URL,
			State:        pr.State,
			LastCommit:   pr.LastCommit,
		})
	}

	return prs, nil
}

// GetPullRequests returns all pull requests for a repository
func (rs *ReviewsService) GetPullRequests(tx *gorm.DB, userID, repoID uint) ([]*responses.GetPullRequestResponse, error) {
	prs, err := rs.reviewsRepository.GetPullRequests(tx, userID, repoID)
	if err != nil {
		return nil, err
	}

	var prsResponse []*responses.GetPullRequestResponse
	for _, pr := range prs {
		prsResponse = append(prsResponse, &responses.GetPullRequestResponse{
			ID:        pr.ID,
			Number:    pr.Number,
			Title:     pr.Title,
			URL:       pr.URL,
			State:     pr.State,
			CreatedAt: pr.CreatedAt,
			UpdatedAt: pr.UpdatedAt,
		})
	}

	return prsResponse, nil
}

// GetReviews returns all reviews for a pull request
func (rs *ReviewsService) GetReviews(tx *gorm.DB, repoID, prID uint) ([]*responses.GetReviewsResponse, error) {
	reviews, err := rs.reviewsRepository.GetReviews(tx, repoID, prID)
	if err != nil {
		return nil, err
	}

	var reviewResponse []*responses.GetReviewsResponse
	for _, review := range reviews {
		reviewResponse = append(reviewResponse, &responses.GetReviewsResponse{
			ID:        review.ID,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
		})
	}

	return reviewResponse, nil
}

// GetReview returns a review
func (rs *ReviewsService) GetReview(tx *gorm.DB, reviewID uint) (*responses.GetReviewResponse, error) {
	review, err := rs.reviewsRepository.GetReview(tx, reviewID)
	if err != nil {
		return nil, err
	}
	response := &responses.GetReviewResponse{
		ID:          review.ID,
		FileReviews: []*responses.GetFileReviewResponse{},
		CreatedAt:   review.CreatedAt,
		UpdatedAt:   review.UpdatedAt,
	}
	for _, fr := range review.FileReviews {
		fileReviewResponse := &responses.GetFileReviewResponse{
			ID:        fr.ID,
			Filename:  fr.Filename,
			Content:   fr.Content,
			CreatedAt: fr.CreatedAt,
			UpdatedAt: fr.UpdatedAt,
		}
		response.FileReviews = append(response.FileReviews, fileReviewResponse)
	}

	return response, nil
}

func (rs *ReviewsService) CreateReview(tx *gorm.DB, ctx context.Context, repoID, prID uint, name string, userID uint) (uint, error) {
	repo, err := rs.reviewsRepository.GetRepository(tx, repoID)
	if err != nil {
		return 0, err
	}

	pr, err := rs.reviewsRepository.GetPullRequest(tx, prID)
	if err != nil {
		return 0, nil
	}

	review := &models.Review{
		Name:          name,
		PullRequestID: pr.ID,
	}
	review, err = rs.reviewsRepository.CreateReview(tx, review)
	if err != nil {
		return 0, nil
	}

	reviewStatus := &models.ReviewStatus{
		ReviewID: review.ID,
		Status:   constants.StatusQueued,
	}
	if err := rs.reviewsRepository.CreateReviewStatus(tx, reviewStatus); err != nil {
		return 0, err
	}

	// fetch file diffs for the PR from github using github client
	// send info over RabbitMQ to call external review service api to retrieve file reviews
	go func() {
		// fetch file diffs for the PR from GitHub using the GitHub client
		fileDiffs, err := utils.FetchFileDiffs(ctx, repo.Name, repo.Owner, pr.Number, userID)
		if err != nil {
			log.Println("Error fetching file diffs:", err)
			return
		}

		// Publish a message to RabbitMQ
		message := requests.FileDiffReviewRequest{
			FileDiffs:      fileDiffs,
			ReviewID:       review.ID,
			ReviewStatusID: reviewStatus.ID,
		}
		err = mq.Publish(config.QueueFileDiffs, &message)
		if err != nil {
			// Handle the error, e.g., log or retry
			log.Println("Error publishing to RabbitMQ:", err)
			return
		}
	}()

	// return the review id
	return review.ID, nil
}

func (rs *ReviewsService) CompleteReview(tx *gorm.DB, req *requests.CompleteReviewRequest) error {
	var fileReviews []*models.FileReview
	for _, review := range req.FileReviews {
		fr := models.FileReview{
			ReviewID: req.ReviewID,
			Filename: review.Filename,
			Content:  review.Content,
		}
		fileReviews = append(fileReviews, &fr)
	}

	if err := rs.reviewsRepository.CreateFileReviews(tx, fileReviews); err != nil {
		return err
	}

	if err := rs.reviewsRepository.UpdateReviewStatus(tx, req.ReviewStatusID, constants.StatusAvailable); err != nil {
		return err
	}

	if err := rs.reviewsRepository.UpdateReviewProgress(tx, req.ReviewStatusID, 100); err != nil {
		return err
	}

	return nil
}
