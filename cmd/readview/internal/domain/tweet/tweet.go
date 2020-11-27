package tweet

import "github.com/martinmhan/tweet-app-api/cmd/readview/internal/domain/user"

type Tweet struct {
	ID       string
	UserID   user.ID
	Username string
	Text     string
}

type Repository interface {
	FindAll() ([]Tweet, error)
}
