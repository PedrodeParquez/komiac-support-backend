package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"komiac-support-backend/internal/auth"
	"komiac-support-backend/internal/config"
	postgres "komiac-support-backend/internal/storage"
)

const RefreshCookieName = "refresh_token"

type Handlers struct {
	Cfg   config.Config
	Users *postgres.UsersRepo
}

func New(cfg config.Config, users *postgres.UsersRepo) *Handlers {
	return &Handlers{Cfg: cfg, Users: users}
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	u, err := h.Users.GetByLogin(c.Request.Context(), req.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if err := auth.CheckPassword(u.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	role := auth.Role(u.Role)

	access, err := auth.Sign(u.ID, role, h.Cfg.AccessSecret, h.Cfg.AccessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	refresh, err := auth.Sign(u.ID, role, h.Cfg.RefreshSecret, h.Cfg.RefreshTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.SetCookie(
		RefreshCookieName,
		refresh,
		int(h.Cfg.RefreshTTL.Seconds()),
		"/",
		h.Cfg.CookieDomain,
		h.Cfg.CookieSecure,
		true,
	)

	var resp LoginResponse
	resp.AccessToken = access
	resp.User.ID = u.ID
	resp.User.Name = u.FirstName + " " + u.LastName
	resp.User.Role = u.Role
	resp.User.Username = u.Username
	resp.User.Email = u.Email

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) Refresh(c *gin.Context) {
	rt, err := c.Cookie(RefreshCookieName)
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh"})
		return
	}

	claims, err := auth.Parse(rt, h.Cfg.RefreshSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh"})
		return
	}

	access, err := auth.Sign(claims.UID, claims.Role, h.Cfg.AccessSecret, h.Cfg.AccessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{AccessToken: access})
}

func (h *Handlers) Logout(c *gin.Context) {
	c.SetCookie(RefreshCookieName, "", -1, "/", h.Cfg.CookieDomain, h.Cfg.CookieSecure, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handlers) Me(c *gin.Context) {
	uid := c.GetInt("uid")

	u, err := h.Users.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       u.ID,
			"name":     u.FirstName + " " + u.LastName,
			"role":     u.Role,
			"username": u.Username,
			"email":    u.Email,
		},
	})
}
