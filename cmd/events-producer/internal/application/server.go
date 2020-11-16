package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/martinmhan/crud-api-golang-grpc/utils"
	event "github.com/martinmhan/tweet-app-api/cmd/events-producer/internal/domain"
	pb "github.com/martinmhan/tweet-app-api/cmd/events-producer/proto"
	"github.com/streadway/amqp"
)

// EventsProducerServer ...
type EventsProducerServer struct {
	pb.UnimplementedEventsProducerServer
	MessageQueueHost string
	MessageQueuePort string
}

func (s *EventsProducerServer) connect() (*amqp.Connection, error) {
	url := "amqp://guest:guest@" + s.MessageQueueHost + ":" + s.MessageQueuePort + "/"

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errors.New("Failed to connect to RabbitMQ")
	}

	return conn, nil
}

func (s *EventsProducerServer) produceEvent(e event.Event) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"crud",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")

	b := &e.Payload
	body, err := json.Marshal(b)
	utils.FailOnError(err, "Failed to read message body")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Type:        e.Type,
			Body:        body,
		},
	)

	return err
}

// ProduceUserCreation ...
func (s *EventsProducerServer) ProduceUserCreation(ctx context.Context, in *pb.UserFields) (*pb.SimpleResponse, error) {
	// TODO

	// s.produceEvent()

	return &pb.SimpleResponse{
		Message: "User creataion event produced",
	}, nil
}

// ProduceTweetCreation ...
func (s *EventsProducerServer) ProduceTweetCreation(ctx context.Context, in *pb.TweetFields) (*pb.SimpleResponse, error) {
	// TODO

	// s.produceEvent()

	return &pb.SimpleResponse{
		Message: "Tweet creation event produced",
	}, nil
}
