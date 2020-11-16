package application

import (
	"context"
	"fmt"

	"github.com/martinmhan/tweet-app-api/cmd/database-access/internal/domain/accesser"
	pb "github.com/martinmhan/tweet-app-api/cmd/database-access/proto"
)

// DatabaseAccessServer ...
type DatabaseAccessServer struct {
	pb.UnimplementedDatabaseAccessServer
	DBAccesser accesser.DBAccesser
}

// InsertUser ...
func (s *DatabaseAccessServer) InsertUser(ctx context.Context, in *pb.User) (*pb.InsertID, error) {
	fmt.Println("Endpoint hit: InsertUser")

	f := accesser.UserFields{
		Username: in.Username,
	}

	u := s.DBAccesser.InsertUser(f)

	r := pb.InsertID{
		InsertID: u.ID,
	}

	return &r, nil
}

// InsertTweet ...
func (s *DatabaseAccessServer) InsertTweet(ctx context.Context, in *pb.Tweet) (*pb.InsertID, error) {
	fmt.Println("Endpoint hit: InsertTweet")

	f := accesser.TweetFields{
		Text: in.Text,
	}

	t := s.DBAccesser.InsertTweet(f)

	r := pb.InsertID{
		InsertID: t.ID,
	}

	return &r, nil
}

// GetUser ...
func (s *DatabaseAccessServer) GetUser(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	fmt.Println("Endpoint hit: GetUser")

	i := in.UserID

	u := s.DBAccesser.GetUser(i)

	r := pb.User{
		Username: u.Username,
		Password: u.Password,
	}

	return &r, nil
}

// GetTweets ...
func (s *DatabaseAccessServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	fmt.Println("Endpoint hit: GetTweets")

	records := s.DBAccesser.GetTweets(in.UserID)

	var tweets []*pb.Tweet
	for _, t := range records {
		tweets = append(tweets, &pb.Tweet{
			ID:   t.ID,
			Text: t.Text,
		})
	}

	r := pb.Tweets{
		Tweets: tweets,
	}

	return &r, nil
}
