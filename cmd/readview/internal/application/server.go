package application

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/infrastructure/datastore"
	pb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// ReadViewServer implements the gRPC ReadViewServer
type ReadViewServer struct {
	Datastore datastore.Datastore
	pb.UnimplementedReadViewServer
}

// AddUser adds a user to the ReadViewServer's  data store
func (s *ReadViewServer) AddUser(ctx context.Context, in *pb.User) (*pb.SimpleResponse, error) {
	u := datastore.User{
		UserID:   datastore.UserID(in.UserID),
		Username: in.Username,
		Password: in.Password,
	}
	err := s.Datastore.AddUser(u)
	if err != nil {
		return &pb.SimpleResponse{Message: "Failed to add user to read view"}, err
	}

	return &pb.SimpleResponse{}, nil
}

// AddTweet adds a tweet to the ReadViewServer's data store
func (s *ReadViewServer) AddTweet(ctx context.Context, in *pb.Tweet) (*pb.SimpleResponse, error) {
	t := datastore.Tweet{
		TweetID:  in.TweetID,
		UserID:   datastore.UserID(in.UserID),
		Username: in.Username,
		Text:     in.Text,
	}

	err := s.Datastore.AddTweet(t)
	if err != nil {
		return &pb.SimpleResponse{Message: "Failed to add tweet to read view"}, err
	}

	return &pb.SimpleResponse{Message: "Successfully added tweet to read view"}, nil
}

// AddFollower adds a user/follower pair to the ReadViewServer's data store
func (s *ReadViewServer) AddFollower(ctx context.Context, in *pb.Follower) (*pb.SimpleResponse, error) {
	err := s.Datastore.AddFollower(datastore.UserID(in.FollowerUserID), datastore.UserID(in.FollowsUserID))
	if err != nil {
		return &pb.SimpleResponse{Message: "Failed to add follower to read view"}, err
	}

	return &pb.SimpleResponse{Message: "Successfully added follower to read view"}, nil
}

// GetUserByUserID returns the user (if any) of the given UserID
func (s *ReadViewServer) GetUserByUserID(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	u, err := s.Datastore.GetUserByUserID(datastore.UserID(in.UserID))
	if err != nil {
		return &pb.User{}, err
	}

	return &pb.User{
		UserID:   string(u.UserID),
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetUserByUsername returns the user (if any) of the given Username
func (s *ReadViewServer) GetUserByUsername(ctx context.Context, in *pb.Username) (*pb.User, error) {
	u, err := s.Datastore.GetUserByUsername(in.Username)
	if err != nil {
		return &pb.User{}, err
	}

	return &pb.User{
		UserID:   string(u.UserID),
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetTweets returns the tweets of the given UserID
func (s *ReadViewServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	tweets, err := s.Datastore.GetTweets(datastore.UserID(in.UserID))
	if err != nil {
		return &pb.Tweets{}, err
	}

	pbTweets := []*pb.Tweet{}
	for _, t := range tweets {
		pbTweets = append(pbTweets, &pb.Tweet{
			TweetID:  t.TweetID,
			UserID:   string(t.UserID),
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}

// GetTimeline returns the tweets of users that the given UserID follows
func (s *ReadViewServer) GetTimeline(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	timeline, err := s.Datastore.GetTimeline(datastore.UserID(in.UserID))
	if err != nil {
		return &pb.Tweets{}, err
	}

	pbTweets := []*pb.Tweet{}
	for _, t := range timeline {
		pbTweets = append(pbTweets, &pb.Tweet{
			TweetID:  t.TweetID,
			UserID:   string(t.UserID),
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}