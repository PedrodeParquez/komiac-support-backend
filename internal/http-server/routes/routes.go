package routes

import (
	"github.com/gin-gonic/gin"

	authapi "komiac-support-backend/internal/api/auth"
	ticketsapi "komiac-support-backend/internal/api/tickets"
	usersapi "komiac-support-backend/internal/api/users"

	"komiac-support-backend/internal/config"
	"komiac-support-backend/internal/http-server/middleware"
	postgres "komiac-support-backend/internal/storage"
)

func Register(r *gin.Engine, cfg config.Config, users *postgres.UsersRepo, tickets *postgres.TicketsRepo) {
	r.Use(middleware.CORS(middleware.CORSConfig{Origin: cfg.CorsOrigin}))

	authH := authapi.New(cfg, users)
	ticketsH := ticketsapi.New(tickets)
	usersH := usersapi.New(users)

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

	authMW := middleware.RequireAuth(middleware.AuthConfig{AccessSecret: cfg.AccessSecret})

	u := r.Group("/users", authMW)
	{
		u.GET("/support", usersH.ListSupportUsers)
	}

	t := r.Group("/tickets", authMW)
	{
		t.GET("", ticketsH.ListTickets)
		t.GET("/my", ticketsH.ListMyTickets)
		t.GET("/my/:id", ticketsH.GetMyTicket)
		t.GET("/:id", ticketsH.GetTicket)
		t.POST("", ticketsH.CreateTicket)
		t.POST("/:id/assign", ticketsH.AssignTicket)
		t.POST("/:id/messages", ticketsH.AddMessage)
		t.GET("/:id/messages", ticketsH.ListMessages)
		t.POST("/:id/reply", ticketsH.ReplyTicket)
		t.POST("/:id/close", ticketsH.CloseTicket)
	}
}
