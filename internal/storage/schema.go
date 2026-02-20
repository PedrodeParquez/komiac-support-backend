package storage

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

		`
		CREATE TABLE IF NOT EXISTS tickets (
			id SERIAL PRIMARY KEY,
			ticket_number VARCHAR(50) NOT NULL UNIQUE,
			title VARCHAR(200) NOT NULL,
			description TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'open'
				CHECK (status IN ('open', 'in_progress', 'resolved', 'closed', 'reopened')),
			priority VARCHAR(10) NOT NULL DEFAULT 'medium'
				CHECK (priority IN ('low', 'medium', 'high')),
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			taken_by INTEGER NULL REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			resolved_at TIMESTAMP NULL,
			closed_at TIMESTAMP NULL
		);
		`,
		`CREATE INDEX IF NOT EXISTS idx_tickets_user_id ON tickets(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tickets_taken_by ON tickets(taken_by);`,
		`CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);`,
		`CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at);`,

		`
		CREATE TABLE IF NOT EXISTS ticket_messages (
			id SERIAL PRIMARY KEY,
			ticket_id INTEGER NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
			author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`CREATE INDEX IF NOT EXISTS idx_ticket_messages_ticket_id ON ticket_messages(ticket_id);`,
	}

	for _, q := range queries {
		if _, err := s.DB.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("ensure schema: %w", err)
		}
	}
	return nil
}
