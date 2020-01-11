package repositorystore

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

// Repository represents the repository in the database
type Repository struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`
	URL        string    `db:"url"`
	CreatedAt  time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"modified_at"`
}

// Create creates current object in the database
func (r *Repository) Create(ctx context.Context, db *sqlx.DB) error {
	if !r.IsValid() {
		return ErrDataMissing
	}

	if r.ID != 0 {
		return ErrIDIsSet
	}

	q := db.Rebind(`INSERT INTO repositories (name, url) VALUES (?, ?)`)

	res, err := db.ExecContext(ctx, q, r.Name, r.URL)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = int(id)

	return r.Read(ctx, db)
}

// Read sets the repository from database by given ID
func (r *Repository) Read(ctx context.Context, db *sqlx.DB) error {
	if r == nil || r.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`SELECT * FROM repositories WHERE id = ?`)

	return db.GetContext(ctx, r, q, r.ID)
}

// Update changes the current object on the database by ID
func (r *Repository) Update(ctx context.Context, db *sqlx.DB) error {
	if !r.IsValid() {
		return ErrDataMissing
	}

	if r.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`UPDATE repositories SET name = ?, url = ? WHERE id = ?`)

	if _, err := db.ExecContext(ctx, q, r.Name, r.URL, r.ID); err != nil {
		return err
	}

	return r.Read(ctx, db)
}

// Delete removes the current object from database by its ID
func (r *Repository) Delete(ctx context.Context, db *sqlx.DB) error {
	if r == nil || r.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`DELETE FROM repositories WHERE id = ?`)

	_, err := db.ExecContext(ctx, q, r.ID)

	return err
}

// IsValid returns true if all mandatory fields are set
func (r *Repository) IsValid() bool {
	if r == nil || r.Name == "" || r.URL == "" {
		return false
	}

	return true
}
