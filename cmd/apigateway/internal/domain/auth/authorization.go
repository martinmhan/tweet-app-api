package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/martinmhan/tweet-app-api/cmd/apigateway/internal/domain/rpcclient"
)

// Authorization provides methods for creating/validating JWTs, passwords, and new usernames
type Authorization struct {
	JWTKey   string
	ReadView rpcclient.ReadView
}

// CreateJWT creates a JSON web token with username and expiration properties given a username and jwtKey
func (a *Authorization) CreateJWT(username string) (string, error) {
	u, err := a.ReadView.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if u.UserID == "" {
		return "", errors.New("User does not exist")
	}

	claims := JWTClaims{
		username,
		u.UserID,
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
func (a *Authorization) ValidateJWT(tokenString string) (*jwt.Token, error) {
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
func (a *Authorization) ValidateUsername(username string) (bool, error) {
	u, err := a.ReadView.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	if u.UserID != "" {
		return false, nil
	}

	return true, nil
}

// ValidatePassword checks if the given password is correct for the given username
func (a *Authorization) ValidatePassword(username string, password string) (bool, error) {
	u, err := a.ReadView.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	if u.UserID == "" {
		return false, nil
	}

	if u.Password != password {
		return false, nil
	}

	return true, nil
}

// JWTClaims TO DO
type JWTClaims struct {
	Username string
	UserID   string
	jwt.StandardClaims
}
