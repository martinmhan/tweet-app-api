package follower

// Follower represents a unique follower/followee relationship between two users
type Follower struct {
	FollowerUserID string
	FolloweeUserID string
}

// Repository interface for fetching users' followers
type Repository interface {
	FindByUserID(userid string) []Follower
}
