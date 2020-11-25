package application

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/eventproducer/internal/domain/event"
	pb "github.com/martinmhan/tweet-app-api/cmd/eventproducer/proto"
)

// EventProducerServer is a struct type containing the fields and methods used by the Events Producer service
type EventProducerServer struct {
	pb.UnimplementedEventProducerServer
	event.Producer
}

// ProduceUserCreation publishes a UserCreation event to the message queue
func (s *EventProducerServer) ProduceUserCreation(ctx context.Context, in *pb.UserConfig) (*pb.SimpleResponse, error) {
	e := event.Event{
		Type:    event.UserCreation,
		Payload: in,
	}

	err := s.Produce(e)
	if err != nil {
		return &pb.SimpleResponse{Message: "User creation failed"}, err
	}

	return &pb.SimpleResponse{Message: "User creation accepted"}, nil
}

// ProduceTweetCreation publishes a TweetCreation event to the message queue
func (s *EventProducerServer) ProduceTweetCreation(ctx context.Context, in *pb.TweetConfig) (*pb.SimpleResponse, error) {
	e := event.Event{
		Type:    event.TweetCreation,
		Payload: in,
	}

	err := s.Produce(e)
	if err != nil {
		return &pb.SimpleResponse{Message: "Tweet creation failed"}, err
	}

	return &pb.SimpleResponse{Message: "Tweet creation accepted"}, nil
}
