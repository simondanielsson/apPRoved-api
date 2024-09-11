package responses

import "time"

type GetRepositoriesResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetPullRequestResponse struct {
	ID        uint      `json:"id"`
	Number    uint      `json:"number"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetReviewsResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetReviewResponse struct {
	ID          uint                     `json:"id"`
	FileReviews []*GetFileReviewResponse `json:"file_reviews"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type GetFileReviewResponse struct {
	ID        uint      `json:"id"`
	Filename  string    `json:"filename"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}