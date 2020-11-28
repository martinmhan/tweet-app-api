package follow

// A Follow represents a unique follower/followee relationship between two users
type Follow struct {
	FollowerUserID   string
	FollowerUsername string
	FolloweeUserID   string
	FolloweeUsername string
}

// Repository is the FollowRepository interface
type Repository interface {
	Save(Follow) (insertID string, err error)
	FindFollowersByUserID(userID string) ([]Follow, error)
	FindFolloweesByUserID(userID string) ([]Follow, error)
	FindAll() ([]Follow, error)
}
