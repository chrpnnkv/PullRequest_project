package repository

import (
	"PR_project/internal/domain"
	"PR_project/internal/service"
	"context"
	"database/sql"
	"errors"
)

type PostgresTeamRepository struct {
	db *sql.DB
}

func NewPostgresTeamRepository(db *sql.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	query := `SELECT 1 FROM teams WHERE team_name = $1`
	row := r.db.QueryRowContext(ctx, query, teamName)
	var count int
	err := row.Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (r *PostgresTeamRepository) GetTeam(ctx context.Context, teamName string) (domain.Team, []domain.User, error) {
	team := domain.Team{Name: teamName}

	selectUsersQuery := `SELECT user_id, username, team_name, is_active
FROM users
WHERE team_name = $1;`
	userRows, err := r.db.QueryContext(ctx, selectUsersQuery, teamName)
	if err != nil {
		return domain.Team{}, nil, err
	}
	defer userRows.Close()

	var users []domain.User
	for userRows.Next() {
		var user_id, username, team_name string
		var is_active bool
		err = userRows.Scan(&user_id, &username, &team_name, &is_active)
		if err != nil {
			return domain.Team{}, nil, err
		}

		user := domain.User{
			ID:       user_id,
			Username: username,
			TeamName: team_name,
			IsActive: is_active,
		}
		users = append(users, user)
	}
	if err := userRows.Err(); err != nil {
		return domain.Team{}, nil, err
	}

	return team, users, nil
}

func (r *PostgresTeamRepository) CreateTeamWithMembers(
	ctx context.Context,
	teamName string,
	members []service.TeamMemberInput,
) (domain.Team, []domain.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Team{}, nil, err
	}
	defer tx.Rollback()

	insertTeamQuery := `INSERT INTO teams (team_name) VALUES ($1)`
	_, err = tx.ExecContext(ctx, insertTeamQuery, teamName)
	if err != nil {
		return domain.Team{}, nil, err
	}

	insertUsersQuery := `INSERT INTO users (user_id, username, team_name, is_active) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id)
DO UPDATE SET
username = EXCLUDED.username,
    team_name = EXCLUDED.team_name,
    is_active = EXCLUDED.is_active;`

	for _, member := range members {
		_, err = tx.ExecContext(
			ctx,
			insertUsersQuery,
			member.UserID,
			member.Username,
			teamName,
			member.IsActive,
		)
		if err != nil {
			return domain.Team{}, nil, err
		}
	}

	team := domain.Team{
		Name: teamName,
	}

	selectUsersQuery := `SELECT user_id, username, team_name, is_active
FROM users
WHERE team_name = $1;`
	rows, err := tx.QueryContext(ctx, selectUsersQuery, teamName)
	if err != nil {
		return domain.Team{}, nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user_id, username, team_name string
		var is_active bool
		err = rows.Scan(&user_id, &username, &team_name, &is_active)
		if err != nil {
			return domain.Team{}, nil, err
		}

		user := domain.User{
			ID:       user_id,
			Username: username,
			TeamName: team_name,
			IsActive: is_active,
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return domain.Team{}, nil, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Team{}, nil, err
	}

	return team, users, nil
}
