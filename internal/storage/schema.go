package postgres

import (
	"context"
	"fmt"
)

func EnsureSchema(ctx context.Context, s *Storage) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS depts (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			phone VARCHAR(20) NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(50) NOT NULL,
			middle_name VARCHAR(50) NULL,
			last_name VARCHAR(50) NOT NULL,
			role VARCHAR(10) NOT NULL CHECK (role IN ('user', 'support')),
			dept_id INTEGER NULL REFERENCES depts(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
	}

	for _, q := range queries {
		if _, err := s.DB.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("ensure schema: %w", err)
		}
	}
	return nil
}
