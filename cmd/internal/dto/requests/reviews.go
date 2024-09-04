package requests

type CreateRepositoryRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
