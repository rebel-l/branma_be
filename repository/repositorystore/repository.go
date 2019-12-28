package repositorystore

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrIDMissing   = errors.New("id is not set")
	ErrIDIsSet     = errors.New("id should be not set for this operation, use update instead")
	ErrDataMissing = errors.New("no data or mandatory data missing")
)

type Repository struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`
	URL        string    `db:"url"`
	CreatedAt  time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"modified_at"`
}

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

func (r *Repository) Read(ctx context.Context, db *sqlx.DB) error {
	if r.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`SELECT * FROM repositories WHERE id = ?`)
	return db.GetContext(ctx, r, q, r.ID)
}

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

func (r *Repository) Delete(ctx context.Context, db *sqlx.DB) error {
	if r.ID == 0 {
		return ErrIDMissing
	}

	q := db.Rebind(`DELETE FROM repositories WHERE id = ?`)

	_, err := db.ExecContext(ctx, q, r.Name, r.URL)
	return err
}

func (r *Repository) IsValid() bool {
	if r == nil || r.Name == "" || r.URL == "" {
		return false
	}

	return true
}
