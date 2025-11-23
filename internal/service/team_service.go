package service

import (
	"PR_project/internal/domain"
	"context"
	"errors"
)

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")
)

type TeamRepository interface {
	TeamExists(ctx context.Context, teamName string) (bool, error)
	GetTeam(ctx context.Context, teamName string) (domain.Team, []domain.User, error)
	CreateTeamWithMembers(
		ctx context.Context,
		teamName string,
		members []TeamMemberInput,
	) (domain.Team, []domain.User, error)
}

type TeamService struct {
	TRepository TeamRepository
}

type TeamMemberInput struct {
	UserID   string
	Username string
	IsActive bool
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (domain.Team, []domain.User, error) {
	ok, err := s.TRepository.TeamExists(ctx, teamName)
	if err != nil {
		return domain.Team{}, nil, err
	}
	if !ok {
		return domain.Team{}, nil, ErrTeamNotFound
	}

	return s.TRepository.GetTeam(ctx, teamName)
}

func (s *TeamService) CreateTeam(
	ctx context.Context,
	teamName string,
	teamMembers []TeamMemberInput) (domain.Team, []domain.User, error) {
	ok, err := s.TRepository.TeamExists(ctx, teamName)
	if err != nil {
		return domain.Team{}, nil, err
	}
	if ok {
		return domain.Team{}, nil, ErrTeamAlreadyExists
	}

	team, users, err := s.TRepository.CreateTeamWithMembers(ctx, teamName, teamMembers)
	if err != nil {
		return team, nil, err
	}

	return team, users, nil
}
