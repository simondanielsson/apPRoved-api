package requests

type CreateRepositoryRequest struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Owner string `json:"owner"`
}

type CreateReviewRequest struct {
	Name string `json:"name"`
}

type FileReviewRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type CompleteReviewRequest struct {
	ReviewID       uint                `json:"review_id"`
	ReviewStatusID uint                `json:"review_status_id"`
	FileReviews    []FileReviewRequest `json:"file_reviews"`
}
