package follow

// A Follow represents a unique follower/followee relationship between two users
type Follow struct {
	FollowerUserID   string
	FollowerUsername string
	FolloweeUserID   string
	FolloweeUsername string
}

// Repository interface for fetching users' followers
type Repository interface {
	FindFollowersByUserID(userID string) ([]Follow, error)
	FindFolloweesByUserID(userID string) ([]Follow, error)
}
