package api

import "time"

type PrReqDTO struct {
	PrID     string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type PrReassignReqDTO struct {
	PrID        string `json:"pull_request_id"`
	OldReviewer string `json:"old_reviewer_id"`
}

type PrDTO struct {
	PrID         string     `json:"pull_request_id"`
	Name         string     `json:"pull_request_name"`
	AuthorID     string     `json:"author_id"`
	Status       string     `json:"status"`
	ReviewersIDs []string   `json:"assigned_reviewers"`
	CreatedAt    time.Time  `json:"createdAt"`
	MergedAt     *time.Time `json:"mergedAt"`
}

type PrShortDTO struct {
	PrID     string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type PrAddResponse struct {
	Pr PrDTO `json:"pr"`
}

type PrReassignedResponse struct {
	Pr         PrDTO  `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}
