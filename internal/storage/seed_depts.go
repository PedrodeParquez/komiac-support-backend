package storage

import (
	"context"
	"fmt"
)

type SeedDeptsConfig struct {
	Enabled bool
}

func EnsureSeedDepts(ctx context.Context, s *Storage, cfg SeedDeptsConfig) error {
	if !cfg.Enabled {
		return nil
	}

	var existing int
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM depts`).Scan(&existing); err != nil {
		return fmt.Errorf("seed depts count: %w", err)
	}
	if existing > 0 {
		return nil
	}

	depts := []struct {
		name  string
		phone string
	}{
		{"IT Support", "+7 (495) 123-45-67"},
		{"Продажи", "+7 (495) 123-45-68"},
		{"Бухгалтерия", "+7 (495) 123-45-69"},
		{"HR", "+7 (495) 123-45-70"},
	}

	for _, d := range depts {
		_, err := s.DB.ExecContext(ctx, `
INSERT INTO depts(name, phone)
VALUES ($1, $2)
ON CONFLICT (name) DO NOTHING;
`, d.name, d.phone)
		if err != nil {
			return fmt.Errorf("seed depts insert: %w", err)
		}
	}

	return nil
}
