package application

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/infrastructure/dbaccess"
	pb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
)

// DatabaseAccessServer contains the fields and gRPC method implementations used by the DatabaseAccess service
type DatabaseAccessServer struct {
	pb.UnimplementedDatabaseAccessServer
	DBAccesser dbaccess.DBAccesser
}

// InsertUser adds a user to the database
func (s *DatabaseAccessServer) InsertUser(ctx context.Context, in *pb.UserConfig) (*pb.InsertID, error) {
	c := dbaccess.UserConfig{Username: in.Username, Password: in.Password}
	i, err := s.DBAccesser.InsertUser(c)
	if err != nil {
		return &pb.InsertID{}, err
	}

	return &pb.InsertID{InsertID: string(i)}, nil
}

// InsertFollower adds a follower (i.e., a unique pair between follower and followee UserIDs) into the database
func (s *DatabaseAccessServer) InsertFollower(ctx context.Context, in *pb.Follower) (*pb.InsertID, error) {
	f := dbaccess.Follower{FollowerUserID: in.FollowerUserID, FolloweeUserID: in.FolloweeUserID}
	i, err := s.DBAccesser.InsertFollower(f)
	if err != nil {
		return &pb.InsertID{}, err
	}

	return &pb.InsertID{InsertID: string(i)}, nil
}

// InsertTweet adds a tweet to the database
func (s *DatabaseAccessServer) InsertTweet(ctx context.Context, in *pb.TweetConfig) (*pb.InsertID, error) {
	c := dbaccess.TweetConfig{UserID: in.UserID, Username: in.Username, Text: in.Text}
	i, err := s.DBAccesser.InsertTweet(c)
	if err != nil {
		return &pb.InsertID{}, err
	}

	return &pb.InsertID{InsertID: string(i)}, nil
}

// GetUser gets a user from the database given a UserID
func (s *DatabaseAccessServer) GetUser(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	i := in.UserID
	u, err := s.DBAccesser.GetUser(i)
	if err != nil {
		return &pb.User{}, err
	}

	return &pb.User{
		UserID:   u.ID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetFollowers gets the followers of a user from the database given a UserID
func (s *DatabaseAccessServer) GetFollowers(ctx context.Context, in *pb.UserID) (*pb.Followers, error) {
	followers, err := s.DBAccesser.GetFollowers(in.UserID)
	if err != nil {
		return &pb.Followers{}, err
	}

	var pbFollowers []*pb.Follower
	for _, f := range followers {
		pbFollowers = append(pbFollowers, &pb.Follower{
			FollowerUserID: f.FollowerUserID,
			FolloweeUserID: f.FolloweeUserID,
		})
	}

	return &pb.Followers{Followers: pbFollowers}, nil
}

// GetTweets gets the tweets of a user from the database given a UserID
func (s *DatabaseAccessServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	tweets, err := s.DBAccesser.GetTweets(in.UserID)
	if err != nil {
		return &pb.Tweets{}, err
	}

	var pbTweets []*pb.Tweet
	for _, t := range tweets {
		pbTweets = append(pbTweets, &pb.Tweet{
			TweetID:  t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}

// GetAllUsers gets all users from the database (only used by the Read View service on cold starts)
func (s *DatabaseAccessServer) GetAllUsers(ctx context.Context, in *pb.GetAllUsersParam) (*pb.Users, error) {
	users, err := s.DBAccesser.GetAllUsers()
	if err != nil {
		return &pb.Users{}, err
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			UserID:   u.ID,
			Username: u.Username,
			Password: u.Password,
		})
	}

	return &pb.Users{Users: pbUsers}, nil
}

// GetAllFollowers gets all followers from the database (only used by the Read View service on cold starts)
func (s *DatabaseAccessServer) GetAllFollowers(ctx context.Context, in *pb.GetAllFollowersParam) (*pb.Followers, error) {
	followers, err := s.DBAccesser.GetAllFollowers()
	if err != nil {
		return &pb.Followers{}, err
	}

	var pbFollowers []*pb.Follower
	for _, f := range followers {
		pbFollowers = append(pbFollowers, &pb.Follower{
			FollowerUserID: f.FollowerUserID,
			FolloweeUserID: f.FolloweeUserID,
		})
	}

	return &pb.Followers{Followers: pbFollowers}, nil
}

// GetAllTweets gets all tweets from the database (only used by the Read View service on cold starts)
func (s *DatabaseAccessServer) GetAllTweets(ctx context.Context, in *pb.GetAllTweetsParam) (*pb.Tweets, error) {
	tweets, err := s.DBAccesser.GetAllTweets()
	if err != nil {
		return &pb.Tweets{}, err
	}

	var pbTweets []*pb.Tweet
	for _, t := range tweets {
		pbTweets = append(pbTweets, &pb.Tweet{
			TweetID:  t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return &pb.Tweets{Tweets: pbTweets}, nil
}
