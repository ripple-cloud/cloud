package broker

import "github.com/ripple-cloud/common/message"

type Broker interface {
	Publish(msg message.Message, slug string)
}
