package authorization

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CustomClaims is a custom claims type containing a Username and UserID
type CustomClaims struct {
	Username string
	UserID   string
	jwt.StandardClaims
}

// CreateJWT creates a JSON web token with username and expiration properties given a username and jwtKey
func CreateJWT(username string, jwtKey string) (string, error) {
	claims := CustomClaims{
		username,
		"PLACEHOLDER",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JSON web token
func ValidateJWT(tokenString string, jwtKey string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	_, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return token, nil
}
