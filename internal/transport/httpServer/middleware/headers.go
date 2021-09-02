package middleware

import "net/http"

func JsonRespHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Add("Content-Type", "application/json")

		h.ServeHTTP(w, r)
	})
}
