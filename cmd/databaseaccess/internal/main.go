package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/infrastructure/repository"
	pb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
)

func main() {
	godotenv.Load()

	port := os.Getenv("DA_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if port == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Database Access server failed to listen: ", err)
	}

	connectionURI := "mongodb://" + dbHost + ":" + dbPort + "/"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal("Failed to create MongoDB client: ", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal("Failed to connect MongoDB client: ", err)
	}

	defer client.Disconnect(context.TODO())

	db := client.Database(dbName)
	ur := repository.UserRepository{Database: db}
	fr := repository.FollowRepository{Database: db}
	tr := repository.TweetRepository{Database: db}

	g := grpc.NewServer()
	s := &application.DatabaseAccessServer{
		UserRepository:   &ur,
		FollowRepository: &fr,
		TweetRepository:  &tr,
	}
	pb.RegisterDatabaseAccessServer(g, s)

	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Database Access server: ", err)
	}
}
