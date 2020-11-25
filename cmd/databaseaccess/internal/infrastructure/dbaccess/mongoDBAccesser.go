package dbaccess

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBAccesser is the MongoDB implementation of the DBAccesser interface
type MongoDBAccesser struct {
	client *mongo.Client
	DBHost string
	DBPort string
	DBName string
}

// Connect establishes a client connection to MongoDB
func (m *MongoDBAccesser) Connect() error {
	if m.DBHost == "" || m.DBPort == "" || m.DBName == "" {
		return errors.New("Missing DB connection config")
	}

	connectionURI := "mongodb://" + m.DBHost + ":" + m.DBPort + "/"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal("Error connection to MongoDB: ", err)
		return err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}

	m.client = client

	return nil
}

// Disconnect TO DO
func (m *MongoDBAccesser) Disconnect() error {
	if m.client == nil {
		return errors.New("MongoDBAccesser not connected")
	}

	err := m.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

// InsertUser adds a user to the MongoDB database
func (m *MongoDBAccesser) InsertUser(c UserConfig) (InsertID, error) {
	if m.client == nil {
		return "", errors.New("DBAccesser is not connected to a database")
	}

	insert := bson.M{"username": c.Username, "password": c.Password}
	res, err := m.client.Database(m.DBName).Collection("users").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return InsertID(res.InsertedID.(primitive.ObjectID).Hex()), nil
}

// InsertFollower adds a follower to the MongoDB database
func (m *MongoDBAccesser) InsertFollower(f Follower) (InsertID, error) {
	if m.client == nil {
		return "", errors.New("DBAccesser is not connected to a database")
	}

	insert := bson.M{"followerUserID": f.FollowerUserID, "followeeUserID": f.FolloweeUserID}
	res, err := m.client.Database(m.DBName).Collection("followers").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return InsertID(res.InsertedID.(primitive.ObjectID).Hex()), nil
}

// InsertTweet adds a tweet to the MongoDB database
func (m *MongoDBAccesser) InsertTweet(c TweetConfig) (InsertID, error) {
	if m.client == nil {
		return "", errors.New("DBAccesser is not connected to a database")
	}

	insert := bson.M{"userID": c.UserID, "username": c.Username, "text": c.Text}
	res, err := m.client.Database(m.DBName).Collection("tweets").InsertOne(context.TODO(), insert)
	if err != nil {
		return "", err
	}

	return InsertID(res.InsertedID.(primitive.ObjectID).Hex()), nil
}

// GetUser returns a user given a UserID
func (m *MongoDBAccesser) GetUser(UserID string) (User, error) {
	if m.client == nil {
		return User{}, errors.New("DBAccesser is not connected to a database")
	}

	_id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return User{}, err
	}

	record := bson.M{}
	f := bson.M{"_id": _id}
	res := m.client.Database(m.DBName).Collection("users").FindOne(context.TODO(), f)
	err = res.Decode(&record)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:       record["_id"].(primitive.ObjectID).Hex(),
		Username: record["username"].(string),
		Password: record["password"].(string),
	}, nil
}

// GetFollowers gets the followers of a given UserID
func (m *MongoDBAccesser) GetFollowers(UserID string) ([]Follower, error) {
	if m.client == nil {
		return []Follower{}, errors.New("DBAccesser is not connected to a database")
	}

	f := bson.M{"followeeUserID": UserID}
	cursor, err := m.client.Database(m.DBName).Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []Follower{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []Follower{}, err
	}

	followers := []Follower{}
	for _, r := range records {
		followers = append(followers, Follower{
			FollowerUserID: r["followerUserID"].(string),
			FolloweeUserID: r["followeeUserID"].(string),
		})
	}

	return followers, nil
}

// GetTweets returns all tweets by a given UserID
func (m *MongoDBAccesser) GetTweets(UserID string) ([]Tweet, error) {
	if m.client == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	f := bson.M{"userID": UserID}
	cursor, err := m.client.Database(m.DBName).Collection("tweets").Find(context.TODO(), f)
	if err != nil {
		return []Tweet{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []Tweet{}, err
	}

	tweets := []Tweet{}
	for _, r := range records {
		tweets = append(tweets, Tweet{
			ID:   r["_id"].(primitive.ObjectID).Hex(),
			Text: r["text"].(string),
		})
	}

	return tweets, nil
}

// GetAllUsers gets all users in the MongoDB database (only used by the Read View service on cold starts)
func (m *MongoDBAccesser) GetAllUsers() ([]User, error) {
	if m.client == nil {
		return []User{}, errors.New("DBAccesser is not connected to a database")
	}

	f := bson.M{}
	cursor, err := m.client.Database(m.DBName).Collection("users").Find(context.TODO(), f)
	if err != nil {
		return []User{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []User{}, err
	}

	users := []User{}
	for _, r := range records {
		users = append(users, User{
			ID:       r["_id"].(primitive.ObjectID).Hex(),
			Username: r["username"].(string),
			Password: r["password"].(string),
		})
	}

	return users, nil
}

// GetAllFollowers gets all followers in the MongoDB database (only used by the Read View service on cold starts)
func (m *MongoDBAccesser) GetAllFollowers() ([]Follower, error) {
	if m.client == nil {
		return []Follower{}, errors.New("DBAccesser is not connected to a database")
	}

	f := bson.M{}
	cursor, err := m.client.Database(m.DBName).Collection("followers").Find(context.TODO(), f)
	if err != nil {
		return []Follower{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []Follower{}, err
	}

	followers := []Follower{}
	for _, r := range records {
		followers = append(followers, Follower{
			FollowerUserID: r["followerUserID"].(string),
			FolloweeUserID: r["followeeUserID"].(string),
		})
	}

	return followers, nil
}

// GetAllTweets TO DO
func (m *MongoDBAccesser) GetAllTweets() ([]Tweet, error) {
	if m.client == nil {
		return []Tweet{}, errors.New("DBAccesser is not connected to a database")
	}

	f := bson.M{}
	cursor, err := m.client.Database(m.DBName).Collection("tweets").Find(context.TODO(), f)
	if err != nil {
		return []Tweet{}, err
	}

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		return []Tweet{}, err
	}

	tweets := []Tweet{}
	for _, r := range records {
		tweets = append(tweets, Tweet{
			ID:       r["_id"].(primitive.ObjectID).Hex(),
			UserID:   r["userID"].(string),
			Username: r["username"].(string),
			Text:     r["text"].(string),
		})
	}

	return tweets, nil
}
