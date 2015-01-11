package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID                int64      `db:"id"`
	Username          string     `db:"username"`
	Email             string     `db:"email"`
	EncryptedPassword string     `db:"encrypted_password"`
	CreatedAt         *time.Time `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
}

func (u *User) EncryptPassword(password string) error {
	// TODO: implement
	return nil
}

func (u *User) Verify(db *sqlx.DB, password string) bool {
	return false
}

func (u *User) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO users
	(username, email, encrypted_password)
	VALUES (:username, :email, :encrypted_password)
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(u).StructScan(u)
	//TODO: handle the possible error cases (like record not unique)
	return err
}

func (u *User) GetByLogin(db *sqlx.DB, login string) error {
	err := db.Get(u, "SELECT * FROM users WHERE username = $1 OR email = $1 LIMIT 1;", login)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return &Error{"record_not_found", "user not found"}
	default:
		return err
	}
}
