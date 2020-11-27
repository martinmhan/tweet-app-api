package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"

	dbaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/infrastructure/repository"
	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

func main() {
	godotenv.Load()

	mqHost := os.Getenv("MQ_HOST")
	mqPort := os.Getenv("MQ_PORT")
	mqName := os.Getenv("MQ_NAME")
	daHost := os.Getenv("DA_HOST")
	daPort := os.Getenv("DA_PORT")
	rvHost := os.Getenv("RV_HOST")
	rvPort := os.Getenv("RV_PORT")
	if mqPort == "" || mqHost == "" || mqName == "" || daHost == "" || daPort == "" || rvHost == "" || rvPort == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	url := "amqp://guest:guest@" + mqHost + ":" + mqPort + "/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	daTarget := daHost + ":" + daPort
	daCtx, daCancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer daCancel()

	daConn, err := grpc.DialContext(daCtx, daTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Could not connect to database access server")
	}
	defer daConn.Close()

	rvTarget := rvHost + ":" + rvPort
	rvCtx, rvCancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer rvCancel()

	rvConn, err := grpc.DialContext(rvCtx, rvTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Could not connect to database access server")
	}
	defer rvConn.Close()

	daClient := dbaccesspb.NewDatabaseAccessClient(daConn)
	rvClient := readviewpb.NewReadViewClient(rvConn)

	ur := repository.UserRepository{DatabaseAccessClient: daClient, ReadViewClient: rvClient}
	fr := repository.FollowerRepository{DatabaseAccessClient: daClient, ReadViewClient: rvClient}
	tr := repository.TweetRepository{DatabaseAccessClient: daClient, ReadViewClient: rvClient}

	s := &application.EventConsumerServer{
		Connection:         conn,
		MessageQueueName:   mqName,
		UserRepository:     &ur,
		FollowerRepository: &fr,
		TweetRepository:    &tr,
	}

	s.Listen()
}
