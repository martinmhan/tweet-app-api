package repository

import (
	dbaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/follower"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"
)

type UserRepository struct {
	dbaccesspb.DatabaseAccessClient
}

func (ur *UserRepository) FindAll() ([]user.User, error) {
	return []user.User{}, nil
}

type FollowerRepository struct {
	dbaccesspb.DatabaseAccessClient
}

func (ur *FollowerRepository) FindAll() ([]follower.Follower, error) {
	return []follower.Follower{}, nil
}

type TweetRepository struct {
	dbaccesspb.DatabaseAccessClient
}

func (ur *TweetRepository) FindAll() ([]tweet.Tweet, error) {
	return []tweet.Tweet{}, nil
}
