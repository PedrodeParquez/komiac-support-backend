package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleUser    Role = "user"
	RoleSupport Role = "support"
)

type Claims struct {
	UID  int  `json:"uid"`
	Role Role `json:"role"`
	jwt.RegisteredClaims
}

func Sign(uid int, role Role, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func Parse(tokenStr, secret string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
