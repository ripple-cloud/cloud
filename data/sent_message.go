package data

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SentMessage struct {
	ID          int64      `db:"id" json:"id"`
	Topic       string     `db:"topic" json:"topic"`
	Message     []byte     `db:"message" json:"message"`
	HubID       int64      `db:"hub_id" json:"hub_id"`
	UserID      int64      `db:"user_id" json:"user_id"`
	RequestedAt *time.Time `db:"requested_at" json:"requested_at"`
}

func (m *SentMessage) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO sent_messages 
	(topic, message, hub_id, user_id, requested_at)
	VALUES (:topic, :message, :hub_id, :user_id, now())
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(m).StructScan(m)
	return err
}
