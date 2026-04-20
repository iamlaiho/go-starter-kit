package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

// TokenPair holds a short-lived access token and a longer-lived refresh token.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GenerateTokenPair creates a signed JWT access and refresh token for the given userID.
func GenerateTokenPair(userID string, secret []byte) (TokenPair, error) {
	access, err := newToken(userID, "access", accessTokenTTL, secret)
	if err != nil {
		return TokenPair{}, err
	}

	refresh, err := newToken(userID, "refresh", refreshTokenTTL, secret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func newToken(userID, kind string, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID,
		"kind": kind,
		"iat":  now.Unix(),
		"exp":  now.Add(ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}
