package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type TicketsRepo struct {
	db *sql.DB
}

func NewTicketsRepo(db *sql.DB) *TicketsRepo {
	return &TicketsRepo{db: db}
}

func (r *TicketsRepo) ListTickets(ctx context.Context, p ListTicketsParams) ([]TicketListItem, error) {
	where := []string{}
	args := []any{}
	n := 1

	if p.Tab != "" && p.Tab != "all" {
		if p.Tab == "new" {
			where = append(where, fmt.Sprintf("t.status = $%d", n))
			args = append(args, "open")
			n++
		} else {
			where = append(where, fmt.Sprintf("t.status = $%d", n))
			args = append(args, p.Tab)
			n++
		}
	}

	if q := strings.TrimSpace(p.Q); q != "" {
		where = append(where, fmt.Sprintf("(t.ticket_number ILIKE $%d OR t.title ILIKE $%d)", n, n))
		args = append(args, "%"+q+"%")
		n++
	}

	w := ""
	if len(where) > 0 {
		w = "WHERE " + strings.Join(where, " AND ")
	}

	query := `
SELECT
  t.id,
  t.ticket_number,
  t.title,
  t.created_at,
  t.priority,
  t.status,
  u2.first_name,
  u2.last_name
FROM tickets t
JOIN users u ON u.id = t.user_id
LEFT JOIN users u2 ON u2.id = t.taken_by
` + w + `
ORDER BY t.created_at DESC
LIMIT 200;
`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TicketListItem
	for rows.Next() {
		var (
			id       int64
			num      string
			title    string
			created  time.Time
			priority string
			status   string
			fn, ln   sql.NullString
		)

		if err := rows.Scan(&id, &num, &title, &created, &priority, &status, &fn, &ln); err != nil {
			return nil, err
		}

		var assigneeName *string
		if fn.Valid || ln.Valid {
			s := strings.TrimSpace(strings.TrimSpace(fn.String) + " " + strings.TrimSpace(ln.String))
			if s != "" {
				assigneeName = &s
			}
		}

		out = append(out, TicketListItem{
			ID:           id,
			TicketNumber: num,
			Title:        title,
			CreatedAt:    created.Format("15:04 02.01.2006"),
			Priority:     priority,
			Status:       status,
			AssigneeName: assigneeName,
		})
	}

	return out, rows.Err()
}

func (r *TicketsRepo) GetTicket(ctx context.Context, id int64) (TicketDetail, error) {
	query := `
SELECT
  t.id,
  t.ticket_number,
  t.title,
  t.description,
  t.created_at,
  t.priority,
  t.status,
  t.taken_by,
  t.support_reply,
  t.replied_at,
  u.first_name,
  u.last_name,
  d.name,
  u.phone,
  u2.first_name,
  u2.last_name
FROM tickets t
JOIN users u ON u.id = t.user_id
LEFT JOIN depts d ON d.id = u.dept_id
LEFT JOIN users u2 ON u2.id = t.taken_by
WHERE t.id = $1
LIMIT 1;
`

	var (
		t TicketDetail

		created time.Time
		desc    string

		takenBy sql.NullInt64

		supportReply sql.NullString
		repliedAt    sql.NullTime

		uFn, uLn  sql.NullString
		deptName  sql.NullString
		userPhone sql.NullString
		aFn, aLn  sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.TicketNumber,
		&t.Title,
		&desc,
		&created,
		&t.Priority,
		&t.Status,
		&takenBy,
		&supportReply,
		&repliedAt,
		&uFn,
		&uLn,
		&deptName,
		&userPhone,
		&aFn,
		&aLn,
	)
	if err != nil {
		return TicketDetail{}, err
	}

	t.CreatedAt = created.Format("15:04 02.01.2006")
	t.Topic = t.Title
	t.Message = desc

	if supportReply.Valid {
		t.SupportReply = supportReply.String
	}
	if repliedAt.Valid {
		s := repliedAt.Time.Format("15:04 02.01.2006")
		t.RepliedAt = &s
	}

	fromName := strings.TrimSpace(strings.TrimSpace(uFn.String) + " " + strings.TrimSpace(uLn.String))
	if fromName == "" {
		fromName = "—"
	}
	t.FromName = fromName

	if deptName.Valid {
		s := deptName.String
		t.Dept = &s
	}
	if userPhone.Valid {
		s := userPhone.String
		t.Phone = &s
	}

	if takenBy.Valid {
		x := takenBy.Int64
		t.AssigneeID = &x

		name := strings.TrimSpace(strings.TrimSpace(aFn.String) + " " + strings.TrimSpace(aLn.String))
		if name != "" {
			t.AssigneeName = &name
		}
	}

	return t, nil
}

func (r *TicketsRepo) AssignTicket(ctx context.Context, id int64, assigneeID int64) (TicketDetail, error) {
	_, err := r.db.ExecContext(ctx, `
UPDATE tickets
SET taken_by = $1,
    status = CASE WHEN status = 'open' THEN 'in_progress' ELSE status END,
    updated_at = NOW()
WHERE id = $2;
`, assigneeID, id)
	if err != nil {
		return TicketDetail{}, err
	}

	return r.GetTicket(ctx, id)
}

func (r *TicketsRepo) AddMessage(ctx context.Context, p AddMessageParams) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO ticket_messages(ticket_id, author_id, message)
VALUES ($1, $2, $3);
`, p.TicketID, p.AuthorID, p.Message)
	return err
}

func (r *TicketsRepo) ListMyTickets(ctx context.Context, userID int64) ([]TicketListItem, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT
  t.id,
  t.ticket_number,
  t.title,
  t.created_at,
  t.priority,
  t.status,
  u2.first_name,
  u2.last_name
FROM tickets t
LEFT JOIN users u2 ON u2.id = t.taken_by
WHERE t.user_id = $1
ORDER BY t.created_at DESC
LIMIT 200;
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TicketListItem
	for rows.Next() {
		var (
			id       int64
			num      string
			title    string
			created  time.Time
			priority string
			status   string
			fn, ln   sql.NullString
		)

		if err := rows.Scan(&id, &num, &title, &created, &priority, &status, &fn, &ln); err != nil {
			return nil, err
		}

		var assigneeName *string
		if fn.Valid || ln.Valid {
			s := strings.TrimSpace(strings.TrimSpace(fn.String) + " " + strings.TrimSpace(ln.String))
			if s != "" {
				assigneeName = &s
			}
		}

		out = append(out, TicketListItem{
			ID:           id,
			TicketNumber: num,
			Title:        title,
			CreatedAt:    created.Format("15:04 02.01.2006"),
			Priority:     priority,
			Status:       status,
			AssigneeName: assigneeName,
		})
	}

	return out, rows.Err()
}

func (r *TicketsRepo) CreateTicket(ctx context.Context, p CreateTicketParams) (TicketDetail, error) {
	var id int64

	err := r.db.QueryRowContext(ctx, `
INSERT INTO tickets(ticket_number, title, description, status, priority, user_id)
VALUES (
  LPAD(CAST((SELECT COALESCE(MAX(id),0)+1 FROM tickets) AS TEXT), 6, '0'),
  $1, $2, 'open', $3, $4
)
RETURNING id;
`, p.Title, p.Description, p.Priority, p.UserID).Scan(&id)
	if err != nil {
		return TicketDetail{}, err
	}

	return r.GetTicket(ctx, id)
}

func (r *TicketsRepo) ListMessages(ctx context.Context, ticketID int64) ([]TicketMessage, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT
  m.id,
  m.author_id,
  COALESCE(u.first_name,'') as fn,
  COALESCE(u.last_name,'') as ln,
  m.message,
  m.created_at
FROM ticket_messages m
JOIN users u ON u.id = m.author_id
WHERE m.ticket_id = $1
ORDER BY m.created_at ASC, m.id ASC;
`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TicketMessage
	for rows.Next() {
		var (
			id       int64
			authorID int64
			fn       string
			ln       string
			msg      string
			created  time.Time
		)
		if err := rows.Scan(&id, &authorID, &fn, &ln, &msg, &created); err != nil {
			return nil, err
		}

		author := strings.TrimSpace(strings.TrimSpace(fn) + " " + strings.TrimSpace(ln))
		if author == "" {
			author = "—"
		}

		out = append(out, TicketMessage{
			ID:        id,
			AuthorID:  authorID,
			Author:    author,
			Message:   msg,
			CreatedAt: created.Format("15:04 02.01.2006"),
		})
	}
	return out, rows.Err()
}

func (r *TicketsRepo) CloseTicket(ctx context.Context, id int64) (TicketDetail, error) {
	_, err := r.db.ExecContext(ctx, `
UPDATE tickets
SET status = 'closed',
    closed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND status = 'in_progress';
`, id)
	if err != nil {
		return TicketDetail{}, err
	}
	return r.GetTicket(ctx, id)
}

func (r *TicketsRepo) SaveSupportReply(ctx context.Context, ticketID int64, assigneeID int64, reply string) (TicketDetail, error) {
	_, err := r.db.ExecContext(ctx, `
UPDATE tickets
SET taken_by = $1,
    support_reply = $2,
    replied_at = NOW(),
    status = 'in_progress',
    updated_at = NOW()
WHERE id = $3
  AND status IN ('open','in_progress');
`, assigneeID, reply, ticketID)
	if err != nil {
		return TicketDetail{}, err
	}

	return r.GetTicket(ctx, ticketID)
}

func (r *TicketsRepo) GetMyTicket(ctx context.Context, userID int64, id int64) (TicketDetail, error) {
	query := `
SELECT
  t.id,
  t.ticket_number,
  t.title,
  t.description,
  t.created_at,
  t.priority,
  t.status,
  t.taken_by,
  t.support_reply,
  t.replied_at,
  u.first_name,
  u.last_name,
  d.name,
  u.phone,
  u2.first_name,
  u2.last_name
FROM tickets t
JOIN users u ON u.id = t.user_id
LEFT JOIN depts d ON d.id = u.dept_id
LEFT JOIN users u2 ON u2.id = t.taken_by
WHERE t.id = $1 AND t.user_id = $2
LIMIT 1;
`

	var (
		t TicketDetail

		created time.Time
		desc    string

		takenBy sql.NullInt64

		supportReply sql.NullString
		repliedAt    sql.NullTime

		uFn, uLn  sql.NullString
		deptName  sql.NullString
		userPhone sql.NullString
		aFn, aLn  sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&t.ID,
		&t.TicketNumber,
		&t.Title,
		&desc,
		&created,
		&t.Priority,
		&t.Status,
		&takenBy,
		&supportReply,
		&repliedAt,
		&uFn,
		&uLn,
		&deptName,
		&userPhone,
		&aFn,
		&aLn,
	)
	if err != nil {
		return TicketDetail{}, err
	}

	t.CreatedAt = created.Format("15:04 02.01.2006")
	t.Topic = t.Title
	t.Message = desc

	if supportReply.Valid {
		t.SupportReply = supportReply.String
	}
	if repliedAt.Valid {
		s := repliedAt.Time.Format("15:04 02.01.2006")
		t.RepliedAt = &s
	}

	fromName := strings.TrimSpace(strings.TrimSpace(uFn.String) + " " + strings.TrimSpace(uLn.String))
	if fromName == "" {
		fromName = "—"
	}
	t.FromName = fromName

	if deptName.Valid {
		s := deptName.String
		t.Dept = &s
	}
	if userPhone.Valid {
		s := userPhone.String
		t.Phone = &s
	}

	if takenBy.Valid {
		x := takenBy.Int64
		t.AssigneeID = &x

		name := strings.TrimSpace(strings.TrimSpace(aFn.String) + " " + strings.TrimSpace(aLn.String))
		if name != "" {
			t.AssigneeName = &name
		}
	}

	return t, nil
}
