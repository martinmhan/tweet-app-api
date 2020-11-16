package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/martinmhan/tweet-app-api/cmd/api-gateway/internal/domain/authorization"
	"github.com/martinmhan/tweet-app-api/cmd/api-gateway/internal/domain/events"
	"github.com/martinmhan/tweet-app-api/cmd/api-gateway/internal/domain/user"
	pb "github.com/martinmhan/tweet-app-api/cmd/api-gateway/proto"
	"google.golang.org/grpc/metadata"
)

// APIGatewayServer is a struct type containing the fields and methods used by the API Gateway service
type APIGatewayServer struct {
	pb.UnimplementedAPIGatewayServer
	JWTKey             string
	EventsProducerHost string
	EventsProducerPort string
	ReadViewHost       string
	ReadViewPort       string
}

// LoginUser provides a JWT given a valid username/password
func (s *APIGatewayServer) LoginUser(ctx context.Context, in *pb.User) (*pb.JWT, error) {
	valid := user.ValidatePassword(in.Username, in.Password)
	if !valid {
		fmt.Println("Invalid username/password")
		return nil, errors.New("Invalid username or password")
	}

	token, err := authorization.CreateJWT(in.Username, s.JWTKey)
	if err != nil {
		fmt.Println("Failed to create JWT", err)
		return &pb.JWT{}, err
	}

	return &pb.JWT{JWT: token}, nil
}

// CreateUser calls the event producer to create a new user (if the username is available), then responds to the initial gRPC
func (s *APIGatewayServer) CreateUser(ctx context.Context, in *pb.User) (*pb.SimpleResponse, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := authorization.ValidateJWT(tokenString, s.JWTKey)
	if err != nil {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, err
	}

	valid := user.ValidateNewUsername(in.Username)
	if !valid {
		return &pb.SimpleResponse{Message: "Username already exists"}, errors.New("Failed to create new user")
	}

	events.ProduceUserCreation(s.EventsProducerHost, s.EventsProducerPort, in.Username, in.Password)

	return &pb.SimpleResponse{
		Message: "User Creation accepted",
	}, nil
}

// CreateTweet calls the event producer to create a new tweet, then responds to the initial gRPC
func (s *APIGatewayServer) CreateTweet(ctx context.Context, in *pb.Tweet) (*pb.SimpleResponse, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := authorization.ValidateJWT(tokenString, s.JWTKey)
	if err != nil {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, err
	}

	claims := token.Claims.(*authorization.CustomClaims)
	uid := claims.UserID

	events.ProduceTweetCreation(s.EventsProducerHost, s.EventsProducerPort, in.Text, uid)

	return &pb.SimpleResponse{
		Message: "Tweet Creation accepted",
	}, nil
}

// GetTweets returns the tweets created by a given UserID
func (s *APIGatewayServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := authorization.ValidateJWT(tokenString, s.JWTKey)
	if err != nil {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, err
	}

	claims := token.Claims.(*authorization.CustomClaims)
	uid := claims.UserID

	// TO DO - call read view, return values
}

// GetTimeline returns the timeline (i.e., following users' tweets) of a given UserID
func (s *APIGatewayServer) GetTimeline(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := authorization.ValidateJWT(tokenString, s.JWTKey)
	if err != nil {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, err
	}

	claims := token.Claims.(*authorization.CustomClaims)
	uid := claims.UserID

	// TO DO - call read view, return values
}
