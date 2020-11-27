package user

type User struct {
	ID       ID
	Username string
	Password string
}

type ID string

type Config struct {
	Username string
	Password string
}

type Repository interface {
	FindAll() ([]User, error)
}
