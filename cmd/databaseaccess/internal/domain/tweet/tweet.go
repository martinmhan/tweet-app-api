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

// Repository is the Tweet Repository interface
type Repository interface {
	Save(Config) (insertID string, err error)
	FindByUserID(userID string) ([]Tweet, error)
	FindAll() ([]Tweet, error)
}
