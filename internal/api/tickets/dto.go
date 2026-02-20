package tickets

type AssignTicketRequest struct {
	AssigneeID int64 `json:"assigneeId"`
}

type AddMessageRequest struct {
	Message string `json:"message"`
}

type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}
