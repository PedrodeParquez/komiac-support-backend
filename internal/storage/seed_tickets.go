package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

type SeedTicketsConfig struct {
	Enabled    bool
	PerUser    int
	AdminLogin string // support
	UserLogin  string // обычный user
}

func EnsureSeedTickets(ctx context.Context, s *Storage, cfg SeedTicketsConfig) error {
	if !cfg.Enabled {
		return nil
	}

	if cfg.PerUser <= 0 {
		cfg.PerUser = 6
	}

	var existing int
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM tickets`).Scan(&existing); err != nil {
		return fmt.Errorf("seed tickets count: %w", err)
	}
	if existing > 0 {
		return nil
	}

	userID, err := getUserIDByUsername(ctx, s.DB, cfg.UserLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("seed tickets: user %q not found (seed users first)", cfg.UserLogin)
		}
		return fmt.Errorf("seed tickets: get user id: %w", err)
	}

	supportID, err := getUserIDByUsername(ctx, s.DB, cfg.AdminLogin)
	if err != nil {
		supportID = 0
	}

	now := time.Now()

	type tSeed struct {
		title      string
		desc       string
		status     string
		priority   string
		takenBy    *int
		resolvedAt *time.Time
		closedAt   *time.Time
	}

	tb := func(v int) *int { return &v }
	tt := func(t time.Time) *time.Time { return &t }

	var seeds []tSeed

	// 1) open
	seeds = append(seeds, tSeed{
		title:    "Не работает вход в систему",
		desc:     "Не удаётся войти: после ввода логина/пароля возвращает ошибку.",
		status:   "open",
		priority: "high",
	})

	// 2) in_progress (назначен support)
	if supportID != 0 {
		seeds = append(seeds, tSeed{
			title:    "Проблемы с доступом к разделу обращений",
			desc:     "Раздел обращений открывается, но список пустой.",
			status:   "in_progress",
			priority: "medium",
			takenBy:  tb(supportID),
		})
	} else {
		seeds = append(seeds, tSeed{
			title:    "Проблемы с доступом к разделу обращений",
			desc:     "Раздел обращений открывается, но список пустой.",
			status:   "in_progress",
			priority: "medium",
		})
	}

	// 3) resolved
	resolvedAt := now.Add(-6 * time.Hour)
	if supportID != 0 {
		seeds = append(seeds, tSeed{
			title:      "Не открываются вложения",
			desc:       "При попытке открыть вложение ничего не происходит.",
			status:     "resolved",
			priority:   "low",
			takenBy:    tb(supportID),
			resolvedAt: tt(resolvedAt),
		})
	} else {
		seeds = append(seeds, tSeed{
			title:      "Не открываются вложения",
			desc:       "При попытке открыть вложение ничего не происходит.",
			status:     "resolved",
			priority:   "low",
			resolvedAt: tt(resolvedAt),
		})
	}

	// 4) closed
	closedAt := now.Add(-24 * time.Hour)
	if supportID != 0 {
		seeds = append(seeds, tSeed{
			title:    "Сброс пароля",
			desc:     "Нужен сброс пароля для учётной записи.",
			status:   "closed",
			priority: "medium",
			takenBy:  tb(supportID),
			closedAt: tt(closedAt),
		})
	} else {
		seeds = append(seeds, tSeed{
			title:    "Сброс пароля",
			desc:     "Нужен сброс пароля для учётной записи.",
			status:   "closed",
			priority: "medium",
			closedAt: tt(closedAt),
		})
	}

	// добьём до cfg.PerUser
	for len(seeds) < cfg.PerUser {
		seeds = append(seeds, tSeed{
			title:    "Тестовое обращение",
			desc:     "Автосгенерированное обращение для демонстрации интерфейса.",
			status:   "open",
			priority: "low",
		})
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("seed tickets begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for i, t := range seeds {
		num := fmt.Sprintf("T-%s-%02d", randHex(4), i+1)

		_, err := tx.ExecContext(ctx, `
			INSERT INTO tickets (
				ticket_number, title, description, status, priority,
				user_id, taken_by, resolved_at, closed_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			ON CONFLICT (ticket_number) DO NOTHING
		`,
			num, t.title, t.desc, t.status, t.priority,
			userID, t.takenBy, t.resolvedAt, t.closedAt,
		)
		if err != nil {
			return fmt.Errorf("seed tickets insert: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed tickets commit: %w", err)
	}

	return nil
}

func getUserIDByUsername(ctx context.Context, db *sql.DB, username string) (int, error) {
	var id int
	hookup := `SELECT id FROM users WHERE username = $1`
	err := db.QueryRowContext(ctx, hookup, username).Scan(&id)
	return id, err
}

func randHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
