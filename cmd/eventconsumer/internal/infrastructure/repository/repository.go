package repository

import (
	"context"

	dbaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/follower"
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
	insertID, err := ur.DatabaseAccessClient.InsertUser(
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

// FollowerRepository implements the follower repository
type FollowerRepository struct {
	dbaccesspb.DatabaseAccessClient
	readviewpb.ReadViewClient
}

// Save adds a new follower (i.e., follower/followee relationship between the two provided user ids)
// to the database, then updates the Read View service
func (fr *FollowerRepository) Save(f follower.Follower) error {
	_, err := fr.DatabaseAccessClient.InsertFollower(
		context.TODO(),
		&dbaccesspb.Follower{
			FollowerUserID: f.FollowerUserID,
			FolloweeUserID: f.FolloweeUserID,
		},
	)
	if err != nil {
		return err
	}

	_, err = fr.ReadViewClient.AddFollower(
		context.TODO(),
		&readviewpb.Follower{
			FollowerUserID: f.FollowerUserID,
			FolloweeUserID: f.FolloweeUserID,
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
	insertID, err := tr.DatabaseAccessClient.InsertTweet(
		context.TODO(),
		&dbaccesspb.TweetConfig{
			UserID:   conf.UserID,
			Username: conf.Username,
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
			Username: conf.Username,
			Text:     conf.Text,
		},
	)
	if err != nil {
		return tweet.Tweet{}, err
	}

	return tweet.Tweet{
		ID:       insertID.InsertID,
		UserID:   conf.UserID,
		Username: conf.Username,
		Text:     conf.Text,
	}, nil
}
