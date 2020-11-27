package follower

import "github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"

type Follower struct {
	FollowerUserID user.ID
	FolloweeUserID user.ID
}

type Repository interface {
	FindAll() ([]Follower, error)
}
