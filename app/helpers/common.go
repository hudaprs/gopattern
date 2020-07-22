package helpers

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// EncodeAuthToken signs authentication token
func EncodeAuthToken(uid uint, name string, email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["userID"] = uid
	claims["name"] = name
	claims["email"] = email
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}
