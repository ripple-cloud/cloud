package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int64      `db:"id" json:"id"`
	Username          string     `db:"username" json:"username"`
	Email             string     `db:"email" json:"email"`
	EncryptedPassword string     `db:"encrypted_password" json:"-"`
	CreatedAt         *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updated_at"`
}

// EncryptPassword accepts a password as a string and encrytps it using bcrypt.
// Encrypted value will be stored in user's EncryptedPassword field.
func (u *User) EncryptPassword(password string) error {
	ep, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = string(ep)
	return nil
}

func (u *User) VerifyPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func (u *User) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO users
	(username, email, encrypted_password, created_at, updated_at)
	VALUES (:username, :email, :encrypted_password, now(), now())
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
