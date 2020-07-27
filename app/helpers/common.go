package helpers

import (
	"math/rand"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// EncodeAuthToken signs authentication token
func EncodeAuthToken(uid uint, name string, email string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["UserID"] = uid
	claims["Name"] = name
	claims["Email"] = email
	claims["RoleName"] = role
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

// RandomStringRunes generate a random string
func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
