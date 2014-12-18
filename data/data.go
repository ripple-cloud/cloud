package data

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"-"`
	Username  string    `json:"username"`
	Email     string    `json:"-"`
	Password  []byte    `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
}

func (user User) AddTo(db *sql.DB) {
	stmt, err := db.Prepare("INSERT into users (username, email, password, token, created_at) values ($1, $2, $3, $4, $5)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, user.Password, user.Token, user.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

func (user User) SetToken(db *sql.DB) {
	stmt, err := db.Prepare("UPDATE users SET token = $1 WHERE username = $2")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Token, user.Username)
	if err != nil {
		log.Fatal(err)
	}
}

func (login *User) GetUserFrom(db *sql.DB) *User {
	rows, err := db.Query("SELECT id, username, email, password, token, created_at FROM users WHERE username = $1", login.Username)
	if err != nil {
		log.Fatal(err)
	}

	user := &User{}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Token, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	return user
}

func GenerateToken() string {
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
