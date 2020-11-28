package follow

import "github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"

// A Follow represents a unique follower/followee relationship between two users
type Follow struct {
	FollowerUserID   user.ID
	FollowerUsername string
	FolloweeUserID   user.ID
	FolloweeUsername string
}

// Repository is the Follow Repository interface
type Repository interface {
	FindAll() ([]Follow, error)
}
