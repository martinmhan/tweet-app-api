package follow

// A Follow represents a unique follower/followee relationship between two users
type Follow struct {
	FollowerUserID   string
	FollowerUsername string
	FolloweeUserID   string
	FolloweeUsername string
}

// Config contains the fields necessary to create a follow
type Config struct {
	FollowerUserID string
	FolloweeUserID string
}

// Repository is the Follower repository interface
type Repository interface {
	Save(Config) error
}
