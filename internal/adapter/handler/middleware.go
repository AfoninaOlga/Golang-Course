package handler

import (
	"context"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/internal/core/port"
	"log"
	"net"
	"net/http"
)

type ctxKey string

func Auth(adminRequired bool, authService port.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, err := authService.GetUserByToken(req.Context(), req.Header.Get("Authorization"))
		if err != nil {
			log.Printf("authentication error: %v\n", err)
			http.Error(w, "Authentication error", http.StatusUnauthorized)
			return
		}

		if adminRequired {
			if !user.IsAdmin {
				log.Printf("user %v has no permission to %v %v\n", user.Name, req.Method, req.URL)
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(req.Context(), ctxKey("user"), user)
		next(w, req.WithContext(ctx))
	})
}

func RateLimiting(limiter *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			err    error
			client string
		)
		user := req.Context().Value(ctxKey("user"))
		if user == nil {
			client, _, err = net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
		} else {
			client = (user.(*domain.User)).Name
		}

		if limiter.Allow(client) {
			next(w, req)
		} else {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}
	})
}

func ConcurrencyLimiting(limiter *ConcurrencyLimiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		limiter.Add()
		defer limiter.Remove()
		next(w, req)
	})
}
