package storage

import (
	"context"
	"database/sql"
)

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
}

type UsersRepo struct {
	DB *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{DB: db}
}

func (r *UsersRepo) GetByLogin(ctx context.Context, login string) (*User, error) {
	const q = `
		SELECT id, username, email, password_hash, first_name, last_name, role
		FROM users
		WHERE username=$1 OR email=$1
	`
	row := r.DB.QueryRowContext(ctx, q, login)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id int) (*User, error) {
	const q = `
		SELECT id, username, email, password_hash, first_name, last_name, role
		FROM users
		WHERE id=$1
	`
	row := r.DB.QueryRowContext(ctx, q, id)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepo) CreateSeedSupportIfNotExists(ctx context.Context, username, email, passHash, first, last string) error {
	const q = `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, 'support')
		ON CONFLICT (username) DO NOTHING
	`
	_, err := r.DB.ExecContext(ctx, q, username, email, passHash, first, last)
	return err
}

func (r *UsersRepo) CreateSeedUserIfNotExists(ctx context.Context, username, email, passHash, first, last string) error {
	const q = `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, 'user')
		ON CONFLICT (username) DO NOTHING
	`
	_, err := r.DB.ExecContext(ctx, q, username, email, passHash, first, last)
	return err
}
