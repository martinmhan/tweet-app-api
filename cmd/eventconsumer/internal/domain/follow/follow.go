package follow

// A Follow represents a unique follower/followee relationship between two users
type Follow struct {
	FollowerUserID   string
	FollowerUsername string
	FolloweeUserID   string
	FolloweeUsername string
}

// Repository is the Follower repository interface
type Repository interface {
	Save(Follow) error
}
