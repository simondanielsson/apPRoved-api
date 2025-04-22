package requests

import "github.com/simondanielsson/apPRoved/pkg/utils/git"

type FileDiffReviewRequest struct {
	ReviewID       uint                             `json:"review_id" validate:"required"`
	ReviewStatusID uint                             `json:"review_status_id" validate:"required"`
	FileDiffs      []*git.GitPullRequestFileChanges `json:"file_diffs" validate:"required"`
}
