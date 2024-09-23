package constants

type ReviewStatus string

// Review Status Constants
const (
	// StatusQueued indicates that the review is queued.
	StatusQueued ReviewStatus = "queued"
	// StatusProcessing indicates that the review is in progress.
	StatusProcessing ReviewStatus = "processing"
	// StatusAvailable indicates that the review is available.
	StatusAvailable ReviewStatus = "available"
)

type PRState string

// Pull Request State Constants
const (
	// PRStateClosed indicates that the pull request is closed.
	PRStateClosed PRState = "closed"
	// PRStateOpen indicates that the pull request is open.
	PRStateOpen PRState = "open"
)
