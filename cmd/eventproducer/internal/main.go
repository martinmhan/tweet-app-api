package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/martinmhan/tweet-app-api/cmd/eventproducer/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/eventproducer/internal/domain/event"
	pb "github.com/martinmhan/tweet-app-api/cmd/eventproducer/proto"

	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	port := os.Getenv("EP_PORT")
	mqhost := os.Getenv("MQ_HOST")
	mqport := os.Getenv("MQ_PORT")
	mqname := os.Getenv("MQ_NAME")
	if port == "" || mqport == "" || mqhost == "" || mqname == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Event Producer failed to listen: ", err)
	}

	p := event.Producer{
		MessageQueueHost: mqhost,
		MessageQueuePort: mqport,
		MessageQueueName: mqname,
	}
	err = p.Connect()
	if err != nil {
		log.Fatal("Event Producer failed to connect: ", err)
	}
	defer p.Disconnect()

	s := &application.EventProducerServer{Producer: p}
	g := grpc.NewServer()
	pb.RegisterEventProducerServer(g, s)

	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Events Producer server: ", err)
	}
}
