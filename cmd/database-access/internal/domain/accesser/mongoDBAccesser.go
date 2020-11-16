package accesser

import (
	"context"
	"log"
	"time"

	"github.com/martinmhan/crud-api-golang-grpc/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBAccesser is the MongoDB implementation of the DBAccesser interface
type MongoDBAccesser struct {
	connection *mongo.Client
	DBName     string
	DBHost     string
	DBPort     string
}

// Connect establishes a client connection to MongoDB and sets it as m.connection
func (m *MongoDBAccesser) Connect() error {
	if m.DBHost == "" || m.DBPort == "" || m.DBName == "" {
		log.Fatal("Missing DB connection config")
	}

	connectionURI := "mongodb://" + m.DBHost + ":" + m.DBPort + "/"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal("Error connection to MongoDB: ", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	m.connection = client

	return nil
}

// InsertUser adds a user to mongodb collection
func (m *MongoDBAccesser) InsertUser(f UserFields) User {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	insertResult, err := m.connection.Database(m.DBName).Collection("items").InsertOne(context.TODO(), f)
	if err != nil {
		log.Fatal("Error inserting item: ", err)
	}

	return User{
		ID:       insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Username: f.Username,
	}
}

// InsertTweet adds a tweet to mongodb tweet collection
func (m *MongoDBAccesser) InsertTweet(f TweetFields) Tweet {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	insertResult, err := m.connection.Database(m.DBName).Collection("tweets").InsertOne(context.TODO(), f)
	if err != nil {
		log.Fatal("Error inserting item: ", err)
	}

	return Tweet{
		ID:   insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Text: f.Text,
	}
}

// GetUser gets a user from the mongodb collection given an ID
func (m *MongoDBAccesser) GetUser(UserID string) User {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	record := bson.M{}

	_id, err := primitive.ObjectIDFromHex(UserID)
	utils.FailOnError(err, "Invalid id")

	f := bson.M{"_id": _id}
	result := m.connection.Database(m.DBName).Collection("users").FindOne(context.TODO(), f)
	err = result.Decode(&record)
	if err != nil {
		log.Printf("Error decoding document: %s", err)
		return User{}
	}

	return User{
		ID:       record["_id"].(primitive.ObjectID).Hex(),
		Username: record["username"].(string),
	}
}

// GetTweets returns all items in the mongodb items collection
func (m *MongoDBAccesser) GetTweets(UserID string) []Tweet {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	f := bson.M{}
	cursor, err := m.connection.Database(m.DBName).Collection("tweets").Find(context.TODO(), f)
	utils.FailOnError(err, "Error getting tweets")

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	utils.FailOnError(err, "Error reading tweets")

	tweets := []Tweet{}
	for _, r := range records {
		tweets = append(tweets, Tweet{
			ID:   r["_id"].(primitive.ObjectID).Hex(),
			Text: r["text"].(string),
		})
	}

	return tweets
}
