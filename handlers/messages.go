package handlers

import (
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/broker"
	"github.com/ripple-cloud/cloud/data"
	res "github.com/ripple-cloud/cloud/jsonrespond"
	"github.com/ripple-cloud/cloud/router"
	"github.com/ripple-cloud/common/message"
)

func SetBroker(b broker.Broker) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, c router.Context) error {
		c.Meta["broker"] = b
		return c.Next(w, r, c)
	}
}

// POST /send/:topic
// Params hub_id, message
func SendMessage(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}
	bk, ok := c.Meta["upstream"].(broker.Broker)
	if !ok {
		return errors.New("broker not set in context")
	}

	// request must be encoded as a JSON object
	if r.Header.Get("Content-Type") != "application/json" {
		return res.UnsupportedMediaType(w, res.ErrorMsg("invalid_message", "message must be encoded as JSON"))
	}

	topic := c.Params.ByName("topic")
	if topic == "" {
		return res.BadRequest(w, res.ErrorMsg("topic_required", "topic required"))
	}

	// hub slug is optional
	slug := c.Params.ByName("slug")

	msg, err := message.Decode(r.Body)
	if err != nil {
		return res.UnprocessableEntity(w, res.ErrorMsg("message_decode_failed", err))
	}
	msg.Meta["topic"] = topic

	// publish the message to broker
	err := bk.Publish(msg, slug)
	if err != nil {
		return err
	}

	// write it to database
	m := data.SentMessage{
		Topic:   topic,
		Message: msg, //must be converted to hStore
		HubID:   hub_id,
		UserID:  c.Meta["user_id"].(int64),
	}
	err := m.Insert(db)
	if err != nil {
		return err
	}

	res.OK(w, m)
}

// GET /received/:topic
// Params hub_id, message
func LastReceivedMessages(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}

	topic := c.Params.ByName("topic")
	if topic == "" {
		return res.BadRequest(w, res.ErrorMsg("topic_required", "topic required"))
	}

	msg := data.ReceivedMessage{}
	err := msg.Get(db, topic)
	if err != nil {
		return err
	}

	return res.OK(w, msg)
}
