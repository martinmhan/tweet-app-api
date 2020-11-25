package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/auth"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/rpcclient"
	pb "github.com/martinmhan/tweet-app-api/cmd/apigateway/proto"
)

func main() {
	godotenv.Load()

	jwtKey := os.Getenv("JWT_KEY")
	port := os.Getenv("AG_PORT")
	ephost := os.Getenv("EP_HOST")
	epport := os.Getenv("EP_PORT")
	rvhost := os.Getenv("RV_HOST")
	rvport := os.Getenv("RV_PORT")
	if jwtKey == "" || port == "" || ephost == "" || epport == "" || rvhost == "" || rvport == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	eventproducerClient := rpcclient.EventProducer{
		Host: ephost,
		Port: epport,
	}
	readviewClient := rpcclient.ReadView{
		Host: rvhost,
		Port: rvport,
	}

	err := eventproducerClient.Connect()
	if err != nil {
		log.Fatal("Failed to connect EventProducer client: ", err)
	}
	err = readviewClient.Connect()
	if err != nil {
		log.Fatal("Failed to connect ReadView client: ", err)
	}

	defer eventproducerClient.Disconnect()
	defer readviewClient.Disconnect()

	a := auth.Authorization{JWTKey: jwtKey, ReadView: readviewClient}

	g := grpc.NewServer()
	s := &application.APIGatewayServer{
		EventProducer: eventproducerClient,
		ReadView:      readviewClient,
		Authorization: a,
	}

	pb.RegisterAPIGatewayServer(g, s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("API Gateway failed to listen: ", err)
	}

	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start API Gateway server: ", err)
	}
}
