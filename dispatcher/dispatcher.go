package dispatcher

import (
	"bytes"
	"log"
	"regexp"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/common/message"
	"github.com/ripple-other/broker"
)

// collects messages from broker and stores them in database

var hubRegexp *regexp.Regexp

func init() {
	hubRegexp = regexp.MustCompile(`^(?:/data/hub/)([^/]+)/(.+)$`) // eg: /data/hub/myhub/my-topic
}

func storeMessage(db *sqlx.DB, msg mqtt.Message) {
	topic := msg.Topic()
	hub := hubRegexp.FindStringSubmatch(topic)[1]
	if hub == "" {
		// invalid message, ignore
		log.Printf("[info] dispatcher: message received without a hub %s", topic)
		return
	}
	// find hub by slug

	rMsg := msg.Payload()
	// decode the message from payload
	dMsg, err := message.Decode(bytes.NewReader(rMsg))
	if err != nil {
		log.Printf("[error] dispatcher: failed to decode the message %v err: %s", err)
	}
	// get and store the metadata
	meta := dMsg.Meta

	m := data.ReceivedMessage{
		Topic:   topic,
		Meta:    dMsg.Meta,
		Message: rMsg,
		HubID:   hubId,
	}
	if err := m.Insert(db); err != nil {
		log.Printf("[error] dispatcher: failed to insert message: %s", err)
	}
}

func Start(db *sqlx.DB, b broker.Broker) error {
	// subscribe to broker
	tf, err := mqtt.NewTopicFilter("data/hub/+/+", byte(mqtt.QOS_ZERO))
	if err != nil {
		return err
	}

	rcpt, err := up.client.StartSubscription(func(_ *mqtt.MqttClient, msg mqtt.Message) {
		storeMessage(db, msg)
	}, tf)
	if err != nil {
		return err
	}

	// return only after receiving the receipt
	<-rcpt
	return nil
}
