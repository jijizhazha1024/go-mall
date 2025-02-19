package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

func (m *InventoryMQ) Product(msg *InventoryReq) error {
	channel, err := m.conn.Channel()
	if err != nil {
		return err
	}
	defer func(channel *amqp.Channel) {
		err = channel.Close()
	}(channel)
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return channel.Publish(
		ExchangeName,
		RoutingKeyName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}
