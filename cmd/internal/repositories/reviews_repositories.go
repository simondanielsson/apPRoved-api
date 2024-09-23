package repositories

import (
	"fmt"
	"log"

	"github.com/simondanielsson/apPRoved/cmd/constants"
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/gorm"
)

type ReviewsRepository struct{}

// NewReviewsRepository creates a new reviews repository
func NewReviewsRepository() *ReviewsRepository {
	return &ReviewsRepository{}
}

// GetRepositories returns all repositories for a user
func (r *ReviewsRepository) GetRepositories(tx *gorm.DB, userID uint) ([]*models.Repository, error) {
	var repos []*models.Repository

	if err := tx.Model(&models.Repository{}).Where(&models.Repository{UserID: userID}).Find(&repos).Error; err != nil {
		return nil, err
	}

	return repos, nil
}

// GetRepository returns a repository
func (r *ReviewsRepository) GetRepository(tx *gorm.DB, repoID uint) (*models.Repository, error) {
	var repo models.Repository

	if err := tx.Model(&models.Repository{}).Where(&models.Repository{ID: repoID}).First(&repo).Error; err != nil {
		return nil, err
	}

	return &repo, nil
}

// CreateRepository inserts a repository into the database
func (r *ReviewsRepository) CreateRepository(tx *gorm.DB, repo *models.Repository) (*models.Repository, error) {
	if err := tx.Create(repo).Error; err != nil {
		return nil, err
	}

	return repo, nil
}

// CreatePullRequests inserts pull requests into the database
func (r *ReviewsRepository) CreatePullRequests(tx *gorm.DB, prs []*models.PullRequest) error {
	if len(prs) == 0 {
		return nil
	}
	log.Printf("Inserting pull requests %v", prs)
	if err := tx.CreateInBatches(prs, 30).Error; err != nil {
		return fmt.Errorf("failed to insert pull requests: %v", err)
	}

	return nil
}

// GetPullRequests returns all pull requests for a repository
func (r *ReviewsRepository) GetPullRequests(tx *gorm.DB, userID, repoID uint) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest

	if err := tx.Model(&models.PullRequest{}).Where(&models.PullRequest{RepositoryID: repoID}).Find(&prs).Error; err != nil {
		return nil, err
	}

	return prs, nil
}

func (r *ReviewsRepository) GetPullRequest(tx *gorm.DB, prID uint) (*models.PullRequest, error) {
	var pr models.PullRequest

	if err := tx.Model(&models.PullRequest{}).Where(&models.PullRequest{ID: prID}).First(&pr).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *ReviewsRepository) UpdatePullRequestStatuses(tx *gorm.DB, prs []*models.PullRequest) error {
	if len(prs) == 0 {
		return nil
	}
	for _, pr := range prs {
		if err := tx.Model(&models.PullRequest{}).Where("id = ?", pr.ID).Updates(models.PullRequest{State: pr.State}).Error; err != nil {
			return fmt.Errorf("failed to update pull requests: %v", err)
		}
	}

	return nil
}

// GetReviews returns all reviews for a pull request
func (r *ReviewsRepository) GetReviews(tx *gorm.DB, repoID, prID uint) ([]*models.Review, error) {
	var reviews []*models.Review

	if err := tx.Model(&models.Review{}).Where(&models.Review{PullRequestID: prID}).Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

// GetReview returns a review for a pull request
func (r *ReviewsRepository) GetReview(tx *gorm.DB, repoID, prID, reviewID uint) (*models.Review, error) {
	var review models.Review

	if err := tx.Model(&models.Review{}).Where(&models.Review{PullRequestID: prID, ID: reviewID}).First(&review).Error; err != nil {
		return nil, err
	}

	return &review, nil
}

// GetFileReviews returns reviews for files
func (r *ReviewsRepository) GetFileReviews(tx *gorm.DB, reviewID uint) (*models.Review, error) {
	var review models.Review

	if err := tx.Model(&models.Review{}).Preload("FileReviews").Where(&models.Review{ID: reviewID}).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// CreateReview inserts a review into the database
func (r *ReviewsRepository) CreateReview(tx *gorm.DB, review *models.Review) (*models.Review, error) {
	if err := tx.Create(review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

// DeleteReview deletes a review from the database
func (r *ReviewsRepository) DeleteReview(tx *gorm.DB, reviewID uint) error {
	var review models.Review

	// preload the id of the review to delete the review status and file reviews
	if err := tx.Preload("FileReviews").Preload("ReviewStatus").First(&review, reviewID).Error; err != nil {
		return fmt.Errorf("failed to find review with id %d: %v", reviewID, err)
	}

	if err := tx.Delete(&review).Error; err != nil {
		return fmt.Errorf("failed to delete review with id %d: %v", reviewID, err)
	}

	return nil
}

// CreateReviewStatus inserts a review status into the database
func (r *ReviewsRepository) CreateReviewStatus(tx *gorm.DB, reviewStatus *models.ReviewStatus) error {
	if err := tx.Create(reviewStatus).Error; err != nil {
		return err
	}
	return nil
}

func (r *ReviewsRepository) GetReviewStatuses(tx *gorm.DB, reviewIDs []uint) (*[]models.ReviewStatus, error) {
	var reviewStatus []models.ReviewStatus

	if err := tx.Model(&models.ReviewStatus{}).Where("review_id IN ?", reviewIDs).Find(&reviewStatus).Error; err != nil {
		return nil, err
	}

	return &reviewStatus, nil
}

func (r *ReviewsRepository) GetReviewStatus(tx *gorm.DB, reviewID uint) (*models.ReviewStatus, error) {
	var reviewStatus models.ReviewStatus

	if err := tx.Model(&models.ReviewStatus{}).Where("review_id = ?", reviewID).First(&reviewStatus).Error; err != nil {
		return nil, err
	}

	return &reviewStatus, nil
}

func (r *ReviewsRepository) CreateFileReviews(tx *gorm.DB, fileReviews []*models.FileReview) error {
	if len(fileReviews) == 0 {
		log.Printf("No file reviews to insert")
		return nil
	}

	log.Printf("Inserting %d file reviews", len(fileReviews))
	if err := tx.CreateInBatches(fileReviews, 30).Error; err != nil {
		return fmt.Errorf("failed to insert file reviews: %v", err)
	}
	log.Printf("Successfully inserted %d file reviews", len(fileReviews))
	return nil
}

// UpdateReviewStatus updates the status of a review
func (r *ReviewsRepository) UpdateReviewStatus(tx *gorm.DB, reviewStatusID uint, status constants.ReviewStatus) error {
	if err := tx.Model(&models.ReviewStatus{}).Where("id = ?", reviewStatusID).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	log.Printf("Updated review status to %s", string(status))
	return nil
}

// UpdateReviewProgress updates the progress of a review
func (r *ReviewsRepository) UpdateReviewProgress(tx *gorm.DB, reviewStatusID uint, progress int) error {
	if err := tx.Model(&models.ReviewStatus{}).Where("id = ?", reviewStatusID).UpdateColumn("progress", progress).Error; err != nil {
		return err
	}
	log.Printf("Updated review progress to %d", progress)
	return nil
}
