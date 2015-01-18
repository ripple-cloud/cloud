package data

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Token struct {
	ID        int64 `db:"id"`
	UserID    int64 `db:"user_id"`
	ExpiresIn int64 `db:"expires_in"`

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
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		default:
			return &Error{err.Code.Name(), "pq error"}
		}
	}
	return err
}

func (t *Token) Get(db *sqlx.DB, id int64) error {
	err := db.Get(t, "SELECT * FROM tokens WHERE id = $1 LIMIT 1;", id)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		default:
			return &Error{err.Code.Name(), "pq error"}
		}
	}

	if err == sql.ErrNoRows {
		return &Error{"record_not_found", "token not found"}
	}
	return err
}

// Encode JWT will return the current token encoded as a JSON web token.
// Note the encoded token is not persisted
func (t *Token) EncodeJWT(tokenSecret []byte) (string, error) {
	j := jwt.New(jwt.SigningMethodHS256)
	j.Claims["iat"] = t.CreatedAt.Unix()                                 // issued at
	j.Claims["exp"] = t.CreatedAt.Add(time.Duration(t.ExpiresIn)).Unix() // expires at
	j.Claims["jti"] = t.ID                                               // token ID
	j.Claims["user_id"] = t.UserID
	//j.Claims["scopes"] = "user,hub,app" // FIXME: should not be hardcoded
	return j.SignedString(tokenSecret)
}
