package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(ctx context.Context, connectionString string) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := newPgxPool(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db, err := wrapPgxPool(ctx, pool)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	const op = "storage.postgres.Stop"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

//go:embed sql/save_user.sql
var saveUserSql string

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.postgres.SaveUser"

	rows, err := s.db.QueryxContext(ctx, saveUserSql, email, passHash)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func(rows *sqlx.Rows) {
		if tempErr := rows.Close(); tempErr != nil {
			err = tempErr
		}
	}(rows)

	if rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return uid, nil
}

//go:embed sql/get_user.sql
var getUserSql string

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User
	err := s.db.GetContext(ctx, &user, getUserSql, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

//go:embed sql/get_app.sql
var getAppSql string

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"

	var app models.App
	err := s.db.GetContext(ctx, &app, getAppSql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func newPgxPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	if poolConfig == nil {
		return nil, errors.New("parsed config is nil")
	}

	return pgxpool.NewWithConfig(ctx, poolConfig)
}

func wrapPgxPool(ctx context.Context, pool *pgxpool.Pool) (*sqlx.DB, error) {
	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), "pgx")
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
