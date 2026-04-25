package videos

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, v *Video) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO videos (id, title, original_filename, status, source_storage_key, manifest_storage_key, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, v.ID, v.Title, v.OriginalFilename, string(v.Status), v.SourceStorageKey, v.ManifestKey, v.ErrorMessage)
	return err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Video, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, original_filename, status, source_storage_key, manifest_storage_key, error_message
		FROM videos
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		err := rows.Scan(&v.ID, &v.Title, &v.OriginalFilename, &v.Status, &v.SourceStorageKey, &v.ManifestKey, &v.ErrorMessage)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Video, error) {
	var v Video
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, original_filename, status, source_storage_key, manifest_storage_key, error_message
		FROM videos
		WHERE id = $1
	`, id).Scan(&v.ID, &v.Title, &v.OriginalFilename, &v.Status, &v.SourceStorageKey, &v.ManifestKey, &v.ErrorMessage)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PostgresRepository) Update(ctx context.Context, v *Video) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE videos
		SET title = $2, original_filename = $3, status = $4, source_storage_key = $5, manifest_storage_key = $6, error_message = $7, updated_at = NOW()
		WHERE id = $1
	`, v.ID, v.Title, v.OriginalFilename, string(v.Status), v.SourceStorageKey, v.ManifestKey, v.ErrorMessage)
	return err
}
