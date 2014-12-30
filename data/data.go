package data

import (
	"database/sql"
	"fmt"
	"log"
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

type Hub struct {
	ID     string
	Hub    string
	UserID string
}

// User Table.

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

func (user User) SetToken(db *sql.DB, token, username string) {
	stmt, err := db.Prepare("UPDATE users SET token = $1 WHERE username = $2")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(generateToken(), username)
	if err != nil {
		fmt.Println(err)
	}
}

func (user *User) Get(db *sql.DB, col, value string) *User {
	rows, err := db.Query(fmt.Sprintf("SELECT id, username, email, password, token, created_at FROM users WHERE %s = $1", col), value)
	if err != nil {
		log.Fatal(err)
	}

	u := &User{}
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Token, &u.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	return u
}

// Hub table.

func (hub *Hub) Add(db *sql.DB) {
	stmt, err := db.Prepare("INSERT into hubs (hub, user_id) values ($1, $2)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(hub.Hub, hub.UserID)
	if err != nil {
		fmt.Println(err)
	}
}

func (hub *Hub) Get(db *sql.DB, col, value string) []*Hub {
	rows, err := db.Query(fmt.Sprintf("SELECT id, hub, user_id FROM hubs WHERE %s = $1", col), value)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var res []*Hub
	for rows.Next() {
		h := &Hub{}
		err := rows.Scan(&h.ID, &h.Hub, &h.UserID)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, h)
	}

	return res
}

func (hub *Hub) Delete(db *sql.DB, col, value string) {
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM hubs WHERE %s = $1", col))
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(value)
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
