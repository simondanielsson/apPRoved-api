package constants

type ReviewStatus string

const (
	StatusQueued     ReviewStatus = "queued"
	StatusProcessing ReviewStatus = "processing"
	StatusAvailable  ReviewStatus = "available"
)