package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QType string

const (
	DurableQType   QType = "durable"
	TransientQType QType = "transient"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, value T) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func DeclareAndBind(conn *amqp.Connection, exchange, qname, key string, qtype QType) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	queue, err := ch.QueueDeclare(qname, qtype == DurableQType, qtype == TransientQType, qtype == TransientQType, false, nil)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, err
	}

	if err := ch.QueueBind(qname, key, exchange, false, nil); err != nil {
		ch.QueueDelete(qname, false, false, false)
		ch.Close()
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
