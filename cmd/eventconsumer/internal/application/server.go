package application

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"

	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/user"
)

// EventConsumerServer listens for and executes events from the message queue
type EventConsumerServer struct {
	Connection       *amqp.Connection
	MessageQueueName string
	UserRepository   user.Repository
	FollowRepository follow.Repository
	TweetRepository  tweet.Repository
}

func (e *EventConsumerServer) createUser(eventPayload []byte) error {
	var conf user.Config

	err := json.Unmarshal(eventPayload, &conf)
	if err != nil {
		return err
	}

	_, err = e.UserRepository.Save(conf)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventConsumerServer) createFollow(eventPayload []byte) error {
	var f follow.Config

	err := json.Unmarshal(eventPayload, &f)
	if err != nil {
		return err
	}

	err = e.FollowRepository.Save(f)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventConsumerServer) createTweet(eventPayload []byte) error {
	var conf tweet.Config

	err := json.Unmarshal(eventPayload, &conf)
	if err != nil {
		return err
	}

	_, err = e.TweetRepository.Save(conf)
	if err != nil {
		return err
	}

	return nil
}

// Listen starts the EventConsumerServer so that it continually listens for new events to process from the message queue
func (e *EventConsumerServer) Listen() error {
	ch, err := e.Connection.Channel()
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
			case "UserCreation":
				e.createUser(d.Body)
			case "TweetCreation":
				e.createTweet(d.Body)
			case "FollowCreation":
				e.createFollow(d.Body)
			}
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
