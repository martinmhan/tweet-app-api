package datastore

import (
	"errors"
	"log"

	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"
)

// Datastore is an in-memory object that stores a copy of all the app's data
type Datastore struct {
	UserRepository   user.Repository
	FollowRepository follow.Repository
	TweetRepository  tweet.Repository

	Users     map[user.ID]user.User
	Followers map[user.ID][]follow.Follow
	Followees map[user.ID][]follow.Follow
	Tweets    map[user.ID][]tweet.Tweet
}

// Initialize populates the in-memory data store by fetching data via the Database Access service (only called when the server starts)
func (ds *Datastore) Initialize() error {
	log.Println("Initializing data store")

	users, err := ds.UserRepository.FindAll()
	if err != nil {
		return err
	}
	follows, err := ds.FollowRepository.FindAll()
	if err != nil {
		return err
	}
	tweets, err := ds.TweetRepository.FindAll()
	if err != nil {
		return err
	}

	ds.Users = map[user.ID]user.User{}
	ds.Followers = map[user.ID][]follow.Follow{}
	ds.Followees = map[user.ID][]follow.Follow{}
	ds.Tweets = map[user.ID][]tweet.Tweet{}

	for _, u := range users {
		ds.Users[u.ID] = u
	}

	for _, f := range follows {
		// add the follow to the followee's list of followers
		followers, ok := ds.Followers[f.FolloweeUserID]
		if !ok {
			followers = []follow.Follow{}
		}
		ds.Followers[f.FolloweeUserID] = append(followers, f)

		// add the follow to the follower's list of followees
		followees, ok := ds.Followees[f.FollowerUserID]
		if !ok {
			followees = []follow.Follow{}
		}
		ds.Followees[f.FollowerUserID] = append(followees, f)
	}

	for _, t := range tweets {
		_, ok := ds.Tweets[t.UserID]
		if !ok {
			ds.Tweets[t.UserID] = []tweet.Tweet{}
		}
		ds.Tweets[t.UserID] = append(ds.Tweets[t.UserID], t)
	}

	log.Println("Data store initialized")

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

	ds.Users[u.ID] = u

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

// AddFollow adds a follow to the datastore (in both the follower's list of followees and followee's list of followers)
func (ds *Datastore) AddFollow(f follow.Follow) error {
	if f.FollowerUserID == "" || f.FollowerUsername == "" || f.FolloweeUserID == "" || f.FolloweeUsername == "" {
		return errors.New("Invalid follow")
	}

	followers, ok := ds.Followers[f.FolloweeUserID]
	if !ok {
		followers = []follow.Follow{}
	}
	ds.Followers[f.FolloweeUserID] = append(followers, f)

	followees, ok := ds.Followees[f.FollowerUserID]
	if !ok {
		followees = []follow.Follow{}
	}
	ds.Followees[f.FollowerUserID] = append(followees, f)

	return nil
}

// GetUserByUserID returns a user given a userID
func (ds *Datastore) GetUserByUserID(userID user.ID) (user.User, error) {
	u, ok := ds.Users[userID]
	if !ok {
		return user.User{}, errors.New("Invalid UserID")
	}

	return user.User{
		ID:       userID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// GetUserByUsername returns a user given a username
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

	return user.User{}, nil
}

// GetFollowers TO DO
func (ds *Datastore) GetFollowers(userID user.ID) ([]follow.Follow, error) {
	followers, ok := ds.Followers[userID]
	if !ok {
		return []follow.Follow{}, nil
	}

	return followers, nil
}

// GetFollowees returns the tweets of the users that the given user follows
func (ds *Datastore) GetFollowees(userID user.ID) ([]follow.Follow, error) {
	followees, ok := ds.Followees[userID]
	if !ok {
		return []follow.Follow{}, nil
	}

	return followees, nil
}

// GetTweets TO DO
func (ds *Datastore) GetTweets(userID user.ID) ([]tweet.Tweet, error) {
	tweets, ok := ds.Tweets[userID]
	if !ok {
		return []tweet.Tweet{}, nil
	}

	return tweets, nil
}

// GetTimeline returns the tweets of the users that the given user follows
func (ds *Datastore) GetTimeline(userID user.ID) ([]tweet.Tweet, error) {
	followees, ok := ds.Followees[userID]
	if !ok {
		return []tweet.Tweet{}, nil
	}

	var timeline []tweet.Tweet
	for _, f := range followees {
		tweets := ds.Tweets[f.FolloweeUserID]
		timeline = append(timeline, tweets...)
	}

	return timeline, nil
}
