package repository

import (
	"context"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/user"
	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// UserRepository implements the user repository
type UserRepository struct {
	readviewpb.ReadViewClient
}

// FindByID fetches a user given a userID
func (ur *UserRepository) FindByID(userID string) (user.User, error) {
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
	for _, t := range pbtweets.Tweets {
		tweets = append(tweets, tweet.Tweet{
			ID:       t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		})
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
	for _, t := range pbtweets.Tweets {
		tweets = append(tweets, tweet.Tweet{
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		})
	}

	return tweets, nil
}

// FollowRepository implements the follower repository
type FollowRepository struct {
	readviewpb.ReadViewClient
}

// FindFollowersByUserID fetches the followers of a given user
func (fr *FollowRepository) FindFollowersByUserID(userID string) ([]follow.Follow, error) {
	uid := readviewpb.UserID{UserID: userID}
	pbFollows, err := fr.ReadViewClient.GetFollowers(context.TODO(), &uid)
	if err != nil {
		return []follow.Follow{}, err
	}

	followers := []follow.Follow{}
	for _, f := range pbFollows.Follows {
		followers = append(followers, follow.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: f.FolloweeUsername,
		})
	}

	return followers, nil
}

// FindFolloweesByUserID fetches the followees of a given user (i.e., other users that the user follows)
func (fr *FollowRepository) FindFolloweesByUserID(userID string) ([]follow.Follow, error) {
	uid := readviewpb.UserID{UserID: userID}
	pbFollows, err := fr.ReadViewClient.GetFollowees(context.TODO(), &uid)
	if err != nil {
		return []follow.Follow{}, err
	}

	followees := []follow.Follow{}
	for _, f := range pbFollows.Follows {
		followees = append(followees, follow.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: f.FolloweeUsername,
		})
	}

	return followees, nil
}
