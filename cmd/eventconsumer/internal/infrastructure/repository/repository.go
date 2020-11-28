package repository

import (
	"context"

	dbaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/user"
	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

// UserRepository implements the user repository
type UserRepository struct {
	dbaccesspb.DatabaseAccessClient
	readviewpb.ReadViewClient
}

// Save inserts a user into the database, then updates the Read View service
func (ur *UserRepository) Save(conf user.Config) (user.User, error) {
	insertID, err := ur.DatabaseAccessClient.SaveUser(
		context.TODO(),
		&dbaccesspb.UserConfig{Username: conf.Username, Password: conf.Password},
	)
	if err != nil {
		return user.User{}, err
	}

	_, err = ur.ReadViewClient.AddUser(
		context.TODO(),
		&readviewpb.User{
			ID:       insertID.InsertID,
			Username: conf.Username,
			Password: conf.Password,
		},
	)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       insertID.InsertID,
		Username: conf.Username,
		Password: conf.Password,
	}, nil
}

// FollowRepository implements the follower repository
type FollowRepository struct {
	dbaccesspb.DatabaseAccessClient
	readviewpb.ReadViewClient
}

// Save adds a new follow (i.e., follower/followee relationship between the two provided user ids)
// to the database, then updates the Read View service
func (fr *FollowRepository) Save(f follow.Config) error {
	follower, err := fr.ReadViewClient.GetUserByUserID(context.TODO(), &readviewpb.UserID{UserID: f.FollowerUserID})
	if err != nil {
		return err
	}

	followee, err := fr.ReadViewClient.GetUserByUserID(context.TODO(), &readviewpb.UserID{UserID: f.FolloweeUserID})
	if err != nil {
		return err
	}

	_, err = fr.DatabaseAccessClient.SaveFollow(
		context.TODO(),
		&dbaccesspb.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: follower.Username,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: followee.Username,
		},
	)
	if err != nil {
		return err
	}

	_, err = fr.ReadViewClient.AddFollow(
		context.TODO(),
		&readviewpb.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: follower.Username,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: followee.Username,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// TweetRepository implements the tweet repository
type TweetRepository struct {
	dbaccesspb.DatabaseAccessClient
	readviewpb.ReadViewClient
}

// Save inserts a tweet into the database, then updates the Read View service
func (tr *TweetRepository) Save(conf tweet.Config) (tweet.Tweet, error) {
	user, err := tr.ReadViewClient.GetUserByUserID(context.TODO(), &readviewpb.UserID{UserID: conf.UserID})
	if err != nil {
		return tweet.Tweet{}, err
	}

	insertID, err := tr.DatabaseAccessClient.SaveTweet(
		context.TODO(),
		&dbaccesspb.TweetConfig{
			UserID:   conf.UserID,
			Username: user.Username,
			Text:     conf.Text,
		},
	)
	if err != nil {
		return tweet.Tweet{}, err
	}

	_, err = tr.ReadViewClient.AddTweet(
		context.TODO(),
		&readviewpb.Tweet{
			ID:       insertID.InsertID,
			UserID:   conf.UserID,
			Username: user.Username,
			Text:     conf.Text,
		},
	)
	if err != nil {
		return tweet.Tweet{}, err
	}

	return tweet.Tweet{
		ID:       insertID.InsertID,
		UserID:   conf.UserID,
		Username: user.Username,
		Text:     conf.Text,
	}, nil
}
