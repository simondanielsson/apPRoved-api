package models

import (
	"time"

	"github.com/simondanielsson/apPRoved/cmd/constants"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PullRequest struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	RepositoryID uint       `json:"repository_id"`
	Repository   Repository `gorm:"foreignKey:RepositoryID" json:"repository"`
	Number       uint
	Title        string
	URL          string
	State        string
	LastCommit   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Review struct {
	ID            uint `gorm:"primary_key" json:"id"`
	Name          string
	PullRequestID uint         `json:"pull_request_id"`
	PullRequest   PullRequest  `gorm:"foreignKey:PullRequestID" json:"pull_request"`
	FileReviews   []FileReview `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE;" json:"file_reviews"`
	ReviewStatus  ReviewStatus `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE;" json:"review_status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type FileReview struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ReviewID  uint      `json:"review_id"`
	Filename  string    `json:"filename"`
	Content   string    `json:"content"`
	Patch     string    `json:"patch"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewStatus struct {
	ID        uint                   `gorm:"primary_key" json:"id"`
	ReviewID  uint                   `json:"review_id"`
	Status    constants.ReviewStatus `json:"status"`
	Progress  int                    `json:"progress"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
