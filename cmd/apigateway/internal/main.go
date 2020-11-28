package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/auth"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/infrastructure/eventproducer"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/infrastructure/repository"
	pb "github.com/martinmhan/tweet-app-api/cmd/apigateway/proto"
	eventproducerpb "github.com/martinmhan/tweet-app-api/cmd/eventproducer/proto"
	readviewpb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

func main() {
	godotenv.Load()

	jwtKey := os.Getenv("JWT_KEY")
	port := os.Getenv("AG_PORT")
	epHost := os.Getenv("EP_HOST")
	epPort := os.Getenv("EP_PORT")
	rvHost := os.Getenv("RV_HOST")
	rvPort := os.Getenv("RV_PORT")
	if jwtKey == "" || port == "" || epHost == "" || epPort == "" || rvHost == "" || rvPort == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	rvTarget := rvHost + ":" + rvPort
	rvCtx, rvCancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer rvCancel()

	rvConn, err := grpc.DialContext(rvCtx, rvTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Failed to connect readview gRPC client")
	}
	defer rvConn.Close()

	epTarget := epHost + ":" + epPort
	epCtx, epCancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer epCancel()

	epConn, err := grpc.DialContext(epCtx, epTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Failed to connect eventproducer gRPC client")
	}
	defer epConn.Close()

	rvClient := readviewpb.NewReadViewClient(rvConn)
	epClient := eventproducerpb.NewEventProducerClient(epConn)

	ur := repository.UserRepository{ReadViewClient: rvClient}
	fr := repository.FollowRepository{ReadViewClient: rvClient}
	tr := repository.TweetRepository{ReadViewClient: rvClient}
	auth := auth.Authorization{JWTKey: jwtKey, UserRepository: &ur}
	ep := eventproducer.EventProducer{EventProducerClient: epClient}
	s := &application.APIGatewayServer{
		UserRepository:   &ur,
		FollowRepository: &fr,
		TweetRepository:  &tr,
		Authorization:    auth,
		EventProducer:    ep,
	}

	g := grpc.NewServer()
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
