package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

type ReceivedMessage struct {
	ID         int64         `db:"id" json:"id"`
	Topic      string        `db:"topic" json:"topic"`
	Meta       hstore.Hstore `db:"meta" json:"-"`
	Message    []byte        `db:"message" json:"message"`
	HubID      int64         `db:"hub_id" json:"hub_id"`
	ReceivedAt *time.Time    `db:"received_at" json:"received_at"`
}

type ReceivedMessages []ReceivedMessage

func (m *ReceivedMessage) Insert(db *sqlx.DB) error {
	nstmt, err := db.PrepareNamed(`INSERT INTO received_messages 
	(topic, meta, message, hub_id, received_at)
	VALUES (:topic, :meta, :message, :hub_id, now())
	RETURNING *;
	`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.QueryRow(m).StructScan(m)
	return err
}

func (m *ReceivedMessage) Get(db *sqlx.DB, topic string) error {
	err := db.Get(msgs, "SELECT * FROM received_messages WHERE topic = $1 ORDER BY received_at desc LIMIT 1;", topic)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return &Error{"no_messages", "no messages"}
	default:
		return err
	}
}
