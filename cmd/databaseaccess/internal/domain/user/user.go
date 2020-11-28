package user

// User represent an existing user
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

// Repository is the User Repository interface
type Repository interface {
	Save(Config) (insertID string, err error)
	FindByID(userID string) (User, error)
	FindAll() ([]User, error)
}
