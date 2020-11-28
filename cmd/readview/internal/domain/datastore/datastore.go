package datastore

import (
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"
)

// Datastore is the data store interface
type Datastore interface {
	Initialize() error
	AddUser(user.User) error
	AddFollow(follow.Follow) error
	AddTweet(tweet.Tweet) error
	GetUserByUserID(user.ID) (user.User, error)
	GetUserByUsername(username string) (user.User, error)
	GetTweets(user.ID) ([]tweet.Tweet, error)
	GetTimeline(user.ID) ([]tweet.Tweet, error)
	GetFollowers(user.ID) ([]follow.Follow, error)
	GetFollowees(user.ID) ([]follow.Follow, error)
}
