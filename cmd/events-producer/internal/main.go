package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/martinmhan/tweet-app-api/cmd/events-producer/internal/application"
	pb "github.com/martinmhan/tweet-app-api/cmd/events-producer/proto"

	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	port := os.Getenv("EP_PORT")
	mqhost := os.Getenv("MQ_HOST")
	mqport := os.Getenv("MQ_PORT")
	if port == "" || mqport == "" || mqhost == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Events Producer failed to listen: ", err)
	}

	g := grpc.NewServer()
	s := &application.EventsProducerServer{
		MessageQueueHost: mqhost,
		MessageQueuePort: mqport,
	}

	pb.RegisterEventsProducerServer(g, s)
	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Events Producer server: ", err)
	}
}
