package service

import (
	"PR_project/internal/domain"
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"slices"
	"time"
)

var (
	ErrPrNotFound        = errors.New("pull_request not found")
	ErrPrAlreadyExists   = errors.New("pull_request already exists")
	ErrPrAlreadyMerged   = errors.New("pull_request is already merged")
	ErrUserIsNotReviewer = errors.New("user is not reviewer")
	ErrNoCandidate       = errors.New("no candidate")
)

type PrRepository interface {
	PRExists(ctx context.Context, prID string) (bool, error)
	CreatePRWithReviewers(
		ctx context.Context,
		pr domain.PullRequest,
		reviewerIDs []string,
	) (domain.PullRequest, error)
	GetPRWithReviewers(ctx context.Context, prID string) (domain.PullRequest, error)
	MarkMerged(ctx context.Context, prID string, mergedAt time.Time) error
	ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error
}

type PrService struct {
	PrRepository PrRepository
	URepository  UserRepository
}

func (s *PrService) CreatePRWithReviewers(ctx context.Context, prID, prName, authorID string) (domain.PullRequest, error) {
	ok, err := s.PrRepository.PRExists(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	if ok {
		return domain.PullRequest{}, ErrPrAlreadyExists
	}

	author, err := s.URepository.GetUserByID(ctx, authorID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	activeTeamMembers, err := s.URepository.GetActiveTeamMembersExcept(ctx, author.TeamName, authorID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	pr := domain.PullRequest{
		ID:        prID,
		Name:      prName,
		AuthorID:  authorID,
		Status:    "OPEN",
		CreatedAt: time.Now(),
		MergedAt:  nil,
	}

	var reviewerIDs []string
	activeMembersLen := len(activeTeamMembers)
	if activeMembersLen > 0 {
		if activeMembersLen > 1 {
			first := rand.Intn(activeMembersLen)
			second := rand.Intn(activeMembersLen)
			for {
				if second != first {
					break
				}
				second = rand.Intn(activeMembersLen)
			}
			reviewerIDs = append(reviewerIDs, activeTeamMembers[first].ID, activeTeamMembers[second].ID)
		} else {
			reviewerIDs = append(reviewerIDs, activeTeamMembers[0].ID)
		}
	}

	pullRequest, err := s.PrRepository.CreatePRWithReviewers(ctx, pr, reviewerIDs)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return pullRequest, nil
}

func (s *PrService) Merge(ctx context.Context, prID string, mergedAt time.Time) (domain.PullRequest, error) {
	pr, err := s.PrRepository.GetPRWithReviewers(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, ErrPrNotFound
		}
		return domain.PullRequest{}, err
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	err = s.PrRepository.MarkMerged(ctx, prID, mergedAt)
	if err != nil {
		return domain.PullRequest{}, err
	}

	updatedPr, err := s.PrRepository.GetPRWithReviewers(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return updatedPr, nil
}

func (s *PrService) Reassign(ctx context.Context, prID, oldRevId string) (domain.PullRequest, string, error) {
	pr, err := s.PrRepository.GetPRWithReviewers(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, "", ErrPrNotFound
		}
		return domain.PullRequest{}, "", err
	}

	if pr.Status == "MERGED" {
		return domain.PullRequest{}, "", ErrPrAlreadyMerged
	}

	if !slices.Contains(pr.ReviewersIDs, oldRevId) {
		return domain.PullRequest{}, "", ErrUserIsNotReviewer
	}

	oldRev, err := s.URepository.GetUserByID(ctx, oldRevId)
	if err != nil {
		return domain.PullRequest{}, "", err
	}
	activeTeamMembers, err := s.URepository.GetActiveTeamMembersExcept(ctx, oldRev.TeamName, oldRevId)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	var candidates []string
	for _, m := range activeTeamMembers {
		if (m.ID != oldRevId) && (m.ID != pr.AuthorID) &&
			(!slices.Contains(pr.ReviewersIDs, m.ID)) {
			candidates = append(candidates, m.ID)
		}
	}
	if len(candidates) == 0 {
		return domain.PullRequest{}, "", ErrNoCandidate
	}

	newRevID := candidates[rand.Intn(len(candidates))]

	err = s.PrRepository.ReplaceReviewer(ctx, prID, oldRevId, newRevID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	pullRequest, err := s.PrRepository.GetPRWithReviewers(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	return pullRequest, newRevID, nil
}
