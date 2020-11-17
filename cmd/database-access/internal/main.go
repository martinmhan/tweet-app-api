package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/martinmhan/tweet-app-api/cmd/database-access/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/database-access/internal/domain/accesser"
	pb "github.com/martinmhan/tweet-app-api/cmd/database-access/proto"

	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	port := os.Getenv("DBA_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if port == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Database Access server failed to listen: ", err)
	}

	g := grpc.NewServer()
	m := &accesser.MongoDBAccesser{
		DBHost: dbHost,
		DBPort: dbPort,
		DBName: dbName,
	}
	s := &application.DatabaseAccessServer{DBAccesser: m}

	pb.RegisterDatabaseAccessServer(g, s)
	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Database Access server: ", err)
	}
}
