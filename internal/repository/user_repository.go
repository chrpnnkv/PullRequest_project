package repository

import (
	"PR_project/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT 1 FROM users WHERE user_id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
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

func (r *PostgresUserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	updateUsersQuery := `UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING user_id, username, team_name, is_active;`
	row := r.db.QueryRowContext(ctx, updateUsersQuery, userID, isActive)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

func (r *PostgresUserRepository) GetReviews(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	selectPrsQuery := `SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
FROM pull_requests pr
JOIN pr_reviewers rev
    ON rev.pull_request_id = pr.pull_request_id
WHERE rev.user_id = $1;`
	rows, err := r.db.QueryContext(ctx, selectPrsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []domain.PullRequest
	for rows.Next() {
		var prID, prName, authorID, status string
		err = rows.Scan(&prID, &prName, &authorID, &status)
		if err != nil {
			return nil, err
		}

		pr := domain.PullRequest{
			ID:       prID,
			Name:     prName,
			AuthorID: authorID,
			Status:   status,
		}
		pullRequests = append(pullRequests, pr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active
FROM users
WHERE user_id = $1;
`
	row := r.db.QueryRowContext(ctx, query, userID)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepository) GetActiveTeamMembersExcept(
	ctx context.Context,
	teamName string,
	excludeID string,
) ([]domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active
FROM users
WHERE team_name = $1
  AND is_active = TRUE
  AND user_id <> $2`

	rows, err := r.db.QueryContext(ctx, query, teamName, excludeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var userID, username, teamname string
		var isActive bool
		err = rows.Scan(&userID, &username, &teamname, &isActive)
		if err != nil {
			return nil, err
		}

		u := domain.User{
			ID:       userID,
			Username: username,
			TeamName: teamname,
			IsActive: isActive,
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
