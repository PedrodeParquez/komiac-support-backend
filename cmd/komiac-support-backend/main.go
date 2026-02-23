package main

import (
	"context"
	"log"

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
	if cfg.AccessSecret == "" || cfg.RefreshSecret == "" {
		log.Fatal("JWT secrets are empty")
	}

	ctx := context.Background()

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

	if cfg.SeedAdmin {
		hash, err := auth.HashPassword(cfg.SeedPass)
		if err != nil {
			log.Fatal(err)
		}
		if err := usersRepo.CreateSeedSupportIfNotExists(ctx, cfg.SeedLogin, cfg.SeedEmail, hash, cfg.SeedFirst, cfg.SeedLast, cfg.SeedPhone, cfg.SeedDept); err != nil {
			log.Fatal(err)
		}
	}

	if cfg.SeedAdmin2 {
		hash, err := auth.HashPassword(cfg.SeedPass2)
		if err != nil {
			log.Fatal(err)
		}
		if err := usersRepo.CreateSeedSupportIfNotExists(ctx, cfg.SeedLogin2, cfg.SeedEmail2, hash, cfg.SeedFirst2, cfg.SeedLast2, cfg.SeedPhone2, cfg.SeedDept2); err != nil {
			log.Fatal(err)
		}
	}

	if cfg.SeedUser {
		hash, err := auth.HashPassword(cfg.SeedUserPass)
		if err != nil {
			log.Fatal(err)
		}
		if err := usersRepo.CreateSeedUserIfNotExists(
			ctx,
			cfg.SeedUserLogin,
			cfg.SeedUserEmail,
			hash,
			cfg.SeedUserFirst,
			cfg.SeedUserLast,
			cfg.SeedUserPhone,
			cfg.SeedUserDept,
		); err != nil {
			log.Fatal(err)
		}
	}

	r := gin.Default()
	routes.Register(r, cfg, usersRepo, ticketsRepo)

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
