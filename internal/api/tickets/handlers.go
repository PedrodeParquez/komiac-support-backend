package tickets

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"komiac-support-backend/internal/storage"
)

type Handlers struct {
	tickets *storage.TicketsRepo
}

func New(tickets *storage.TicketsRepo) *Handlers {
	return &Handlers{tickets: tickets}
}

func uidFromCtx(c *gin.Context) (int64, bool) {
	v, ok := c.Get("uid")
	if !ok {
		return 0, false
	}

	switch x := v.(type) {
	case int:
		return int64(x), true
	case int64:
		return x, true
	default:
		return 0, false
	}
}

func roleFromCtx(c *gin.Context) string {
	v, _ := c.Get("role")
	s, _ := v.(string)
	return s
}

func (h *Handlers) ListTickets(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	items, err := h.tickets.ListTickets(c.Request.Context(), storage.ListTicketsParams{
		Tab: c.Query("tab"),
		Q:   c.Query("q"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if items == nil {
		items = make([]storage.TicketListItem, 0)
	}

	c.JSON(http.StatusOK, gin.H{"tickets": items})
}

func (h *Handlers) ListMyTickets(c *gin.Context) {
	uid, ok := uidFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := h.tickets.ListMyTickets(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if items == nil {
		items = make([]storage.TicketListItem, 0)
	}

	c.JSON(http.StatusOK, gin.H{"tickets": items})
}

func (h *Handlers) GetTicket(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	t, err := h.tickets.GetTicket(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": t})
}

func (h *Handlers) AssignTicket(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.AssigneeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	t, err := h.tickets.AssignTicket(c.Request.Context(), id, req.AssigneeID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": t})
}

func (h *Handlers) AddMessage(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	uid, ok := uidFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	var req AddMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	if err := h.tickets.AddMessage(c.Request.Context(), storage.AddMessageParams{
		TicketID: id,
		AuthorID: uid,
		Message:  strings.TrimSpace(req.Message),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handlers) CreateTicket(c *gin.Context) {
	if roleFromCtx(c) != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	uid, ok := uidFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Priority = strings.TrimSpace(req.Priority)

	if req.Title == "" || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title/description required"})
		return
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	t, err := h.tickets.CreateTicket(c.Request.Context(), storage.CreateTicketParams{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		UserID:      uid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ticket": t})
}

func (h *Handlers) ListMessages(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	items, err := h.tickets.ListMessages(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if items == nil {
		items = make([]storage.TicketMessage, 0)
	}

	c.JSON(http.StatusOK, gin.H{"messages": items})
}

func (h *Handlers) ReplyTicket(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	var req ReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	req.Reply = strings.TrimSpace(req.Reply)
	if req.AssigneeID <= 0 || req.Reply == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assigneeId and reply required"})
		return
	}

	t, err := h.tickets.SaveSupportReply(c.Request.Context(), id, req.AssigneeID, req.Reply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": t})
}

func (h *Handlers) CloseTicket(c *gin.Context) {
	if roleFromCtx(c) != "support" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	t, err := h.tickets.CloseTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": t})
}
