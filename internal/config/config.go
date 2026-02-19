package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string

	CorsOrigin   string
	CookieSecure bool
	CookieDomain string

	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration

	SeedAdmin bool
	SeedLogin string
	SeedEmail string
	SeedPass  string
	SeedFirst string
	SeedLast  string

	SeedUser      bool
	SeedUserLogin string
	SeedUserEmail string
	SeedUserPass  string
	SeedUserFirst string
	SeedUserLast  string
}

func Load() Config {
	_ = godotenv.Load()

	accessMin := mustInt(getenv("JWT_ACCESS_TTL_MIN", "15"))
	refreshDays := mustInt(getenv("JWT_REFRESH_TTL_DAYS", "14"))

	return Config{
		DatabaseURL: getenv("DATABASE_URL", ""),

		CorsOrigin:   getenv("CORS_ORIGIN", "http://localhost:5173"),
		CookieSecure: getenv("COOKIE_SECURE", "false") == "true",
		CookieDomain: getenv("COOKIE_DOMAIN", ""),

		AccessSecret:  getenv("JWT_ACCESS_SECRET", ""),
		RefreshSecret: getenv("JWT_REFRESH_SECRET", ""),
		AccessTTL:     time.Duration(accessMin) * time.Minute,
		RefreshTTL:    time.Duration(refreshDays) * 24 * time.Hour,

		SeedAdmin: getenv("SEED_ADMIN", "true") == "true",
		SeedLogin: getenv("SEED_ADMIN_LOGIN", "admin"),
		SeedEmail: getenv("SEED_ADMIN_EMAIL", "admin@local.test"),
		SeedPass:  getenv("SEED_ADMIN_PASSWORD", "admin12345"),
		SeedFirst: getenv("SEED_ADMIN_FIRST", "Марина"),
		SeedLast:  getenv("SEED_ADMIN_LAST", "Шпегель"),

		SeedUser:      getenv("SEED_USER", "true") == "true",
		SeedUserLogin: getenv("SEED_USER_LOGIN", "user1"),
		SeedUserEmail: getenv("SEED_USER_EMAIL", "user1@local.test"),
		SeedUserPass:  getenv("SEED_USER_PASSWORD", "user12345"),
		SeedUserFirst: getenv("SEED_USER_FIRST", "Иван"),
		SeedUserLast:  getenv("SEED_USER_LAST", "Петров"),
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func mustInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
