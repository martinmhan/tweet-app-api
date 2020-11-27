package eventproducer

import (
	"encoding/json"
	"errors"

	"github.com/streadway/amqp"

	"github.com/martinmhan/tweet-app-api/cmd/eventproducer/internal/domain/event"
)

// EventProducer produces events by publishing message queue
type EventProducer struct {
	MessageQueueName string
	Connection       *amqp.Connection
}

// Produce publishes an event to the message queue
func (p *EventProducer) Produce(e event.Event) error {
	if p.Connection.IsClosed() {
		return errors.New("EventProducer is not connected")
	}

	if e.Type.String() == "" {
		return errors.New("Invalid Event Type")
	}

	ch, err := p.Connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		p.MessageQueueName, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	b := &e.Payload
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Type:        e.Type.String(),
			Body:        body,
		},
	)

	return err
}
