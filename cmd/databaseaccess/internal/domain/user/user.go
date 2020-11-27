package user

// User TO DO
type User struct {
	ID       string
	Username string
	Password string
}

// Config TO DO
type Config struct {
	Username string
	Password string
}

// Repository TO DO
type Repository interface {
	Save(Config) (insertID string, err error)
	FindByID(userID string) (User, error)
	FindAll() ([]User, error)
}
