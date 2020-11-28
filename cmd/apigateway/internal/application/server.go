package application

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/auth"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/follow"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/tweet"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/user"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/infrastructure/eventproducer"
	pb "github.com/martinmhan/tweet-app-api/cmd/apigateway/proto"
)

// APIGatewayServer contains the fields and gRPC method implementations used by the API Gateway service
type APIGatewayServer struct {
	pb.UnimplementedAPIGatewayServer
	UserRepository   user.Repository
	TweetRepository  tweet.Repository
	FollowRepository follow.Repository
	auth.Authorization
	eventproducer.EventProducer
}

// LoginUser provides a JWT given a valid username/password
func (s *APIGatewayServer) LoginUser(ctx context.Context, in *pb.User) (*pb.JWT, error) {
	valid, err := s.ValidatePassword(in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("Invalid username or password")
	}

	token, err := s.CreateJWT(in.Username)
	if err != nil {
		return &pb.JWT{}, err
	}

	return &pb.JWT{JWT: token}, nil
}

// CreateUser validates the new username, then calls the event producer to publish a CreateUser event and responds to the initial gRPC
func (s *APIGatewayServer) CreateUser(ctx context.Context, in *pb.User) (*pb.SimpleResponse, error) {
	valid, err := s.ValidateUsername(in.Username)
	if err != nil {
		return nil, err
	}

	if !valid {
		return &pb.SimpleResponse{Message: "Username already exists"}, errors.New("Failed to create new user")
	}

	c := user.Config{Username: in.Username, Password: in.Password}
	s.ProduceUserCreation(c)

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
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, err
	}

	claims := token.Claims.(*auth.JWTClaims)
	userID := claims.UserID

	c := tweet.Config{UserID: userID, Username: in.Username, Text: in.Text}
	err = s.ProduceTweetCreation(c)
	if err != nil {
		return &pb.SimpleResponse{Message: "Failed to create tweet"}, err
	}

	return &pb.SimpleResponse{Message: "Tweet Creation accepted"}, nil
}

// CreateFollow calls the event producer to make the current user a follower of the given UserID
func (s *APIGatewayServer) CreateFollow(ctx context.Context, in *pb.UserID) (*pb.SimpleResponse, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		return &pb.SimpleResponse{Message: "Invalid JWT"}, err
	}

	claims := token.Claims.(*auth.JWTClaims)
	currentUserID := claims.UserID
	currentUsername := claims.Username

	followeeUserID := in.UserID
	if currentUserID == followeeUserID {
		return &pb.SimpleResponse{Message: "A user cannot follow him/her self"}, errors.New("Failed to follow user")
	}

	followee, err := s.UserRepository.FindByID(followeeUserID)
	if followee.ID == "" {
		return &pb.SimpleResponse{Message: "Invalid UserID"}, errors.New("Failed to follow user : Invalid UserID")
	}

	f := follow.Follow{
		FollowerUserID:   currentUserID,
		FollowerUsername: currentUsername,
		FolloweeUserID:   followeeUserID,
		FolloweeUsername: followee.Username,
	}
	err = s.ProduceFollowCreation(f)
	if err != nil {
		return &pb.SimpleResponse{Message: "Failed to follow user"}, err
	}

	return &pb.SimpleResponse{Message: "Follow accepted"}, nil
}

// GetFollowers returns the followers of a given UserID
func (s *APIGatewayServer) GetFollowers(ctx context.Context, in *pb.UserID) (*pb.Follows, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.Follows{}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		return &pb.Follows{}, err
	}

	claims := token.Claims.(*auth.JWTClaims)
	userID := claims.UserID

	followers, err := s.FollowRepository.FindFollowersByUserID(userID)
	if err != nil {
		return &pb.Follows{}, err
	}

	var pbFollows pb.Follows
	for i, f := range followers {
		pbFollows.Follows[i] = &pb.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: f.FolloweeUsername,
		}
	}

	return &pbFollows, nil
}

// GetFollowees returns the followees of a given UserID (i.e., users that the user follows)
func (s *APIGatewayServer) GetFollowees(ctx context.Context, in *pb.UserID) (*pb.Follows, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.Follows{}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		return &pb.Follows{}, err
	}

	claims := token.Claims.(*auth.JWTClaims)
	userID := claims.UserID

	followees, err := s.FollowRepository.FindFolloweesByUserID(userID)
	if err != nil {
		return &pb.Follows{}, err
	}

	var pbFollows pb.Follows
	for i, f := range followees {
		pbFollows.Follows[i] = &pb.Follow{
			FollowerUserID:   f.FollowerUserID,
			FollowerUsername: f.FollowerUsername,
			FolloweeUserID:   f.FolloweeUserID,
			FolloweeUsername: f.FolloweeUsername,
		}
	}

	return &pbFollows, nil
}

// GetTweets returns the tweets created by a given UserID
func (s *APIGatewayServer) GetTweets(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.Tweets{}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		return &pb.Tweets{}, err
	}

	claims := token.Claims.(*auth.JWTClaims)
	if claims.UserID != in.UserID {
		return &pb.Tweets{}, errors.New("Unauthorized")
	}

	tweets, err := s.TweetRepository.FindByUserID(in.UserID)
	if err != nil {
		return &pb.Tweets{}, err
	}

	var pbTweets pb.Tweets
	for i, t := range tweets {
		pbTweets.Tweets[i] = &pb.Tweet{
			ID:       t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		}
	}

	return &pbTweets, nil
}

// GetTimeline returns the timeline (i.e., tweets of users that this user follows) of a given UserID
func (s *APIGatewayServer) GetTimeline(ctx context.Context, in *pb.UserID) (*pb.Tweets, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.Tweets{}, errors.New("Failed to find metadata headers from context")
	}

	tokenString := headers["authorization"][0]
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		return &pb.Tweets{}, err
	}

	claims := token.Claims.(*auth.JWTClaims)
	if claims.UserID != in.UserID {
		return &pb.Tweets{}, errors.New("Unauthorized")
	}

	tweets, err := s.TweetRepository.FindByUserID(in.UserID)
	if err != nil {
		return &pb.Tweets{}, err
	}

	var pbTweets pb.Tweets
	for i, t := range tweets {
		pbTweets.Tweets[i] = &pb.Tweet{
			ID:       t.ID,
			UserID:   t.UserID,
			Username: t.Username,
			Text:     t.Text,
		}
	}

	return &pbTweets, nil
}
