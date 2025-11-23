package service

import (
	"PR_project/internal/domain"
	"context"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	UserExists(ctx context.Context, userID string) (bool, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
	GetReviews(ctx context.Context, userID string) ([]domain.PullRequest, error)
	GetUserByID(ctx context.Context, userID string) (domain.User, error)
	GetActiveTeamMembersExcept(ctx context.Context, teamName string, excludeID string) ([]domain.User, error)
}

type UserService struct {
	URepository UserRepository
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	ok, err := s.URepository.UserExists(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	if !ok {
		return domain.User{}, ErrUserNotFound
	}

	return s.URepository.SetIsActive(ctx, userID, isActive)
}

func (s *UserService) GetReviews(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	ok, err := s.URepository.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUserNotFound
	}

	return s.URepository.GetReviews(ctx, userID)
}
