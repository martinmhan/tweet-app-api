package application

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/follower"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/infrastructure/datastore"
	pb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// ReadViewServer implements the gRPC ReadViewServer
type ReadViewServer struct {
	pb.UnimplementedReadViewServer
	Datastore datastore.Datastore
}

// AddUser adds a user to the ReadViewServer's  data store
func (s *ReadViewServer) AddUser(ctx context.Context, in *pb.User) (*pb.SimpleResponse, error) {
	u := user.User{
		ID:       user.ID(in.ID),
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
	t := tweet.Tweet{
		ID:       in.ID,
		UserID:   user.ID(in.UserID),
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
	f := follower.Follower{
		FollowerUserID: user.ID(in.FollowerUserID),
		FolloweeUserID: user.ID(in.FolloweeUserID),
	}
	err := s.Datastore.AddFollower(f)
	if err != nil {
		return &pb.SimpleResponse{Message: "Failed to add follower to read view"}, err
	}

	return &pb.SimpleResponse{Message: "Successfully added follower to read view"}, nil
}

// GetUserByUserID returns the user (if any) of the given UserID
func (s *ReadViewServer) GetUserByUserID(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	u, err := s.Datastore.GetUserByUserID(user.ID(in.UserID))
	if err != nil {
		return &pb.User{}, err
	}

	return &pb.User{
		ID:       string(u.ID),
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
		ID:       string(u.ID),
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetTweets returns the tweets of the given UserID
func (s *ReadViewServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	tweets, err := s.Datastore.GetTweets(user.ID(in.UserID))
	if err != nil {
		return &pb.Tweets{}, err
	}

	pbTweets := []*pb.Tweet{}
	for _, t := range tweets {
		pbTweets = append(pbTweets, &pb.Tweet{
			ID:       t.ID,
			UserID:   string(t.UserID),
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}

// GetTimeline returns the tweets of users that the given UserID follows
func (s *ReadViewServer) GetTimeline(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	timeline, err := s.Datastore.GetTimeline(user.ID(in.UserID))
	if err != nil {
		return &pb.Tweets{}, err
	}

	pbTweets := []*pb.Tweet{}
	for _, t := range timeline {
		pbTweets = append(pbTweets, &pb.Tweet{
			ID:       t.ID,
			UserID:   string(t.UserID),
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}
