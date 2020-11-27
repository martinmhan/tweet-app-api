package datastore

import (
	"errors"
	"log"

	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/follower"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"
)

// Datastore is an in-memory object that stores a copy of all the data for this app
type Datastore struct {
	UserRepository     user.Repository
	FollowerRepository follower.Repository
	TweetRepository    tweet.Repository

	Users     map[user.ID]user.Config
	Followers map[user.ID][]user.ID
	Tweets    map[user.ID][]tweet.Tweet
}

// Initialize populates the in-memory data store by fetching data via the Database Access service (only called when the server starts)
func (ds *Datastore) Initialize() error {
	log.Println("Initializing data store")

	users, err := ds.UserRepository.FindAll()
	if err != nil {
		return err
	}
	followers, err := ds.FollowerRepository.FindAll()
	if err != nil {
		return err
	}
	tweets, err := ds.TweetRepository.FindAll()
	if err != nil {
		return err
	}

	for _, u := range users {
		ds.Users[u.ID] = user.Config{
			Username: u.Username,
			Password: u.Password,
		}
	}

	for _, f := range followers {
		userID := f.FolloweeUserID
		followerID := f.FollowerUserID

		_, ok := ds.Followers[userID]
		if !ok {
			ds.Followers[userID] = []user.ID{}
		}

		ds.Followers[userID] = append(ds.Followers[userID], followerID)
	}

	for _, t := range tweets {
		_, ok := ds.Tweets[t.UserID]
		if !ok {
			ds.Tweets[t.UserID] = []tweet.Tweet{}
		}

		ds.Tweets[t.UserID] = append(ds.Tweets[t.UserID], t)
	}

	return nil
}

// AddUser adds a user to the datastore
func (ds *Datastore) AddUser(u user.User) error {
	if u.ID == "" || u.Username == "" {
		return errors.New("Invalid user")
	}

	_, ok := ds.Users[u.ID]
	if ok {
		return errors.New("User already exists")
	}

	ds.Users[u.ID] = user.Config{Username: u.Username, Password: u.Password}

	return nil
}

// AddTweet adds a tweets to the datastore
func (ds *Datastore) AddTweet(t tweet.Tweet) error {
	if t.ID == "" || t.UserID == "" || t.Username == "" || t.Text == "" {
		return errors.New("Invalid tweet")
	}

	tweets, ok := ds.Tweets[t.UserID]
	if !ok {
		tweets = []tweet.Tweet{}
	}

	ds.Tweets[t.UserID] = append(tweets, t)

	return nil
}

// AddFollower adds a follower to the datastore
func (ds *Datastore) AddFollower(f follower.Follower) error {
	if f.FollowerUserID == "" || f.FolloweeUserID == "" {
		return errors.New("Invalid userID(s)")
	}

	followers, ok := ds.Followers[f.FolloweeUserID]
	if !ok {
		followers = []user.ID{}
	}

	ds.Followers[f.FolloweeUserID] = append(followers, f.FollowerUserID)

	return nil
}

// GetUserByUserID TO DO
func (ds *Datastore) GetUserByUserID(uid user.ID) (user.User, error) {
	u, ok := ds.Users[uid]
	if !ok {
		return user.User{}, errors.New("Invalid UserID")
	}

	return user.User{
		ID:       uid,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetUserByUsername TO DO
func (ds *Datastore) GetUserByUsername(username string) (user.User, error) {
	for uid, u := range ds.Users {
		if u.Username == username {
			return user.User{
				ID:       uid,
				Username: u.Username,
				Password: u.Password,
			}, nil
		}
	}

	return user.User{}, errors.New("User does not exist")
}

// GetTweets TO DO
func (ds *Datastore) GetTweets(uid user.ID) ([]tweet.Tweet, error) {
	tweets, ok := ds.Tweets[uid]
	if !ok {
		return []tweet.Tweet{}, nil
	}

	return tweets, nil
}

// GetTimeline TO DO
func (ds *Datastore) GetTimeline(uid user.ID) ([]tweet.Tweet, error) {
	followers, ok := ds.Followers[uid]
	if !ok {
		return []tweet.Tweet{}, nil
	}

	var timeline []tweet.Tweet
	for _, fid := range followers {
		tweets := ds.Tweets[fid]
		timeline = append(timeline, tweets...)
	}

	return timeline, nil
}
