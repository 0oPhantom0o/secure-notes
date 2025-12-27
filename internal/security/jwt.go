package security

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrMissingToken = errors.New("missing token")
)

type JWTManager struct {
	Secret []byte
	Issuer string
	TTL    time.Duration
}

func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		Secret: []byte(secret),
		Issuer: issuer,
		TTL:    ttl,
	}
}

// Sign creates a JWT for a user ID and returns token + expiry time.
func (m *JWTManager) Sign(userID int64) (tokenString string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(m.TTL)

	claims := jwt.RegisteredClaims{
		Issuer:    m.Issuer,
		Subject:   strconv.FormatInt(userID, 10), // user id in sub
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = t.SignedString(m.Secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expiresAt, nil
}

// Parse verifies a token and extracts userID from "sub".
func (m *JWTManager) Parse(tokenString string) (int64, error) {
	if tokenString == "" {
		return 0, ErrMissingToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		// Enforce HS256
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.Secret, nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	uid, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || uid <= 0 {
		return 0, ErrInvalidToken
	}

	// Optional strict issuer check
	if m.Issuer != "" && claims.Issuer != m.Issuer {
		return 0, ErrInvalidToken
	}

	return uid, nil
}
