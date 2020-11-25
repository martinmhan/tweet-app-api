package datastore

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	databaseaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
)

// Datastore TO DO
type Datastore struct {
	DatabaseAccessHost string
	DatabaseAccessPort string
	Users              map[UserID]UserConfig
	Followers          map[UserID][]UserID
	Tweets             map[UserID][]Tweet
}

// Initialize populates the in-memory data store by fetching data via the Database Access service (only called when the server starts)
func (ds *Datastore) Initialize() error {
	conn, err := grpc.Dial(ds.DatabaseAccessHost+":"+ds.DatabaseAccessPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := databaseaccesspb.NewDatabaseAccessClient(conn)

	users, err := client.GetAllUsers(context.TODO(), &databaseaccesspb.GetAllUsersParam{})
	if err != nil {
		return err
	}
	followers, err := client.GetAllFollowers(context.TODO(), &databaseaccesspb.GetAllFollowersParam{})
	if err != nil {
		return err
	}
	tweets, err := client.GetAllTweets(context.TODO(), &databaseaccesspb.GetAllTweetsParam{})
	if err != nil {
		return err
	}

	for _, u := range users.Users {
		ds.Users[UserID(u.UserID)] = UserConfig{
			Username: u.Username,
			Password: u.Password,
		}
	}

	for _, f := range followers.Followers {
		userID := UserID(f.FolloweeUserID)
		followerID := UserID(f.FollowerUserID)

		_, ok := ds.Followers[userID]
		if !ok {
			ds.Followers[userID] = []UserID{}
		}

		ds.Followers[userID] = append(ds.Followers[userID], followerID)
	}

	for _, t := range tweets.Tweets {
		userID := UserID(t.UserID)
		tweet := Tweet{
			TweetID:  t.TweetID,
			UserID:   UserID(t.UserID),
			Username: t.Username,
			Text:     t.Text,
		}

		_, ok := ds.Tweets[userID]
		if !ok {
			ds.Tweets[userID] = []Tweet{}
		}

		ds.Tweets[userID] = append(ds.Tweets[userID], tweet)
	}

	return nil
}

// AddUser adds a user to the datastore
func (ds *Datastore) AddUser(u User) error {
	if u.UserID == "" || u.Username == "" {
		return errors.New("Invalid user")
	}

	_, ok := ds.Users[u.UserID]
	if ok {
		return errors.New("User already exists")
	}

	ds.Users[u.UserID] = UserConfig{Username: u.Username, Password: u.Password}

	return nil
}

// AddTweet adds a tweets to the datastore
func (ds *Datastore) AddTweet(t Tweet) error {
	if t.TweetID == "" || t.UserID == "" || t.Username == "" || t.Text == "" {
		return errors.New("Invalid tweet")
	}

	tweets, ok := ds.Tweets[t.UserID]
	if !ok {
		tweets = []Tweet{}
	}

	ds.Tweets[t.UserID] = append(tweets, t)

	return nil
}

// AddFollower adds a follower to the datastore
func (ds *Datastore) AddFollower(followerUserID UserID, followsUserID UserID) error {
	if followerUserID == "" || followsUserID == "" {
		return errors.New("Invalid userID(s)")
	}

	followers, ok := ds.Followers[followsUserID]
	if !ok {
		followers = []UserID{}
	}

	ds.Followers[followsUserID] = append(followers, followerUserID)

	return nil
}

// GetUserByUserID TO DO
func (ds *Datastore) GetUserByUserID(uid UserID) (User, error) {
	u, ok := ds.Users[uid]
	if !ok {
		return User{}, errors.New("Invalid UserID")
	}

	return User{
		UserID:   uid,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetUserByUsername TO DO
func (ds *Datastore) GetUserByUsername(username string) (User, error) {
	for uid, u := range ds.Users {
		if u.Username == username {
			return User{
				UserID:   uid,
				Username: u.Username,
				Password: u.Password,
			}, nil
		}
	}

	return User{}, nil
}

// GetTweets TO DO
func (ds *Datastore) GetTweets(uid UserID) ([]Tweet, error) {
	tweets, ok := ds.Tweets[uid]
	if !ok {
		return []Tweet{}, nil
	}

	return tweets, nil
}

// GetTimeline TO DO
func (ds *Datastore) GetTimeline(uid UserID) ([]Tweet, error) {
	followers, ok := ds.Followers[uid]
	if !ok {
		return []Tweet{}, nil
	}

	var timeline []Tweet
	for _, fid := range followers {
		tweets := ds.Tweets[fid]
		timeline = append(timeline, tweets...)
	}

	return timeline, nil
}

// User TO DO
type User struct {
	UserID   UserID
	Username string
	Password string
}

// UserID TO DO
type UserID string

// UserConfig TO DO
type UserConfig struct {
	Username string
	Password string
}

// Tweet TO DO
type Tweet struct {
	TweetID  string
	UserID   UserID
	Username string
	Text     string
}
