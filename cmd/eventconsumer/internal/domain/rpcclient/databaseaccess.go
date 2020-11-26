package rpcclient

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	databaseaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
)

// DatabaseAccess implements the methods of the DatabaseAccess gRPC client used by this service
// and provides an abstraction layer for establishing connections and transforming protobuf data types
type DatabaseAccess struct {
	Host string
	Port string
	conn *grpc.ClientConn
}

// Connect establishes a gRPC client connection
func (da *DatabaseAccess) Connect() error {
	target := da.Host + ":" + da.Port
	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.New("Could not connect to database access server")
	}

	da.conn = conn

	return nil
}

// Disconnect closes the gRPC client connection
func (da *DatabaseAccess) Disconnect() error {
	if da.conn == nil {
		return errors.New("DatabaseAccess not connected")
	}

	err := da.conn.Close()
	if err != nil {
		return err
	}

	da.conn = nil

	return nil
}

// InsertUser makes an RPC to the Database Access service to insert a user into the database
func (da *DatabaseAccess) InsertUser(u UserConfig) (insertID string, err error) {
	if da.conn == nil {
		return "", errors.New("DatabaseAccess not connected")
	}

	c := databaseaccesspb.NewDatabaseAccessClient(da.conn)

	uc := databaseaccesspb.UserConfig{Username: u.Username, Password: u.Password}
	id, err := c.InsertUser(context.TODO(), &uc)
	if err != nil {
		return "", err
	}

	return id.InsertID, nil
}

// InsertTweet sends an RPC to the Database Access service to insert a tweet into the database
func (da *DatabaseAccess) InsertTweet(t TweetConfig) (insertID string, err error) {
	if da.conn == nil {
		return "", errors.New("DatabaseAccess not connected")
	}

	c := databaseaccesspb.NewDatabaseAccessClient(da.conn)

	tc := databaseaccesspb.TweetConfig{UserID: t.UserID, Username: t.Username, Text: t.Text}
	id, err := c.InsertTweet(context.TODO(), &tc)
	if err != nil {
		return "", err
	}

	return id.InsertID, nil
}

// UserConfig provides the fields necessary to create a user - must match databaseaccesspb.UserConfig
type UserConfig struct {
	Username string
	Password string
}

// TweetConfig provides the fields necessary to create a tweet - must match databaseaccesspb.TweetConfig
type TweetConfig struct {
	UserID   string
	Username string
	Text     string
}
