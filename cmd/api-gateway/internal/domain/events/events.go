package events

import (
	"context"
	"errors"

	eppb "github.com/martinmhan/tweet-app-api/cmd/events-producer/proto"

	"google.golang.org/grpc"
)

// ProduceUserCreation calls the events producer service to create a user
func ProduceUserCreation(eventsProducerHost string, eventsProducerPort string, username string, password string) error {
	conn, err := grpc.Dial(eventsProducerHost+":"+eventsProducerPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.New("Could not connect to events producer server")
	}

	defer conn.Close()
	c := eppb.NewEventsProducerClient(conn)

	uf := eppb.UserFields{
		Username: username,
		Password: password,
	}

	c.ProduceUserCreation(context.TODO(), &uf)

	return nil
}

// ProduceTweetCreation calls the events producer service to create a tweet
func ProduceTweetCreation(eventsProducerHost string, eventsProducerPort string, text string, userID string) error {
	conn, err := grpc.Dial(eventsProducerHost+":"+eventsProducerPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.New("Could not connect to events producer server")
	}

	defer conn.Close()
	c := eppb.NewEventsProducerClient(conn)

	tf := eppb.TweetFields{
		UserID: userID,
		Text:   text,
	}

	c.ProduceTweetCreation(context.TODO(), &tf)

	return nil
}
