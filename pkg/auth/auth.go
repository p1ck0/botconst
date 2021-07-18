package auth

import (
	"errors"
	"fmt"
	"time"

	jwtv3 "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type TokenManager interface {
	NewJWT(userID string, ttl time.Duration) (string, error)
	Parse(token string) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signingKey")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwtv3.NewWithClaims(jwt.SigningMethodHS256, jwtv3.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwtv3.Parse(accessToken, func(token *jwtv3.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwtv3.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwtv3.MapClaims)
	if !ok {
		return "", fmt.Errorf("error getting user claims from token")
	}

	return claims["sub"].(string), nil
}

func ParseToken(c *fiber.Ctx) string {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	return claims["sub"].(string)
}
