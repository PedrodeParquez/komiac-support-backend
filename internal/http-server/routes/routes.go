package routes

import (
	"github.com/gin-gonic/gin"

	authapi "komiac-support-backend/internal/api/auth"
	"komiac-support-backend/internal/config"
	"komiac-support-backend/internal/http-server/middleware"
	postgres "komiac-support-backend/internal/storage"
)

func Register(r *gin.Engine, cfg config.Config, users *postgres.UsersRepo) {
	r.Use(middleware.CORS(middleware.CORSConfig{Origin: cfg.CorsOrigin}))

	authH := authapi.New(cfg, users)

	g := r.Group("/auth")
	{
		g.POST("/login", authH.Login)
		g.POST("/refresh", authH.Refresh)
		g.POST("/logout", authH.Logout)

		g.GET("/me",
			middleware.RequireAuth(middleware.AuthConfig{AccessSecret: cfg.AccessSecret}),
			authH.Me,
		)
	}
}
