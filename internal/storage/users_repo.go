package storage

import (
	"context"
	"database/sql"
	"strings"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
	DeptID       *int64
	DeptName     *string
	Role         string
}

type SupportUser struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) GetByLogin(ctx context.Context, login string) (*User, error) {
	const q = `
SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.dept_id, d.name, u.role
FROM users u
LEFT JOIN depts d ON d.id = u.dept_id
WHERE u.username = $1 OR u.email = $1
LIMIT 1;
`
	row := r.db.QueryRowContext(ctx, q, login)

	var u User
	var deptID sql.NullInt64
	var deptName sql.NullString
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone, &deptID, &deptName, &u.Role); err != nil {
		return nil, err
	}
	if deptID.Valid {
		u.DeptID = &deptID.Int64
	}
	if deptName.Valid {
		u.DeptName = &deptName.String
	}
	return &u, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	const q = `
SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.dept_id, d.name, u.role
FROM users u
LEFT JOIN depts d ON d.id = u.dept_id
WHERE u.id = $1
LIMIT 1;
`
	row := r.db.QueryRowContext(ctx, q, id)

	var u User
	var deptID sql.NullInt64
	var deptName sql.NullString
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone, &deptID, &deptName, &u.Role); err != nil {
		return nil, err
	}
	if deptID.Valid {
		u.DeptID = &deptID.Int64
	}
	if deptName.Valid {
		u.DeptName = &deptName.String
	}
	return &u, nil
}

func (r *UsersRepo) CreateSeedSupportIfNotExists(ctx context.Context, username, email, passHash, first, last, phone, deptName string) error {
	var deptID *int64

	if deptName != "" {
		const q = `SELECT id FROM depts WHERE name = $1 LIMIT 1;`
		var id int64
		if err := r.db.QueryRowContext(ctx, q, deptName).Scan(&id); err == nil {
			deptID = &id
		}
	}

	const q = `
INSERT INTO users (username, email, password_hash, first_name, last_name, phone, dept_id, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'support')
ON CONFLICT (username) DO NOTHING;
`
	_, err := r.db.ExecContext(ctx, q, username, email, passHash, first, last, phone, deptID)
	return err
}

func (r *UsersRepo) CreateSeedUserIfNotExists(ctx context.Context, username, email, passHash, first, last, phone, deptName string) error {
	var deptID *int64

	if deptName != "" {
		const q = `SELECT id FROM depts WHERE name = $1 LIMIT 1;`
		var id int64
		if err := r.db.QueryRowContext(ctx, q, deptName).Scan(&id); err == nil {
			deptID = &id
		}
	}

	const q = `
INSERT INTO users (username, email, password_hash, first_name, last_name, phone, dept_id, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'user')
ON CONFLICT (username) DO NOTHING;
`
	_, err := r.db.ExecContext(ctx, q, username, email, passHash, first, last, phone, deptID)
	return err
}

func (r *UsersRepo) ListSupportUsers(ctx context.Context) ([]SupportUser, error) {
	const q = `
SELECT id, first_name, last_name
FROM users
WHERE role = 'support'
ORDER BY last_name, first_name, id;
`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SupportUser
	for rows.Next() {
		var (
			id int64
			fn sql.NullString
			ln sql.NullString
		)
		if err := rows.Scan(&id, &fn, &ln); err != nil {
			return nil, err
		}

		name := strings.TrimSpace(strings.TrimSpace(fn.String) + " " + strings.TrimSpace(ln.String))
		if name == "" {
			name = "—"
		}

		out = append(out, SupportUser{
			ID:   id,
			Name: name,
		})
	}

	return out, rows.Err()
}
