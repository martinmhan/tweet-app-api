package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/user"
)

// Authorization is an interface containing auth methods
type Authorization interface {
	CreateJWT(username string) (string, error)
	ValidateJWT(tokenString string) (*jwt.Token, error)
	ValidateUsername(username string) (bool, error)
	ValidatePassword(username string, password string) (bool, error)
}

// New returns an Authorization object
func New(jwtKey string, ur user.Repository) Authorization {
	return &auth{jwtKey, ur}
}

type auth struct {
	JWTKey         string
	UserRepository user.Repository
}

// CreateJWT creates a JSON web token with username and expiration properties given a username and jwtKey
func (a *auth) CreateJWT(username string) (string, error) {
	u, err := a.UserRepository.FindByUsername(username)
	if err != nil {
		return "", err
	}

	if u.ID == "" {
		return "", errors.New("User does not exist")
	}

	claims := JWTClaims{
		username,
		u.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, err := token.SignedString([]byte(a.JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JSON web token
func (a *auth) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(a.JWTKey), nil
	})
	if err != nil {
		return nil, err
	}

	_, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return token, nil
}

// ValidateUsername checks if a user already exists with the given username
func (a *auth) ValidateUsername(username string) (bool, error) {
	u, err := a.UserRepository.FindByUsername(username)
	if err != nil {
		return false, err
	}

	if u.ID != "" {
		return false, nil
	}

	return true, nil
}

// ValidatePassword checks if the given password is correct for the given username
func (a *auth) ValidatePassword(username string, password string) (bool, error) {
	u, err := a.UserRepository.FindByUsername(username)
	if err != nil {
		return false, err
	}

	if u.ID == "" {
		return false, nil
	}

	if u.Password != password {
		return false, nil
	}

	return true, nil
}

// JWTClaims contains the fields stored in a JWT
type JWTClaims struct {
	Username string
	UserID   string
	jwt.StandardClaims
}
