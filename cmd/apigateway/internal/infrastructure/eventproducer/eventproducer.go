package eventproducer

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/user"
	eventproducerpb "github.com/martinmhan/tweet-app-api/cmd/eventproducer/proto"
)

// EventProducer implements the methods of the EventProducer gRPC client utilizing the domain objects
type EventProducer struct {
	eventproducerpb.EventProducerClient
}

// ProduceUserCreation tells the event producer service via gRPC to publish a Create User event to the message queue
func (ep *EventProducer) ProduceUserCreation(u user.Config) error {
	uc := eventproducerpb.UserConfig{Username: u.Username, Password: u.Password}

	_, err := ep.EventProducerClient.ProduceUserCreation(context.TODO(), &uc)
	if err != nil {
		return err
	}

	return nil
}

// ProduceTweetCreation sends a gRPC to the event producer service to publish a CreateTweet event to the message queue
func (ep *EventProducer) ProduceTweetCreation(t tweet.Config) error {
	tc := eventproducerpb.TweetConfig{UserID: t.UserID, Text: t.Text}

	_, err := ep.EventProducerClient.ProduceTweetCreation(context.TODO(), &tc)
	if err != nil {
		return err
	}

	return nil
}

// ProduceFollowCreation sends a gRPC to the event producer service to publish a CreateFollow event to the message queue
func (ep *EventProducer) ProduceFollowCreation(f follow.Config) error {
	fo := eventproducerpb.FollowConfig{
		FollowerUserID: f.FollowerUserID,
		FolloweeUserID: f.FolloweeUserID,
	}

	_, err := ep.EventProducerClient.ProduceFollowCreation(context.TODO(), &fo)
	if err != nil {
		return err
	}

	return nil
}
