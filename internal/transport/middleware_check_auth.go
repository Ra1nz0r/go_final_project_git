package transport

import (
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CheckAuth(endpoint http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passFromEnv, exists := os.LookupEnv("TODO_PASSWORD")
		if exists && passFromEnv != "" {
			var passHash string
			cookie, errCook := r.Cookie("token")
			if errCook == nil {
				passHash = cookie.Value
			}
			if errBC := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(passFromEnv)); errBC != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		endpoint(w, r)
	})
}
