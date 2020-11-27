package user

// Config contains the fields necessary to create a user
type Config struct {
	Username string
	Password string
}

// Create TO DO
func Create(c Config) error {
	// TO DO
	return nil
}

// User represents an existing user
type User struct {
	ID       string
	Username string
	Password string
}

// Repository TO DO
type Repository interface {
	FindById(UserID string) (User, error)
	FindByUsername(Username string) (User, error)
}
