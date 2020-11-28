package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/domain/user"
)

// UserRepository implements the User Repository
type UserRepository struct {
	Database *mongo.Database
}

// Save inserts a user into the database
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

// FollowRepository TO DO
type FollowRepository struct {
	Database *mongo.Database
}

// Save TO DO
func (fr *FollowRepository) Save(f follow.Follow) (insertID string, err error) {
	insert := bson.M{
		"followerUserID":   f.FollowerUserID,
		"followerUsername": f.FollowerUsername,
		"followeeUserID":   f.FolloweeUserID,
		"followeeUsername": f.FolloweeUsername,
	}
	res, err := fr.Database.Collection("followers").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindFollowersByUserID TO DO
func (fr *FollowRepository) FindFollowersByUserID(userID string) ([]follow.Follow, error) {
	f := bson.M{"followeeUserID": userID}
	cursor, err := fr.Database.Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []follow.Follow{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []follow.Follow{}, err
	}

	followers := []follow.Follow{}
	for _, r := range records {
		followers = append(followers, follow.Follow{
			FollowerUserID:   r["followerUserID"].(string),
			FollowerUsername: r["followerUsername"].(string),
			FolloweeUserID:   r["followeeUserID"].(string),
			FolloweeUsername: r["followeeUsername"].(string),
		})
	}

	return followers, nil
}

// FindFolloweesByUserID TO DO
func (fr *FollowRepository) FindFolloweesByUserID(userID string) ([]follow.Follow, error) {
	f := bson.M{"followerUserID": userID}
	cursor, err := fr.Database.Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []follow.Follow{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []follow.Follow{}, err
	}

	followers := []follow.Follow{}
	for _, r := range records {
		followers = append(followers, follow.Follow{
			FollowerUserID:   r["followerUserID"].(string),
			FollowerUsername: r["followerUsername"].(string),
			FolloweeUserID:   r["followeeUserID"].(string),
			FolloweeUsername: r["followeeUsername"].(string),
		})
	}

	return followers, nil
}

// FindAll TO DO
func (fr *FollowRepository) FindAll() ([]follow.Follow, error) {
	f := bson.M{}
	cursor, err := fr.Database.Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []follow.Follow{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []follow.Follow{}, err
	}

	followers := []follow.Follow{}
	for _, r := range records {
		followers = append(followers, follow.Follow{
			FollowerUserID:   r["followerUserID"].(string),
			FollowerUsername: r["followerUsername"].(string),
			FolloweeUserID:   r["followeeUserID"].(string),
			FolloweeUsername: r["followeeUsername"].(string),
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
