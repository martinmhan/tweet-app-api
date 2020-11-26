package application

import (
	"encoding/json"
	"log"

	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/rpcclient"
	"github.com/streadway/amqp"
)

// EventConsumer listens for and executes events from the message queue
type EventConsumer struct {
	MessageQueueHost string
	MessageQueuePort string
	MessageQueueName string
	rpcclient.DatabaseAccess
	rpcclient.ReadView
}

func (e *EventConsumer) createUser(eventPayload []byte) error {
	var uc rpcclient.UserConfig

	err := json.Unmarshal(eventPayload, &uc)
	if err != nil {
		return err
	}

	id, err := e.InsertUser(uc) // creates user in db
	if err != nil {
		return err
	}

	u := rpcclient.User{ID: id, Username: uc.Username, Password: uc.Password}
	err = e.AddUser(u) // adds created user to read view
	if err != nil {
		return err
	}

	return nil
}

func (e *EventConsumer) createTweet(eventPayload []byte) error {
	var tc rpcclient.TweetConfig

	err := json.Unmarshal(eventPayload, &tc)
	if err != nil {
		return err
	}

	id, err := e.InsertTweet(tc) // creates user in db
	if err != nil {
		return err
	}

	tweet := rpcclient.Tweet{
		TweetID:  id,
		UserID:   tc.UserID,
		Username: tc.Username,
		Text:     tc.Text,
	}
	err = e.AddTweet(tweet) // adds created user to read view
	if err != nil {
		return err
	}

	return nil
}

func (e *EventConsumer) connect() (*amqp.Connection, error) {
	url := "amqp://guest:guest@" + e.MessageQueueHost + ":" + e.MessageQueuePort + "/"
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Init starts the EventConsumer so that it continually listens for new events to process from the message queue
func (e *EventConsumer) Init() error {
	conn, err := e.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		e.MessageQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println("Received a message")
			log.Printf("Message Type: %s", d.Type)
			log.Printf("Message Body: %s", d.Body)

			switch d.Type {
			case "CreateUser":
				e.createUser(d.Body)
			case "CreateTweet":
				e.createTweet(d.Body)
			}
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
