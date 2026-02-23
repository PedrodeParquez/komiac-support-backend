package users

import (
	"net/http"

	"komiac-support-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	users *storage.UsersRepo
}

func New(users *storage.UsersRepo) *Handlers {
	return &Handlers{users: users}
}

func roleFromCtx(c *gin.Context) string {
	v, _ := c.Get("role")
	s, _ := v.(string)
	return s
}

func (h *Handlers) ListSupportUsers(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	items, err := h.users.ListSupportUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if items == nil {
		items = make([]storage.SupportUser, 0)
	}

	c.JSON(http.StatusOK, gin.H{"users": items})
}
