package tweet

// Config contains the fields necessary to create a tweet
type Config struct {
	UserID   string
	Username string
	Text     string
}

// Tweet represents an existing tweet
type Tweet struct {
	ID       string
	UserID   string
	Username string
	Text     string
}

// Repository TO DO
type Repository interface {
	FindByUserID(UserID string) ([]Tweet, error)
	FindTimelineByUserID(UserID string) ([]Tweet, error)
}
