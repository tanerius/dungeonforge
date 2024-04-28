package tokenmanager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: change this key! Dont be an idiot to forget
var jwtKey = []byte("this_key_will_change_in_time")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateTokens generates both access and refresh JWT tokens for a given username
func GenerateTokens(username string) (accessToken string, accessTokenExpiresIn int64, refreshToken string, refreshTokenExpiresIn int64, err error) {
	accessTokenExpiresIn = 15 * 60       // 15 minutes
	refreshTokenExpiresIn = 24 * 60 * 60 // 24 hours

	// Access Token
	accessTokenClaims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessTokenExpiresIn) * time.Second)),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err = at.SignedString(jwtKey)
	if err != nil {
		return "", 0, "", 0, err
	}

	// Refresh Token
	refreshTokenClaims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(refreshTokenExpiresIn) * time.Second)),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshToken, err = rt.SignedString(jwtKey)
	if err != nil {
		return "", 0, "", 0, err
	}

	return accessToken, accessTokenExpiresIn, refreshToken, refreshTokenExpiresIn, nil
}

// ValidateToken checks the validity of the token and returns the username if it's valid
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
