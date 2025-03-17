package credential

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateAccessToken(uid int64, jwtKey []byte) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(uid int64, jwtKey []byte) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(jwtToken []byte, tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtToken, nil
	})
	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Проверяем срок действия токена
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
