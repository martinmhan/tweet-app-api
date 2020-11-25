package event

import (
	"encoding/json"
	"errors"

	"github.com/streadway/amqp"
)

// Event is a generic struct type passed to the message queue containing a Type and Payload
type Event struct {
	Type    Type
	Payload interface{}
}

// Type specifies the type of event to pass to the message queue
type Type int

const (
	// UserCreation is an event type
	UserCreation Type = iota
	// TweetCreation is an event type
	TweetCreation
)

func (t Type) String() string {
	types := [...]string{
		"UserCreation",
		"TweetCreation",
	}

	return types[t]
}

// Producer contains methods used to publish events to the message queue
type Producer struct {
	MessageQueueHost string
	MessageQueuePort string
	MessageQueueName string
	conn             *amqp.Connection
}

// Connect establishes a RabbitMQ connection
func (p *Producer) Connect() error {
	url := "amqp://guest:guest@" + p.MessageQueueHost + ":" + p.MessageQueuePort + "/"
	conn, err := amqp.Dial(url)
	if err != nil {
		return errors.New("Failed to connect to RabbitMQ")
	}

	p.conn = conn

	return nil
}

// Disconnect closes the RabbitMQ connection
func (p *Producer) Disconnect() error {
	if p.conn == nil {
		return errors.New("Producer is not connected")
	}

	err := p.conn.Close()
	if err != nil {
		return err
	}

	p.conn = nil

	return nil
}

// Produce publishes an event to the message queue
func (p *Producer) Produce(e Event) error {
	if p.conn == nil {
		return errors.New("Producer is not connected")
	}

	if e.Type.String() == "" {
		return errors.New("Invalid Event Type")
	}

	ch, err := p.conn.Channel()
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
