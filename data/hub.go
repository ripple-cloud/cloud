package data

import (
	//	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Hub struct {
	ID        int64      `db:"id" json:"id"`
	Hub       string     `db:"hub" json:"hub"`
	UserID    int64      `db:"user_id" json:"user_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (h *Hub) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO hubs 
	(hub, user_id, created_at, updated_at)
	VALUES (:hub, :user_id, now(), now())
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

// func (hub *Hub) Add(db *sql.DB) {
// 	stmt, err := db.Prepare("INSERT into hubs (hub, user_id) values ($1, $2)")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec(hub.Hub, hub.UserID)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
//
// func (hub *Hub) GetByUser(db *sql.DB) []*Hub {
// 	rows, err := db.Query("SELECT id, hub, user_id FROM hubs WHERE user = $1", hub.UserID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
//
// 	var res []*Hub
// 	for rows.Next() {
// 		h := &Hub{}
// 		err := rows.Scan(&h.ID, &h.Hub, &h.UserID)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		res = append(res, h)
// 	}
//
// 	return res
// }
//
// func (hub *Hub) GetByHub(db *sql.DB) []*Hub {
// 	rows, err := db.Query("SELECT id, hub, user_id FROM hubs WHERE id = $1", hub.Hub)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
//
// 	var res []*Hub
// 	for rows.Next() {
// 		h := &Hub{}
// 		err := rows.Scan(&h.ID, &h.Hub, &h.UserID)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		res = append(res, h)
// 	}
//
// 	return res
// }
//
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
