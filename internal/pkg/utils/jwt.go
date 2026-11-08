package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"gogogo/internal/model"
)

type AccessClaims struct {
	Role         string `json:"role"`
	TokenVersion uint32 `json:"token_version"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(secret string, ttl time.Duration, user model.User) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	claims := AccessClaims{
		Role:         user.Role,
		TokenVersion: user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(user.ID, 10),
			ID:        NewUUID(),
			Issuer:    "timeseries-platform",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func ParseAccessToken(secret, raw string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(
		raw,
		&AccessClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}
