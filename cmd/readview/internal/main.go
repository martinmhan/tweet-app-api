package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/infrastructure/datastore"
	pb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

func main() {
	godotenv.Load()

	port := os.Getenv("RV_PORT")
	dahost := os.Getenv("DA_HOST")
	daport := os.Getenv("DA_PORT")
	if port == "" || dahost == "" || daport == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	ds := datastore.Datastore{DatabaseAccessHost: dahost, DatabaseAccessPort: daport}
	err := ds.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize data store: ", err)
	}

	g := grpc.NewServer()
	s := &application.ReadViewServer{
		Datastore: ds,
	}
	pb.RegisterReadViewServer(g, s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to start Read View server: ", err)
	}

	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Read View server: ", err)
	}
}
