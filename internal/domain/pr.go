package domain

import "time"

type PullRequest struct {
	ID           string
	Name         string
	AuthorID     string
	Status       string
	ReviewersIDs []string
	CreatedAt    time.Time
	MergedAt     *time.Time
}
