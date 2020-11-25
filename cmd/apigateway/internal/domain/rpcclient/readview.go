package rpcclient

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// ReadView implements the same methods as the ReadView gRPC client and provides an abstraction layer for establishing connections and transforming protobuf data types
type ReadView struct {
	Host string
	Port string
	conn *grpc.ClientConn
}

// Connect creates a connection to the Read View service
func (rv *ReadView) Connect() error {
	conn, err := grpc.Dial(rv.Host+":"+rv.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	rv.conn = conn

	return nil
}

// Disconnect closes the gRPC client connection
func (rv *ReadView) Disconnect() error {
	if rv.conn == nil {
		return errors.New("ReadView not connected")
	}

	err := rv.conn.Close()
	if err != nil {
		return err
	}

	rv.conn = nil

	return nil
}

// GetUserByUsername TO DO
func (rv *ReadView) GetUserByUsername(username string) (User, error) {
	if rv.conn == nil {
		return User{}, errors.New("ReadView not connected")
	}

	c := readviewpb.NewReadViewClient(rv.conn)
	un := readviewpb.Username{Username: username}
	u, err := c.GetUserByUsername(context.TODO(), &un)
	if err != nil {
		return User{}, err
	}

	return User{
		UserID:   u.UserID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetUserByUserID TO DO
func (rv *ReadView) GetUserByUserID(userID string) (User, error) {
	if rv.conn == nil {
		return User{}, errors.New("ReadView not connected")
	}

	c := readviewpb.NewReadViewClient(rv.conn)
	uid := readviewpb.UserID{UserID: userID}
	u, err := c.GetUserByUserID(context.TODO(), &uid)
	if err != nil {
		return User{}, err
	}

	return User{
		UserID:   u.UserID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetTweets TO DO
func (rv *ReadView) GetTweets(userID string) (Tweets, error) {
	if rv.conn == nil {
		return Tweets{}, errors.New("ReadView not connected")
	}

	c := readviewpb.NewReadViewClient(rv.conn)
	uid := readviewpb.UserID{UserID: userID}
	pbtweets, err := c.GetTweets(context.TODO(), &uid)
	if err != nil {
		return Tweets{}, err
	}

	tweets := Tweets{}
	for i, t := range pbtweets.Tweets {
		tweets[i] = Tweet{
			UserID: t.UserID,
			Text:   t.Text,
		}
	}

	return tweets, nil
}

// GetTimeline TO DO
func (rv *ReadView) GetTimeline(userID string) (Tweets, error) {
	if rv.conn == nil {
		return Tweets{}, errors.New("ReadView not connected")
	}

	c := readviewpb.NewReadViewClient(rv.conn)
	uid := readviewpb.UserID{UserID: userID}
	pbtweets, err := c.GetTimeline(context.TODO(), &uid)
	if err != nil {
		return Tweets{}, err
	}

	tweets := Tweets{}
	for i, t := range pbtweets.Tweets {
		tweets[i] = Tweet{
			UserID: t.UserID,
			Text:   t.Text,
		}
	}

	return tweets, nil
}

// User TO DO
type User struct {
	UserID   string
	Username string
	Password string
}

// Tweet TO DO
type Tweet struct {
	TweetID  string
	UserID   string
	Username string
	Text     string
}

// Tweets TO DO
type Tweets []Tweet
