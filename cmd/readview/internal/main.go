package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	dbaccesspb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/infrastructure/datastore"
	"github.com/martinmhan/tweet-app-api/cmd/readview/internal/infrastructure/repository"
	pb "github.com/martinmhan/tweet-app-api/cmd/readview/proto"
)

func main() {
	godotenv.Load()

	port := os.Getenv("RV_PORT")
	daHost := os.Getenv("DA_HOST")
	daPort := os.Getenv("DA_PORT")
	if port == "" || daHost == "" || daPort == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	target := daHost + ":" + daPort
	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Failed to connect to Database Access service")
	}

	defer conn.Close()

	daClient := dbaccesspb.NewDatabaseAccessClient(conn)
	ur := repository.UserRepository{DatabaseAccessClient: daClient}
	fr := repository.FollowerRepository{DatabaseAccessClient: daClient}
	tr := repository.TweetRepository{DatabaseAccessClient: daClient}

	ds := datastore.Datastore{
		UserRepository:     &ur,
		FollowerRepository: &fr,
		TweetRepository:    &tr,
	}

	err = ds.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize data store: ", err)
	}

	g := grpc.NewServer()
	s := &application.ReadViewServer{Datastore: ds}
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
