package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

var jwtKey []byte

func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

// CreateJWTToken creates a JWT token with the given user ID
func CreateJWTToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return token_string, nil
}

// ParseJWTToken parses a JWT token and returns the claims
func ParseJWTToken(token string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	return claims, nil
}
