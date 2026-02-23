package storage

type ListTicketsParams struct {
	Tab string
	Q   string
}

type TicketListItem struct {
	ID           int64   `json:"id"`
	TicketNumber string  `json:"ticketNumber"`
	Title        string  `json:"title"`
	CreatedAt    string  `json:"createdAt"`
	Priority     string  `json:"priority"`
	Status       string  `json:"status"`
	AssigneeName *string `json:"assigneeName,omitempty"`
}

type TicketDetail struct {
	ID           int64   `json:"id"`
	TicketNumber string  `json:"ticketNumber"`
	Title        string  `json:"title"`
	CreatedAt    string  `json:"createdAt"`
	Priority     string  `json:"priority"`
	Status       string  `json:"status"`
	AssigneeID   *int64  `json:"assigneeId,omitempty"`
	AssigneeName *string `json:"assigneeName,omitempty"`
	SupportReply string  `json:"supportReply"`
	RepliedAt    *string `json:"repliedAt,omitempty"`
	Topic        string  `json:"topic"`
	FromName     string  `json:"fromName"`
	Dept         *string `json:"dept,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Message      string  `json:"message"`
}

type AddMessageParams struct {
	TicketID int64
	AuthorID int64
	Message  string
}

type CreateTicketParams struct {
	Title       string
	Description string
	Priority    string
	UserID      int64
}

type TicketMessage struct {
	ID        int64  `json:"id"`
	AuthorID  int64  `json:"authorId"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}
