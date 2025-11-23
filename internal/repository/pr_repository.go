package repository

import (
	"PR_project/internal/domain"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostgresPrRepository struct {
	db *sql.DB
}

func NewPostgresPrRepository(db *sql.DB) *PostgresPrRepository {
	return &PostgresPrRepository{db: db}
}

func (r *PostgresPrRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	query := `SELECT 1 FROM pull_requests WHERE pull_request_id = $1`
	row := r.db.QueryRowContext(ctx, query, prID)
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

func (r *PostgresPrRepository) CreatePRWithReviewers(
	ctx context.Context,
	pr domain.PullRequest,
	reviewerIDs []string,
) (domain.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.PullRequest{}, err
	}
	defer tx.Rollback()

	insertPrQuery := `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, insertPrQuery, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return domain.PullRequest{}, err
	}

	insertReviewersQuery := `INSERT INTO pr_reviewers (pull_request_id, user_id)
VALUES ($1, $2)`

	for _, rev := range reviewerIDs {
		_, err = tx.ExecContext(
			ctx,
			insertReviewersQuery,
			pr.ID,
			rev,
		)
		if err != nil {
			return domain.PullRequest{}, err
		}
	}

	prUpdated := domain.PullRequest{
		ID:           pr.ID,
		Name:         pr.Name,
		AuthorID:     pr.AuthorID,
		Status:       pr.Status,
		CreatedAt:    pr.CreatedAt,
		MergedAt:     pr.MergedAt,
		ReviewersIDs: reviewerIDs,
	}

	if err := tx.Commit(); err != nil {
		return domain.PullRequest{}, err
	}

	return prUpdated, nil
}

func (r *PostgresPrRepository) GetPRWithReviewers(ctx context.Context, prID string) (domain.PullRequest, error) {
	selectPrsQuery := `SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at 
FROM pull_requests
WHERE pull_request_id = $1`
	row := r.db.QueryRowContext(ctx, selectPrsQuery, prID)

	var pr domain.PullRequest

	err := row.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PullRequest{}, err
	}
	if err != nil {
		return domain.PullRequest{}, err
	}

	selectReviewersQuery := `SELECT user_id 
FROM pr_reviewers 
WHERE pull_request_id = $1`
	revRows, err := r.db.QueryContext(ctx, selectReviewersQuery, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	defer revRows.Close()

	var reviewerIDs []string
	for revRows.Next() {
		var userID string
		err = revRows.Scan(&userID)
		if err != nil {
			return domain.PullRequest{}, err
		}

		reviewerIDs = append(reviewerIDs, userID)
	}

	if err := revRows.Err(); err != nil {
		return domain.PullRequest{}, err
	}

	pullRequest := domain.PullRequest{
		ID:           pr.ID,
		Name:         pr.Name,
		AuthorID:     pr.AuthorID,
		Status:       pr.Status,
		CreatedAt:    pr.CreatedAt,
		MergedAt:     pr.MergedAt,
		ReviewersIDs: reviewerIDs,
	}

	return pullRequest, nil
}

func (r *PostgresPrRepository) MarkMerged(ctx context.Context, prID string, mergedAt time.Time) error {
	updatePrsQuery := `UPDATE pull_requests 
SET status = 'MERGED',
    merged_at = $2
WHERE pull_request_id = $1`
	_, err := r.db.ExecContext(ctx, updatePrsQuery, prID, mergedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresPrRepository) ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteOldQuery := `DELETE FROM pr_reviewers
WHERE pull_request_id = $1 AND user_id = $2;`
	_, err = tx.ExecContext(ctx, deleteOldQuery, prID, oldUserID)
	if err != nil {
		return err
	}

	insertNewQuery := `INSERT INTO pr_reviewers (pull_request_id, user_id)
VALUES ($1, $2);`
	_, err = tx.ExecContext(ctx, insertNewQuery, prID, newUserID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
