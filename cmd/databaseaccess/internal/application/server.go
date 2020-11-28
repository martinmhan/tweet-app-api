package application

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/user"
	pb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
)

// DatabaseAccessServer contains the fields and gRPC method implementations used by the DatabaseAccess service
type DatabaseAccessServer struct {
	pb.UnimplementedDatabaseAccessServer
	UserRepository   user.Repository
	FollowRepository follow.Repository
	TweetRepository  tweet.Repository
}

// SaveUser adds a user to the database
func (s *DatabaseAccessServer) SaveUser(ctx context.Context, in *pb.UserConfig) (*pb.InsertID, error) {
	conf := user.Config{Username: in.Username, Password: in.Password}
	i, err := s.UserRepository.Save(conf)
	if err != nil {
		return &pb.InsertID{}, err
	}

	return &pb.InsertID{InsertID: string(i)}, nil
}

// SaveFollow adds a follow (i.e., a unique pair between follower and followee UserIDs) into the database
func (s *DatabaseAccessServer) SaveFollow(ctx context.Context, in *pb.Follow) (*pb.InsertID, error) {
	f := follow.Follow{
		FollowerUserID:   in.FollowerUserID,
		FollowerUsername: in.FollowerUsername,
		FolloweeUserID:   in.FolloweeUserID,
		FolloweeUsername: in.FolloweeUsername,
	}
	insertID, err := s.FollowRepository.Save(f)
	if err != nil {
		return &pb.InsertID{}, err
	}

	return &pb.InsertID{InsertID: insertID}, nil
}

// SaveTweet adds a tweet to the database
func (s *DatabaseAccessServer) SaveTweet(ctx context.Context, in *pb.TweetConfig) (*pb.InsertID, error) {
	conf := tweet.Config{UserID: in.UserID, Username: in.Username, Text: in.Text}
	insertID, err := s.TweetRepository.Save(conf)
	if err != nil {
		return &pb.InsertID{}, err
	}

	return &pb.InsertID{InsertID: insertID}, nil
}

// GetUser gets a user from the database given a UserID
func (s *DatabaseAccessServer) GetUser(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	u, err := s.UserRepository.FindByID(in.UserID)
	if err != nil {
		return &pb.User{}, err
	}

	return &pb.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetFollowers gets the followers of a user from the database given a UserID
func (s *DatabaseAccessServer) GetFollowers(ctx context.Context, in *pb.UserID) (*pb.Follows, error) {
	followers, err := s.FollowRepository.FindFollowersByUserID(in.UserID)
	if err != nil {
		return &pb.Follows{}, err
	}

	var pbFollows []*pb.Follow
	for _, f := range followers {
		pbFollows = append(pbFollows, &pb.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: f.FolloweeUsername,
		})
	}

	return &pb.Follows{Follows: pbFollows}, nil
}

// GetTweets gets the tweets of a user from the database given a UserID
func (s *DatabaseAccessServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	tweets, err := s.TweetRepository.FindByUserID(in.UserID)
	if err != nil {
		return &pb.Tweets{}, err
	}

	var pbTweets []*pb.Tweet
	for _, t := range tweets {
		pbTweets = append(pbTweets, &pb.Tweet{
			ID:       t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}

// GetAllUsers gets all users from the database (only used by the Read View service on cold starts)
func (s *DatabaseAccessServer) GetAllUsers(ctx context.Context, in *pb.GetAllUsersParam) (*pb.Users, error) {
	users, err := s.UserRepository.FindAll()
	if err != nil {
		return &pb.Users{}, err
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
		})
	}

	return &pb.Users{Users: pbUsers}, nil
}

// GetAllFollows gets all follows from the database (only used by the Read View service on cold starts)
func (s *DatabaseAccessServer) GetAllFollows(ctx context.Context, in *pb.GetAllFollowsParam) (*pb.Follows, error) {
	follows, err := s.FollowRepository.FindAll()
	if err != nil {
		return &pb.Follows{}, err
	}

	var pbFollows []*pb.Follow
	for _, f := range follows {
		pbFollows = append(pbFollows, &pb.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   f.FolloweeUserID,
		})
	}

	return &pb.Follows{Follows: pbFollows}, nil
}

// GetAllTweets gets all tweets from the database (only used by the Read View service on cold starts)
func (s *DatabaseAccessServer) GetAllTweets(ctx context.Context, in *pb.GetAllTweetsParam) (*pb.Tweets, error) {
	tweets, err := s.TweetRepository.FindAll()
	if err != nil {
		return &pb.Tweets{}, err
	}

	var pbTweets []*pb.Tweet
	for _, t := range tweets {
		pbTweets = append(pbTweets, &pb.Tweet{
			ID:       t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}
