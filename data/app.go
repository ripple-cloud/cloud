package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type App struct {
	ID        int64      `db:"id" json:"id"`
	Slug      string     `db:"slug" json:"slug"`
	HubID     int64      `db:"hub_id" json:"hub_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (a *App) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO apps
	(slug, hub_id, created_at, updated_at)
	VALUES (:slug, :hub_id, now(), now())
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(a).StructScan(a)
	//TODO: handle the possible error cases (like record not unique)
	return err
}

func (a *App) Get(db *sqlx.DB, slug string, hub_id int64) error {
	err := db.Get(a, "SELECT * FROM users WHERE slug = $1 AND hub_id = $2 LIMIT 1;", slug, hub_id)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return &Error{"record_not_found", "app not found"}
	default:
		return err
	}
}

func (a *App) Delete(db *sqlx.DB) error {
	_, err := db.Exec(`DELETE FROM app WHERE id = $1`, a.ID)
	return err
}
