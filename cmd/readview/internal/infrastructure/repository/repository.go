package repository

import (
	"context"

	dbaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"
)

// UserRepository implements the User repository
type UserRepository struct {
	dbaccesspb.DatabaseAccessClient
}

// FindAll TO DO
func (ur *UserRepository) FindAll() ([]user.User, error) {
	pbUsers, err := ur.DatabaseAccessClient.GetAllUsers(context.TODO(), &dbaccesspb.GetAllUsersParam{})
	if err != nil {
		return []user.User{}, err
	}

	var users []user.User
	for _, u := range pbUsers.Users {
		users = append(users, user.User{
			ID:       user.ID(u.ID),
			Username: u.Username,
			Password: u.Password,
		})
	}

	return users, nil
}

// FollowRepository implements the Follow repository
type FollowRepository struct {
	dbaccesspb.DatabaseAccessClient
}

// FindAll TO DO
func (ur *FollowRepository) FindAll() ([]follow.Follow, error) {
	pbFollows, err := ur.DatabaseAccessClient.GetAllFollows(context.TODO(), &dbaccesspb.GetAllFollowsParam{})
	if err != nil {
		return []follow.Follow{}, err
	}

	var follows []follow.Follow
	for _, f := range pbFollows.Follows {
		follows = append(follows, follow.Follow{
			FollowerUserID:   user.ID(f.FollowerUserID),
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   user.ID(f.FolloweeUserID),
			FolloweeUsername: f.FolloweeUsername,
		})
	}

	return follows, nil
}

// TweetRepository implements the Tweet repository
type TweetRepository struct {
	dbaccesspb.DatabaseAccessClient
}

// FindAll TO DO
func (ur *TweetRepository) FindAll() ([]tweet.Tweet, error) {
	pbTweets, err := ur.DatabaseAccessClient.GetAllTweets(context.TODO(), &dbaccesspb.GetAllTweetsParam{})
	if err != nil {
		return []tweet.Tweet{}, err
	}

	var tweets []tweet.Tweet
	for _, t := range pbTweets.Tweets {
		tweets = append(tweets, tweet.Tweet{
			ID:       t.ID,
			UserID:   user.ID(t.UserID),
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return tweets, nil
}
