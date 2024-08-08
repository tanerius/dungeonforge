package tokenmanager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: change this key! Dont be an idiot to forget
var jwtKey = []byte("this_key_will_change_in_time")

type Tokens struct {
	Access         string
	AccessExires   int64
	Refresh        string
	RefreshExpires int64
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateTokens generates both access and refresh JWT tokens for a given username
func GenerateTokens(username string) (tok *Tokens, err error) {
	tok = &Tokens{}
	tok.AccessExires = 15 * 60        // 15 minutes
	tok.RefreshExpires = 24 * 60 * 60 // 24 hours

	// Access Token
	accessTokenClaims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tok.AccessExires) * time.Second)),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tok.Access, err = at.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshTokenClaims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tok.RefreshExpires) * time.Second)),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tok.Refresh, err = rt.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return tok, nil
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
