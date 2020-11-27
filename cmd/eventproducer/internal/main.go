package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"

	"github.com/martinmhan/tweet-app-api/cmd/eventproducer/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/eventproducer/internal/infrastructure/eventproducer"
	pb "github.com/martinmhan/tweet-app-api/cmd/eventproducer/proto"
)

func main() {
	godotenv.Load()

	port := os.Getenv("EP_PORT")
	mqHost := os.Getenv("MQ_HOST")
	mqPort := os.Getenv("MQ_PORT")
	mqName := os.Getenv("MQ_NAME")
	if port == "" || mqPort == "" || mqHost == "" || mqName == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Event Producer failed to listen: ", err)
	}

	url := "amqp://guest:guest@" + mqHost + ":" + mqPort + "/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	ep := eventproducer.EventProducer{
		MessageQueueName: mqName,
		Connection:       conn,
	}

	s := &application.EventProducerServer{Producer: &ep}
	g := grpc.NewServer()
	pb.RegisterEventProducerServer(g, s)

	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Events Producer server: ", err)
	}
}
