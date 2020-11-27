package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/follower"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/user"
)

// UserRepository TO DO
type UserRepository struct {
	Database *mongo.Database // replace Client? TO DO
}

// Save TO DO
func (ur *UserRepository) Save(conf user.Config) (insertID string, err error) {
	insert := bson.M{"username": conf.Username, "password": conf.Password}
	res, err := ur.Database.Collection("users").InsertOne(context.TODO(), insert)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindByID TO DO
func (ur *UserRepository) FindByID(userID string) (user.User, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user.User{}, err
	}

	record := bson.M{}
	f := bson.M{"_id": _id}
	res := ur.Database.Collection("users").FindOne(context.TODO(), f)
	err = res.Decode(&record)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       record["_id"].(primitive.ObjectID).Hex(),
		Username: record["username"].(string),
		Password: record["password"].(string),
	}, nil
}

// FindAll TO DO
func (ur *UserRepository) FindAll() ([]user.User, error) {
	f := bson.M{}
	cursor, err := ur.Database.Collection("users").Find(context.TODO(), f)
	if err != nil {
		return []user.User{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []user.User{}, err
	}

	users := []user.User{}
	for _, r := range records {
		users = append(users, user.User{
			ID:       r["_id"].(primitive.ObjectID).Hex(),
			Username: r["username"].(string),
			Password: r["password"].(string),
		})
	}

	return users, nil
}

// FollowerRepository TO DO
type FollowerRepository struct {
	Database *mongo.Database
}

// Save TO DO
func (fr *FollowerRepository) Save(f follower.Follower) (insertID string, err error) {
	insert := bson.M{"followerUserID": f.FollowerUserID, "followeeUserID": f.FolloweeUserID}
	res, err := fr.Database.Collection("followers").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindByUserID TO DO
func (fr *FollowerRepository) FindByUserID(userID string) ([]follower.Follower, error) {
	f := bson.M{"followeeUserID": userID}
	cursor, err := fr.Database.Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []follower.Follower{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []follower.Follower{}, err
	}

	followers := []follower.Follower{}
	for _, r := range records {
		followers = append(followers, follower.Follower{
			FollowerUserID: r["followerUserID"].(string),
			FolloweeUserID: r["followeeUserID"].(string),
		})
	}

	return followers, nil
}

// FindAll TO DO
func (fr *FollowerRepository) FindAll() ([]follower.Follower, error) {
	f := bson.M{}
	cursor, err := fr.Database.Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []follower.Follower{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []follower.Follower{}, err
	}

	followers := []follower.Follower{}
	for _, r := range records {
		followers = append(followers, follower.Follower{
			FollowerUserID: r["followerUserID"].(string),
			FolloweeUserID: r["followeeUserID"].(string),
		})
	}

	return followers, nil
}

// TweetRepository TO DO
type TweetRepository struct {
	Database *mongo.Database
}

// Save TO DO
func (tr *TweetRepository) Save(conf tweet.Config) (insertID string, err error) {
	insert := bson.M{"userID": conf.UserID, "username": conf.Username, "text": conf.Text}
	res, err := tr.Database.Collection("tweets").InsertOne(context.TODO(), insert)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindByUserID TO DO
func (tr *TweetRepository) FindByUserID(userID string) ([]tweet.Tweet, error) {
	f := bson.M{"userID": userID}
	cursor, err := tr.Database.Collection("tweets").Find(context.TODO(), f)
	if err != nil {
		return []tweet.Tweet{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []tweet.Tweet{}, err
	}

	tweets := []tweet.Tweet{}
	for _, r := range records {
		tweets = append(tweets, tweet.Tweet{
			ID:       r["_id"].(primitive.ObjectID).Hex(),
			UserID:   r["userID"].(string),
			Username: r["username"].(string),
			Text:     r["text"].(string),
		})
	}

	return tweets, nil
}

// FindAll TO DO
func (tr *TweetRepository) FindAll() ([]tweet.Tweet, error) {
	f := bson.M{}
	cursor, err := tr.Database.Collection("tweets").Find(context.TODO(), f)
	if err != nil {
		return []tweet.Tweet{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []tweet.Tweet{}, err
	}

	tweets := []tweet.Tweet{}
	for _, r := range records {
		tweets = append(tweets, tweet.Tweet{
			ID:       r["_id"].(primitive.ObjectID).Hex(),
			UserID:   r["userID"].(string),
			Username: r["username"].(string),
			Text:     r["text"].(string),
		})
	}

	return tweets, nil
}
