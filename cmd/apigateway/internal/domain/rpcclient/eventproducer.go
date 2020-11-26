package rpcclient

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	eventproducerpb "github.com/martinmhan/tweet-app-api/cmd/eventproducer/proto"
)

// EventProducer implements the same methods as the EventProducer gRPC client and provides an abstraction layer for establishing connections and transforming protobuf data types
type EventProducer struct {
	Host string
	Port string
	conn *grpc.ClientConn
}

// Connect establishes a gRPC client connection to the Event Producer service
func (ep *EventProducer) Connect() error {
	target := ep.Host + ":" + ep.Port
	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	ep.conn = conn

	return nil
}

// Disconnect closes the gRPC client connection
func (ep *EventProducer) Disconnect() error {
	if ep.conn == nil {
		return errors.New("EventProducer not connected")
	}

	err := ep.conn.Close()
	if err != nil {
		return err
	}

	ep.conn = nil

	return nil
}

// ProduceUserCreation tells the event producer service via gRPC to publish a Create User event to the message queue
func (ep *EventProducer) ProduceUserCreation(username string, password string) error {
	if ep.conn == nil {
		return errors.New("EventProducer not connected")
	}

	c := eventproducerpb.NewEventProducerClient(ep.conn)
	uf := eventproducerpb.UserConfig{
		Username: username,
		Password: password,
	}

	c.ProduceUserCreation(context.TODO(), &uf)

	return nil
}

// ProduceTweetCreation tells the events producer service via gRPC to publish a Create Tweet event to the message queue
func (ep *EventProducer) ProduceTweetCreation(text string, userID string) error {
	if ep.conn == nil {
		return errors.New("EventProducer not connected")
	}

	c := eventproducerpb.NewEventProducerClient(ep.conn)
	tf := eventproducerpb.TweetConfig{
		UserID: userID,
		Text:   text,
	}

	c.ProduceTweetCreation(context.TODO(), &tf)

	return nil
}
