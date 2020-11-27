package tweet

// Tweet represents an existing tweet
type Tweet struct {
	ID       string
	UserID   string
	Username string
	Text     string
}

// Config contains the fields necessary to create a tweet
type Config struct {
	UserID   string
	Username string
	Text     string
}

// Repository is the Tweet repository interface
type Repository interface {
	Save(Config) (Tweet, error)
}
