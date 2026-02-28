package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID      uuid.UUID   `json:"user_id"`
	Email       string      `json:"email"`
	Roles       []string    `json:"roles"`
	Permissions []string    `json:"permissions"`
}

type JWTUtil struct {
	secret        string
	expirationTime time.Duration
}

func NewJWTUtil(secret string, expirationTime time.Duration) *JWTUtil {
	return &JWTUtil{
		secret:        secret,
		expirationTime: expirationTime,
	}
}

func (j *JWTUtil) GenerateToken(userID uuid.UUID, email string, roles, permissions []string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expirationTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "ms-ga-identifier",
			ID:        uuid.New().String(),
		},
		UserID:      userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTUtil) GetJTI(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.ID, nil
}

func (j *JWTUtil) GetExpiration(tokenString string) (time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}
