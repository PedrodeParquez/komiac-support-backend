package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type SeedAccount struct {
	Enabled  bool
	Login    string
	Email    string
	Password string
	First    string
	Last     string
	Phone    string
	Dept     string
}

type Config struct {
	DatabaseURL string

	CorsOrigin   string
	CookieSecure bool
	CookieDomain string

	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration

	SeedAdmin SeedAccount
	SeedUser  SeedAccount
	SeedDepts bool
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		DatabaseURL: mustEnv("DATABASE_URL"),

		CorsOrigin:   env("CORS_ORIGIN", ""),
		CookieDomain: env("COOKIE_DOMAIN", ""),
		CookieSecure: envBool("COOKIE_SECURE", false),

		AccessSecret:  mustEnv("JWT_ACCESS_SECRET"),
		RefreshSecret: mustEnv("JWT_REFRESH_SECRET"),

		AccessTTL:  envDurationMinutes("JWT_ACCESS_TTL_MIN", 15),
		RefreshTTL: envDurationDays("JWT_REFRESH_TTL_DAYS", 30),

		SeedDepts: envBool("SEED_DEPTS", true),
		SeedAdmin: SeedAccount{
			Enabled:  envBool("SEED_ADMIN", true),
			Login:    env("SEED_ADMIN_LOGIN", ""),
			Email:    env("SEED_ADMIN_EMAIL", ""),
			Password: env("SEED_ADMIN_PASSWORD", ""),
			First:    env("SEED_ADMIN_FIRST", ""),
			Last:     env("SEED_ADMIN_LAST", ""),
			Phone:    env("SEED_ADMIN_PHONE", ""),
			Dept:     env("SEED_ADMIN_DEPT", ""),
		},
		SeedUser: SeedAccount{
			Enabled:  envBool("SEED_USER", false),
			Login:    env("SEED_USER_LOGIN", ""),
			Email:    env("SEED_USER_EMAIL", ""),
			Password: env("SEED_USER_PASSWORD", ""),
			First:    env("SEED_USER_FIRST", ""),
			Last:     env("SEED_USER_LAST", ""),
			Phone:    env("SEED_USER_PHONE", ""),
			Dept:     env("SEED_USER_DEPT", ""),
		},
	}
}

func env(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	return v
}

func envBool(k string, def bool) bool {
	v, ok := os.LookupEnv(k)
	if !ok || v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envDurationMinutes(k string, defMin int) time.Duration {
	v, ok := os.LookupEnv(k)
	if !ok || v == "" {
		return time.Duration(defMin) * time.Minute
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return time.Duration(defMin) * time.Minute
	}
	return time.Duration(n) * time.Minute
}

func envDurationDays(k string, defDays int) time.Duration {
	v, ok := os.LookupEnv(k)
	if !ok || v == "" {
		return time.Duration(defDays) * 24 * time.Hour
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return time.Duration(defDays) * 24 * time.Hour
	}
	return time.Duration(n) * 24 * time.Hour
}
