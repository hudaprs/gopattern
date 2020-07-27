package middlewares

import (
	"net/http"
	"os"
	"strings"

	"gopattern/app/helpers"
	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type BaseMiddleware struct {
	DB *gorm.DB
}

// SetContentTypeHeader to JSON
func SetContentTypeHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// AuthJwtVerify verify token and add UserID to the request context
func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = map[string]interface{}{"Status": "failed", "Message": "Missing Authorization Token"}

		var header = r.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			helpers.JSON(w, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			resp["Status"] = "failed"
			resp["Message"] = "Invalid token, please login"
			helpers.JSON(w, http.StatusForbidden, resp)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		context.Set(r, "UserID", claims["UserID"])
		context.Set(r, "RoleName", claims["RoleName"])
		next.ServeHTTP(w, r)
	})
}

// OnlyAdmin can access
func OnlyHighAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleName := context.Get(r, "RoleName")

		if roleName != "High Admin"  {
			helpers.JSON(w, http.StatusUnauthorized, map[string]interface{}{"Status": "Unauthorized", "Message": "You can't access this page"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
