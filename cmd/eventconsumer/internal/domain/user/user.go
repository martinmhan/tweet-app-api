package user

// User represents an existing user
type User struct {
	ID       string
	Username string
	Password string
}

// Config contains the fields necessary to create a user
type Config struct {
	Username string
	Password string
}

// Repository is the user repository interface
type Repository interface {
	Save(Config) (User, error)
}
