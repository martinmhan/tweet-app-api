package follower

// Follower represents a unique follower/followee relationship between two users
type Follower struct {
	FollowerUserID string
	FolloweeUserID string
}

// Repository is the Follower repository interface
type Repository interface {
	Save(Follower) error
}
