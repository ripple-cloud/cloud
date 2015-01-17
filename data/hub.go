package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	//TODO: handle the possible error cases (like record not unique)
	return err
}

func (h *Hubs) GetByUserId(db *sqlx.DB, userid int64) error {
	err := db.Select(h, "SELECT slug FROM hubs WHERE user_id = $1;", userid)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return &Error{"record_not_found", "hub not found"}
	default:
		return err
	}
}

// func (hub *Hub) Delete(db *sql.DB) {
// 	stmt, err := db.Prepare("DELETE FROM hubs WHERE id = $1", hub.ID)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec(value)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
