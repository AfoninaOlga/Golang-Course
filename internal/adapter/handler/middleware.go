package handler

import (
	"context"
	"github.com/AfoninaOlga/xkcd/internal/core/port"
	"log"
	"net/http"
)

func Auth(adminRequired bool, authService port.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, err := authService.GetUserByToken(req.Context(), req.Header.Get("Authorization"))
		if err != nil {
			log.Printf("authentification error: %v\n", err)
			http.Error(w, "Authentification error", http.StatusUnauthorized)
			return
		}

		if adminRequired {
			if !user.IsAdmin {
				log.Printf("user %v has no permission to %v %v\n", user.Name, req.Method, req.URL)
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(req.Context(), "user", user)
		next(w, req.WithContext(ctx))
	})
}
