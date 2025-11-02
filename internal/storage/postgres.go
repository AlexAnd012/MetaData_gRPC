package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	mediametav1 "github.com/AlexAnd012/mediameta/gen/go/mediameta/v1"
)

// Реализация интерфейса Repository для Postgres
type PostgresRepo struct{ db *sql.DB }

func NewPostgresRepo(db *sql.DB) *PostgresRepo { return &PostgresRepo{db: db} }

func (r *PostgresRepo) Insert(ctx context.Context, m *mediametav1.Metadata) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO metadata (id, filename, size_bytes, content_type, owner_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		m.Id, m.Filename, m.SizeBytes, m.ContentType, m.OwnerId, m.CreatedAt)
	return err
}

func (r *PostgresRepo) Get(ctx context.Context, id string) (*mediametav1.Metadata, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, filename, size_bytes, content_type, owner_id, created_at
		FROM metadata WHERE id=$1`, id)
	m := &mediametav1.Metadata{}
	if err := row.Scan(&m.Id, &m.Filename, &m.SizeBytes, &m.ContentType, &m.OwnerId, &m.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}
	return m, nil
}

func (r *PostgresRepo) List(ctx context.Context, limit, offset int) ([]*mediametav1.Metadata, int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, filename, size_bytes, content_type, owner_id, created_at
		FROM metadata
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*mediametav1.Metadata

	for rows.Next() {
		m := &mediametav1.Metadata{}
		if err := rows.Scan(&m.Id, &m.Filename, &m.SizeBytes, &m.ContentType, &m.OwnerId, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, m)
	}
	//Если вернули ровно limit строк, возможно есть следующая страница
	//next = offset + limit, иначе 0.
	next := 0
	if len(out) == limit {
		next = offset + limit
	}
	return out, next, nil
}

func (r *PostgresRepo) Update(ctx context.Context, m *mediametav1.Metadata) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE metadata SET filename=$2, content_type=$3
		WHERE id=$1`,
		m.Id, m.Filename, m.ContentType)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM metadata WHERE id=$1`, id)
	return err
}
