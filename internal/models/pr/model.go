package pr

import (
	"AvitoProject/internal/models/user"
	"time"
)

type ID string

type PullRequest struct {
	ID               ID    `json:"pull_request_id"`
	Name             string    `json:"pull_request_name"`
	AuthorID         user.ID    `json:"author_id"`
	Status           string    `json:"status"`
	AssignedReviewers []user.ID `json:"assigned_reviewers"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	MergedAt         *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	ID       ID    `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID user.ID `json:"author_id"`
	Status   string `json:"status"`
}