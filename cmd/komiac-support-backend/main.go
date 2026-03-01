package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"komiac-support-backend/internal/auth"
	"komiac-support-backend/internal/config"
	"komiac-support-backend/internal/http-server/routes"
	postgres "komiac-support-backend/internal/storage"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store, err := postgres.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = store.Close() }()

	if err := postgres.EnsureSchema(ctx, store); err != nil {
		log.Fatal(err)
	}

	if err := postgres.EnsureSeedDepts(ctx, store, postgres.SeedDeptsConfig{
		Enabled: cfg.SeedDepts,
	}); err != nil {
		log.Fatal(err)
	}

	usersRepo := postgres.NewUsersRepo(store.DB)
	ticketsRepo := postgres.NewTicketsRepo(store.DB)

	r := gin.Default()
	routes.Register(r, cfg, usersRepo, ticketsRepo)

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func seedUsers(ctx context.Context, cfg config.Config, usersRepo *postgres.UsersRepo) error {
	if cfg.SeedAdmin.Enabled {
		hash, err := auth.HashPassword(cfg.SeedAdmin.Password)
		if err != nil {
			return err
		}

		if err := usersRepo.CreateSeedSupportIfNotExists(
			ctx,
			cfg.SeedAdmin.Login,
			cfg.SeedAdmin.Email,
			hash,
			cfg.SeedAdmin.First,
			cfg.SeedAdmin.Last,
			cfg.SeedAdmin.Phone,
			cfg.SeedAdmin.Dept,
		); err != nil {
			return err
		}
	}

	if cfg.SeedUser.Enabled {
		hash, err := auth.HashPassword(cfg.SeedUser.Password)
		if err != nil {
			return err
		}

		if err := usersRepo.CreateSeedUserIfNotExists(
			ctx,
			cfg.SeedUser.Login,
			cfg.SeedUser.Email,
			hash,
			cfg.SeedUser.First,
			cfg.SeedUser.Last,
			cfg.SeedUser.Phone,
			cfg.SeedUser.Dept,
		); err != nil {
			return err
		}
	}

	return nil
}
