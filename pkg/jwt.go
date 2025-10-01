package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int
	jwt.RegisteredClaims
}

func NewJWTClaims(u int) *Claims {
	return &Claims{
		UserId: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func (c *Claims) GenerateToken() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("jwt secret not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}

func (c *Claims) VerifyToken(tokenString string) error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("jwt secret not found")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return jwt.ErrTokenExpired
	}

	iss, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return err
	}

	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
