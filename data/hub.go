package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Hub struct {
	ID        int64      `db:"id" json:"id"`
	Slug      string     `db:"slug" json:"slug"`
	UserID    int64      `db:"user_id" json:"user_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

type Hubs []string

func (h *Hub) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO hubs 
	(slug, user_id, created_at, updated_at)
	VALUES (:slug, :user_id, now(), now())
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(h).StructScan(h)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		case "unique_violation":
			return &Error{"unique_violation", "hub exists"}
		default:
			return &Error{err.Code.Name(), "pq error"}
		}
	}
	return err
}

func (h *Hubs) SelectByUserId(db *sqlx.DB, userid int64) error {
	err := db.Select(h, "SELECT slug FROM hubs WHERE user_id = $1;", userid)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		default:
			return &Error{err.Code.Name(), "pq error"}
		}
	}

	if *h == nil {
		return &Error{"record_not_found", "hub not found"}
	}
	return err
}

func (h *Hub) Get(db *sqlx.DB, slug string) error {
	err := db.Get(h, "SELECT user_id FROM hubs WHERE slug = $1;", slug)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		default:
			return &Error{err.Code.Name(), "pq error"}
		}
	}

	if err == sql.ErrNoRows {
		return &Error{"record_not_found", "hub not found"}
	}
	return err
}

func (h *Hub) Delete(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`DELETE FROM hubs
	WHERE slug = (:slug)
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(h).StructScan(h)
	// TODO?: Handle possible error cases? Eg...
	// hub record_not_found error would have already been caught in Get Method above.
	return err
}
