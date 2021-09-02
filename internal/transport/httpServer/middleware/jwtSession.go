package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/DrIhor/test_task/internal/service/jwt"
)

func JwtSessionCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		check, err := strconv.ParseBool(os.Getenv("CHECK_TOKEN"))
		fmt.Println(check, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if check {
			token := r.Header.Get("token")
			ok, err := jwt.ValidateGoogleExpireJWT(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
