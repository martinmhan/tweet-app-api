package repository

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/follower"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/user"
	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// UserRepository implements the user repository
type UserRepository struct {
	readviewpb.ReadViewClient
}

// FindByUserID fetches a user given a userID
func (ur *UserRepository) FindByUserID(userID string) (user.User, error) {
	uid := readviewpb.UserID{UserID: userID}
	u, err := ur.ReadViewClient.GetUserByUserID(context.TODO(), &uid)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// FindByUsername fetches a user given a username
func (ur *UserRepository) FindByUsername(username string) (user.User, error) {
	un := readviewpb.Username{Username: username}
	u, err := ur.ReadViewClient.GetUserByUsername(context.TODO(), &un)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

// TweetRepository implements the tweet repository
type TweetRepository struct {
	readviewpb.ReadViewClient
}

// FindByUserID fetches tweets of a given user
func (tr *TweetRepository) FindByUserID(userID string) ([]tweet.Tweet, error) {
	uid := readviewpb.UserID{UserID: userID}
	pbtweets, err := tr.ReadViewClient.GetTweets(context.TODO(), &uid)
	if err != nil {
		return []tweet.Tweet{}, err
	}

	tweets := []tweet.Tweet{}
	for i, t := range pbtweets.Tweets {
		tweets[i] = tweet.Tweet{UserID: t.UserID, Text: t.Text}
	}

	return tweets, nil
}

// FindTimelineByUserID fetches the tweets of users followed by a given user
func (tr *TweetRepository) FindTimelineByUserID(userID string) ([]tweet.Tweet, error) {
	uid := readviewpb.UserID{UserID: userID}
	pbtweets, err := tr.ReadViewClient.GetTimeline(context.TODO(), &uid)
	if err != nil {
		return []tweet.Tweet{}, err
	}

	tweets := []tweet.Tweet{}
	for i, t := range pbtweets.Tweets {
		tweets[i] = tweet.Tweet{
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		}
	}

	return tweets, nil
}

// FollowerRepository implements the follower repository
type FollowerRepository struct {
	readviewpb.ReadViewClient
}

// FindByUserID fetches the followers of a given user (i.,e., the followee)
func (fr *FollowerRepository) FindByUserID(userID string) ([]follower.Follower, error) {
	uid := readviewpb.UserID{UserID: userID}
	pbfollowers, err := fr.ReadViewClient.GetFollowers(context.TODO(), &uid)
	if err != nil {
		return []follower.Follower{}, err
	}

	followers := []follower.Follower{}
	for i, f := range pbfollowers.Followers {
		followers[i] = follower.Follower{
			FollowerUserID: f.FollowerUserID,
			FolloweeUserID: f.FolloweeUserID,
		}
	}

	return followers, nil
}
