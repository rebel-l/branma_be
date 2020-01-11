package branchstore

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	// ErrIDMissing will be thrown if an ID is expected but not set
	ErrIDMissing = errors.New("id is mandatory for this operation")

	// ErrIDIsSet will be thrown if no ID is expected but already set
	ErrIDIsSet = errors.New("id should be not set for this operation, use update instead")

	// ErrDataMissing will be thrown if mandatory data is not set
	ErrDataMissing = errors.New("no data or mandatory data missing")
)

// Branch represents the branch in the database
type Branch struct {
	ID             int       `db:"id"`
	Name           string    `db:"branch_name"`
	TicketID       string    `db:"ticket_id"`
	ParentTicketID string    `db:"parent_ticket_id"`
	RepositoryID   int       `db:"repository_id"`
	TicketSummary  string    `db:"ticket_summary"`
	TicketStatus   string    `db:"ticket_status"`
	TicketType     string    `db:"ticket_type"`
	Closed         bool      `db:"closed"`
	CreatedAt      time.Time `db:"created_at"`
	ModifiedAt     time.Time `db:"modified_at"`
}

// Create creates current branch in the database
func (b *Branch) Create(ctx context.Context, db *sqlx.DB) error {
	if !b.IsValid() {
		return ErrDataMissing
	}

	if b.ID != 0 {
		return ErrIDIsSet
	}

	q := db.Rebind(`
		INSERT INTO branches (
			ticket_id,
			parent_ticket_id,
			repository_id,
			ticket_summary,
			ticket_status,
			ticket_type,
			branch_name,
			closed
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`)

	res, err := db.ExecContext(ctx, q, b.getCreateArgs()...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	b.ID = int(id)

	return b.Read(ctx, db)
}

// Read sets the branch from database by given ID
func (b *Branch) Read(ctx context.Context, db *sqlx.DB) error {
	if b == nil || b.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`SELECT * FROM branches WHERE id = ?`)

	return db.GetContext(ctx, b, q, b.ID)
}

// Update changes the current branch on the database by ID
func (b *Branch) Update(ctx context.Context, db *sqlx.DB) error {
	if !b.IsValid() {
		return ErrDataMissing
	}

	if b.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`
		UPDATE branches 
		SET ticket_id = ?,
			parent_ticket_id = ?,
			repository_id = ?,
			ticket_summary = ?,
			ticket_status = ?,
			ticket_type = ?,
			branch_name = ?,
			closed = ?
		WHERE id = ?
	`)

	if _, err := db.ExecContext(ctx, q, b.getUpdateArgs()...); err != nil {
		return err
	}

	return b.Read(ctx, db)
}

// Delete removes the current branch from database by its ID
func (b *Branch) Delete(ctx context.Context, db *sqlx.DB) error {
	if b == nil || b.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`DELETE FROM branches WHERE id = ?`)

	_, err := db.ExecContext(ctx, q, b.ID)

	return err
}

// IsValid returns true if all mandatory fields are set
func (b *Branch) IsValid() bool {
	if b == nil || b.Name == "" || b.RepositoryID == 0 {
		return false
	}

	return true
}

func (b *Branch) getCreateArgs() []interface{} {
	return []interface{}{
		b.TicketID,
		b.ParentTicketID,
		b.RepositoryID,
		b.TicketSummary,
		b.TicketStatus,
		b.TicketType,
		b.Name,
		b.Closed,
	}
}

func (b *Branch) getUpdateArgs() []interface{} {
	args := b.getCreateArgs()
	args = append(args, b.ID)

	return args
}
