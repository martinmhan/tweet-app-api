package application

import (
	"fmt"
)

type user struct {
	username string
	password string
}

type tweet struct {
	username string
	userID   string
	text     string
}

type followers struct {
	userID    string
	followers []string
}

type datastore struct {
	users     map[string]user
	followers map[string][]user
	tweets    map[string][]tweet
}

type ReadViewServer struct {
	pb.UnimplementedReadViewServer
	DataAccesHost  string
	DataAccessPort string
	datastore      datastore
}

func (s *ReadViewServer) GetAllUsers() {
	fmt.Println("Endpoint hit: GetAllUsers")
	// TO DO
}

func (s *ReadViewServer) GetAllTweets() {
	fmt.Println("Endpoint hit: GetAllTweets")
	// TO DO
}

func (s *ReadViewServer) GetAllUserFollowers() {
	fmt.Println("Endpoint hit: GetAllUserFollowers")
	// TO DO
}

func (s *ReadViewServer) init() {
	// TO DO - call data access server to populate in-memory data store
}
