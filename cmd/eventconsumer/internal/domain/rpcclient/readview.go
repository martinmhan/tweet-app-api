package rpcclient

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// ReadView implements the methods of the ReadView gRPC client used by this service
// and provides an abstraction layer for establishing connections and transforming protobuf data types
type ReadView struct {
	Host string
	Port string
	conn *grpc.ClientConn
}

// Connect establishes a gRPC client connection
func (rv *ReadView) Connect() error {
	target := rv.Host + ":" + rv.Port
	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.New("Could not connect to Read View server")
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

// AddUser adds a user to the Read View service
func (rv *ReadView) AddUser(u User) error {
	if rv.conn == nil {
		return errors.New("ReadView not connected")
	}

	c := readviewpb.NewReadViewClient(rv.conn)

	user := readviewpb.User{
		UserID:   u.ID,
		Username: u.Username,
		Password: u.Password,
	}

	c.AddUser(context.TODO(), &user)

	return nil
}

// AddTweet TO DO
func (rv *ReadView) AddTweet(t Tweet) error {
	if rv.conn == nil {
		return errors.New("ReadView not connected")
	}

	c := readviewpb.NewReadViewClient(rv.conn)

	tweet := readviewpb.Tweet{
		TweetID:  t.TweetID,
		UserID:   t.UserID,
		Username: t.Username,
		Text:     t.Text,
	}
	c.AddTweet(context.TODO(), &tweet)

	return nil
}

// User TO DO
type User struct {
	ID       string
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

// Follower TO DO
type Follower struct {
	FollowerUserID string
	FolloweeUserID string
}
