package models

import "time"

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
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
