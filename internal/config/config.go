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
	SeedPhone string
	SeedDept  string

	SeedAdmin2 bool
	SeedLogin2 string
	SeedEmail2 string
	SeedPass2  string
	SeedFirst2 string
	SeedLast2  string
	SeedPhone2 string
	SeedDept2  string

	SeedUser      bool
	SeedUserLogin string
	SeedUserEmail string
	SeedUserPass  string
	SeedUserFirst string
	SeedUserLast  string
	SeedUserPhone string
	SeedUserDept  string

	SeedTickets        bool
	SeedTicketsPerUser int

	SeedDepts bool
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
		SeedPhone: getenv("SEED_ADMIN_PHONE", "+7 (495) 123-45-67"),
		SeedDept:  getenv("SEED_ADMIN_DEPT", ""),

		SeedAdmin2: getenv("SEED_ADMIN2", "false") == "true",
		SeedLogin2: getenv("SEED_ADMIN2_LOGIN", "support"),
		SeedEmail2: getenv("SEED_ADMIN2_EMAIL", "support@local.test"),
		SeedPass2:  getenv("SEED_ADMIN2_PASSWORD", "support12345"),
		SeedFirst2: getenv("SEED_ADMIN2_FIRST", "Сергей"),
		SeedLast2:  getenv("SEED_ADMIN2_LAST", "Иванов"),
		SeedPhone2: getenv("SEED_ADMIN2_PHONE", "+7 (495) 999-88-88"),
		SeedDept2:  getenv("SEED_ADMIN2_DEPT", "IT Support"),

		SeedUser:      getenv("SEED_USER", "true") == "true",
		SeedUserLogin: getenv("SEED_USER_LOGIN", "user1"),
		SeedUserEmail: getenv("SEED_USER_EMAIL", "user1@local.test"),
		SeedUserPass:  getenv("SEED_USER_PASSWORD", "user12345"),
		SeedUserFirst: getenv("SEED_USER_FIRST", "Иван"),
		SeedUserLast:  getenv("SEED_USER_LAST", "Петров"),
		SeedUserPhone: getenv("SEED_USER_PHONE", "+7 (495) 999-88-77"),
		SeedUserDept:  getenv("SEED_USER_DEPT", ""),

		SeedTickets:        getenv("SEED_TICKETS", "true") == "true",
		SeedTicketsPerUser: mustInt(getenv("SEED_TICKETS_PER_USER", "6")),

		SeedDepts: getenv("SEED_DEPTS", "true") == "true",
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
