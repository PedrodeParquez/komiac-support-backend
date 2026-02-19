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

	usersRepo := postgres.NewUsersRepo(store.DB)

	if cfg.SeedAdmin {
		hash, err := auth.HashPassword(cfg.SeedPass)
		if err != nil {
			log.Fatal(err)
		}
		if err := usersRepo.CreateSeedSupportIfNotExists(ctx, cfg.SeedLogin, cfg.SeedEmail, hash, cfg.SeedFirst, cfg.SeedLast); err != nil {
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
		); err != nil {
			log.Fatal(err)
		}
	}

	r := gin.Default()
	routes.Register(r, cfg, usersRepo)

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
