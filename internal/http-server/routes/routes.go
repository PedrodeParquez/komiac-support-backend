package routes

import (
	"github.com/gin-gonic/gin"

	authapi "komiac-support-backend/internal/api/auth"
	ticketsapi "komiac-support-backend/internal/api/tickets"
	"komiac-support-backend/internal/config"
	"komiac-support-backend/internal/http-server/middleware"
	postgres "komiac-support-backend/internal/storage"
)

func Register(r *gin.Engine, cfg config.Config, users *postgres.UsersRepo, tickets *postgres.TicketsRepo) {
	r.Use(middleware.CORS(middleware.CORSConfig{Origin: cfg.CorsOrigin}))

	authH := authapi.New(cfg, users)
	ticketsH := ticketsapi.New(tickets)

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

	t := r.Group("/tickets",
		middleware.RequireAuth(middleware.AuthConfig{AccessSecret: cfg.AccessSecret}),
	)
	{
		t.GET("", ticketsH.ListTickets)
		t.GET("/my", ticketsH.ListMyTickets)
		t.GET("/:id", ticketsH.GetTicket)
		t.POST("", ticketsH.CreateTicket)
		t.POST("/:id/assign", ticketsH.AssignTicket)
		t.POST("/:id/messages", ticketsH.AddMessage)
	}
}
