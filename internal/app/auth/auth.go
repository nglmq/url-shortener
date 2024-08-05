package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// TokenExp expiration time for JWT token
const TokenExp = time.Hour * 3

// SecretKey secret key for signing JWT tokens
const SecretKey = "supersecretkey"

// Claims structure for JWT cookie
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildJWTString builds a JWT string with the given claims
func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID decode given JWT token and returns the UserID
func GetUserID(tokenString string) string {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return ""
	}

	if !token.Valid {
		fmt.Println("token is not valid")
		return ""
	}

	return claims.UserID
}
