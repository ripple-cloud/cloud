package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Token struct {
	ID        int64         `db:"id"`
	UserID    int64         `db:"user_id"`
	ExpiresIn time.Duration `db:"expires_in"`

	CreatedAt *time.Time `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}

func (t *Token) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO tokens
	(user_id, expires_in)
	VALUES (:user_id, :expires_in)
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(t).StructScan(t)
	//TODO: handle the possible error cases (like record not unique)
	return err
}

func (t *Token) Get(db *sqlx.DB, id int64) error {
	err := db.Get(t, "SELECT * FROM tokens WHERE id = $1 LIMIT 1;", id)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return &Error{"record_not_found", "token not found"}
	default:
		return err
	}
}