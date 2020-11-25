package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/databaseaccess/internal/infrastructure/dbaccess"
	pb "github.com/martinmhan/tweet-app-api/cmd/databaseaccess/proto"
)

func main() {
	godotenv.Load()

	port := os.Getenv("DA_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if port == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Database Access server failed to listen: ", err)
	}

	m := &dbaccess.MongoDBAccesser{
		DBHost: dbHost,
		DBPort: dbPort,
		DBName: dbName,
	}
	err = m.Connect()
	if err != nil {
		log.Fatal("Failed to connect MongoDB client")
	}

	defer m.Disconnect()

	g := grpc.NewServer()
	s := &application.DatabaseAccessServer{DBAccesser: m}
	pb.RegisterDatabaseAccessServer(g, s)
	err = g.Serve(lis)
	if err != nil {
		log.Fatal("Failed to start Database Access server: ", err)
	}
}
