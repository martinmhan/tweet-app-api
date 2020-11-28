package user

// A User represents an existing user
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

// Repository interface for fetching users
type Repository interface {
	FindByID(UserID string) (User, error)
	FindByUsername(Username string) (User, error)
}
