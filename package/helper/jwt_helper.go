package helper

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("Yuzsahj!!!22383930")

type Claims struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	SocketGroupID string `json:"socket_group_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, username, SocketID string) (string, error) {
	claims := &Claims{
		UserID:        userID,
		Email:         email,
		Username:      username,
		SocketGroupID: SocketID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
