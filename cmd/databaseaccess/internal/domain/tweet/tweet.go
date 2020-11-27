package tweet

// Config TO DO
type Config struct {
	UserID   string
	Username string
	Text     string
}

// Tweet TO DO
type Tweet struct {
	ID       string
	UserID   string
	Username string
	Text     string
}

// Repository TO DO
type Repository interface {
	Save(Config) (insertID string, err error)
	FindByUserID(userID string) ([]Tweet, error)
	FindAll() ([]Tweet, error)
}
