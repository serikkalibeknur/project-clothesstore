package middleware

import (
	"net/http"

	"github.com/serikkalibeknur/project-clothesstore/internal/utils"
)

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role")
		if role != "admin" {
			utils.ErrorResponse(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
