package user

import (
	"GoVoteApi/models"
	"fmt"
	"time"

	// "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
)

type jwtClaims struct {
	id   string
	role models.UserRole
	jwt.RegisteredClaims
}

func genJWT(id, token string, role models.UserRole) (string, error) {
	expTime := time.Now().Add(1 * time.Hour)
	claims := jwtClaims{
		id:   id,
		role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jc := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return jc.SignedString([]byte(token))

}

func ValidateToken(tokenStr, secret string, role models.UserRole) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return false, fmt.Errorf("invalid token")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}

	if token.Valid && token.Claims.(jwtClaims).role == role {
		return true, nil
	}

	return false, fmt.Errorf("invalid token")
}
