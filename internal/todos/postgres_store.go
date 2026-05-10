package todos

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func OpenPostgres(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS todos (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (s *PostgresStore) List(ctx context.Context) ([]Todo, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, completed
		FROM todos
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (s *PostgresStore) Get(ctx context.Context, id int64) (Todo, error) {
	var todo Todo
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, completed
		FROM todos
		WHERE id = $1
	`, id).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *PostgresStore) Create(ctx context.Context, title string) (Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Todo{}, errors.New("title is required")
	}

	var todo Todo
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO todos (title)
		VALUES ($1)
		RETURNING id, title, completed
	`, title).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *PostgresStore) Update(ctx context.Context, id int64, title string, completed bool) (Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Todo{}, errors.New("title is required")
	}

	var todo Todo
	err := s.db.QueryRowContext(ctx, `
		UPDATE todos
		SET title = $1, completed = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, completed
	`, title, completed, id).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *PostgresStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
