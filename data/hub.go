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

func (hub *Hub) GetByUser(db *sql.DB) []*Hub {
	rows, err := db.Query("SELECT id, hub, user_id FROM hubs WHERE user = $1", hub.UserID)
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

func (hub *Hub) GetByHub(db *sql.DB) []*Hub {
	rows, err := db.Query("SELECT id, hub, user_id FROM hubs WHERE id = $1", hub.Hub)
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

func (hub *Hub) Delete(db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM hubs WHERE id = $1", hub.ID)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(value)
	if err != nil {
		fmt.Println(err)
	}
}
