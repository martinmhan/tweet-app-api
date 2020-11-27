package follower

// Follower TO DO
type Follower struct {
	FollowerUserID string
	FolloweeUserID string
}

// Repository TO DO
type Repository interface {
	Save(Follower) (insertID string, err error)
	FindByUserID(followeeUserID string) ([]Follower, error)
	FindAll() ([]Follower, error)
}
