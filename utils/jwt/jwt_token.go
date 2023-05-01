package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stuton/xm-golang-exercise/internal/application/configuration"
)

type JWT struct {
	config configuration.Configuration
}

func New(config configuration.Configuration) JWT {
	return JWT{config: config}
}

func (t JWT) GenerateToken(userID string) (string, error) {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(t.config.TokenHourLifespan).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.config.ApiSecret))
}

func (t JWT) TokenValid(bearerToken string) error {

	if _, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.config.ApiSecret), nil
	}); err != nil {
		return err
	}

	return nil
}
