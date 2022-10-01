package auth

import (
	"GoVoteApi/config"
	"GoVoteApi/models"
	"GoVoteApi/service"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type auth struct {
	cfg config.Secrets
}

type JwtClaims struct {
	// ID is mongodb ID
	ID       string
	Username string
	Role     models.UserRole
	jwt.RegisteredClaims
}

// GenerateToken implements service.AuthService
func (a *auth) GenerateToken(id, username string, role models.UserRole) (string, error) {
	expTime := time.Now().Add(time.Duration(a.cfg.ExpTime) * time.Minute)
	claims := JwtClaims{
		ID:       id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jc := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return jc.SignedString([]byte(a.cfg.JwtSecret))
}

// VerifyToken implements service.AuthService
// this function check that user has the given role from token string
// and return token claims
func (a *auth) ClaimsFromToken(tokenStr string) (any, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token")
		}

		return []byte(a.cfg.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	_, ok := token.Claims.(*JwtClaims)

	if !ok {
		return nil, fmt.Errorf("cant get jwt claims")
	}

	return token.Claims, nil
}
func New(cfg config.Secrets) service.AuthService {
	return &auth{
		cfg: cfg,
	}
}
