package security

import (
	"RealWorldWeb/config"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJWT(username, email string) (string, error) {
	//usages from https://golang-jwt.github.io/jwt/usage/create/
	key := []byte(config.GetSecret())
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

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	var err error
	var bytes []byte

	bytes, err = os.ReadFile(config.GetPrivateKeyLocation())
	if err != nil {
		panic(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}

	bytes, err = os.ReadFile(config.GetPublicKeyLocation())
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}

}
func GenerateJWTByRS256(username, email string) (string, error) {
	tokenDuration := 24 * time.Hour
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"user": map[string]string{
				"email":    email,
				"username": username,
			},
			"iat": now.Unix(),
			"exp": now.Add(tokenDuration).Unix(),
		})
	return t.SignedString(privateKey)
}

func GenerateJWTByHS256(username, email string) (string, error) {
	//usages from https://golang-jwt.github.io/jwt/usage/create/
	key := []byte(config.GetSecret())
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
