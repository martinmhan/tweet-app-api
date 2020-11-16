package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/martinmhan/tweet-app-api/cmd/api-gateway/internal/application"
	pb "github.com/martinmhan/tweet-app-api/cmd/api-gateway/proto"
	"github.com/martinmhan/tweet-app-api/util"

	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	jwtKey := os.Getenv("JWT_KEY")
	port := os.Getenv("AG_PORT")
	eph := os.Getenv("EP_HOST")
	epp := os.Getenv("EP_PORT")
	if jwtKey == "" || port == "" || eph == "" || epp == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	util.FailOnError(err, "Failed to listen")

	g := grpc.NewServer()
	s := &application.APIGatewayServer{
		JWTKey:             jwtKey,
		EventsProducerHost: eph,
		EventsProducerPort: epp,
	}

	pb.RegisterAPIGatewayServer(g, s)
	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start API Gateway server: ", err)
	}
}
