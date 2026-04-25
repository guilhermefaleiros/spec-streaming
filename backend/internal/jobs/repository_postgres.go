package jobs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, j *Job) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO transcoding_jobs (id, video_id, status, attempts, error_message)
		VALUES ($1, $2, $3, $4, $5)
	`, j.ID, j.VideoID, string(j.Status), j.Attempts, j.ErrorMessage)
	return err
}

func (r *PostgresRepository) ClaimPending(ctx context.Context) (*Job, error) {
	// Use advisory lock or SKIP LOCKED for concurrency safety
	var job Job
	err := r.pool.QueryRow(ctx, `
		UPDATE transcoding_jobs
		SET status = 'processing', attempts = attempts + 1, started_at = NOW()
		WHERE id = (
			SELECT id FROM transcoding_jobs
			WHERE status = 'pending'
			ORDER BY created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, video_id, status, attempts, error_message
	`).Scan(&job.ID, &job.VideoID, &job.Status, &job.Attempts, &job.ErrorMessage)
	if err != nil {
		// No pending jobs
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("claim pending: %w", err)
	}
	return &job, nil
}

func (r *PostgresRepository) Update(ctx context.Context, j *Job) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE transcoding_jobs
		SET status = $2, attempts = $3, error_message = $4, finished_at = CASE WHEN $2 IN ('completed', 'failed') THEN NOW() ELSE finished_at END, updated_at = NOW()
		WHERE id = $1
	`, j.ID, string(j.Status), j.Attempts, j.ErrorMessage)
	return err
}
