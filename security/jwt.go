package security

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(username, email string) (string, error) {
	//usages from https://golang-jwt.github.io/jwt/usage/create/
	key := []byte("secret")
	tokenDuration := 24 * time.Hour
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user": map[string]string{
				"email":    email,
				"username": username,
			},
			"iat": now.Unix(),
			"exp": now.Add(tokenDuration).Unix(),
		})
	return t.SignedString(key)
}
