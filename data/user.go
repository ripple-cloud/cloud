package data

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  []byte
	Token     string
	CreatedAt time.Time
}

func (user User) Add(db *sql.DB) {
	stmt, err := db.Prepare("INSERT into users (username, email, password, token, created_at) values ($1, $2, $3, $4, $5)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, user.Password, user.Token, user.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
}

func (user *User) GetByID(db *sql.DB) *User {
	rows, err := db.Query("SELECT id, username, email, password, token, created_at FROM users WHERE id = $1", user.ID)
	if err != nil {
		fmt.Println(err)
	}

	u := &User{}
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Token, &u.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
	}
	return u
}

func (user *User) GetByToken(db *sql.DB) *User {
	rows, err := db.Query("SELECT id, username, email, password, token, created_at FROM users WHERE token = $1", user.Token)
	if err != nil {
		fmt.Println(err)
	}

	u := &User{}
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Token, &u.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
	}
	return u
}

func (user User) SetToken(db *sql.DB) {
	stmt, err := db.Prepare("UPDATE users SET token = $1 WHERE username = $2")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(generateToken(), user.Username)
	if err != nil {
		fmt.Println(err)
	}
}

// TODO: base64?
func generateToken() string {
	rand.Seed(time.Now().UTC().UnixNano())

	char := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 64)
	for i := range b {
		b[i] = char[rand.Intn(len(char))]
	}

	return string(b)
}

func Encrypt(plaintext string) []byte {
	cryptext, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return cryptext
}
