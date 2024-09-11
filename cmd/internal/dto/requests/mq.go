package requests

import "github.com/simondanielsson/apPRoved/pkg/utils"

type FileDiffReviewRequest struct {
	ReviewID       uint                                  `json:"review_id" validate:"required"`
	ReviewStatusID uint                                  `json:"review_status_id" validate:"required"`
	FileDiffs      []*utils.GithubPullRequestFileChanges `json:"file_diffs" validate:"required"`
}
