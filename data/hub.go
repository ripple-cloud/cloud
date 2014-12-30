package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Hub struct {
	ID     string
	Hub    string
	UserID string
}

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
