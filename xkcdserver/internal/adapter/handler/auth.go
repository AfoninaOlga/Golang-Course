package handler

import (
	"encoding/json"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/port"
	"log"
	"net/http"
)

type AuthHandler struct {
	svc port.AuthService
}

func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, req *http.Request) {
	user := domain.User{}
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Printf("error decoding request: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := ah.svc.Login(req.Context(), user)
	if err != nil {
		log.Printf("error logging in: %v\n", err)
		http.Error(w, "Error logging in", http.StatusUnauthorized)
		return
	}
	resp := struct {
		Token string `json:"token"`
	}{Token: token}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Panic("error encoding response:", err)
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, req *http.Request) {
	user := domain.User{}
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Printf("error decoding request: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	added, err := ah.svc.Register(req.Context(), user)

	if err != nil {
		log.Printf("error registering user: %v\n", err)
		http.Error(w, "Error registering", http.StatusInternalServerError)
		return
	}

	//user with requested name already exists
	if !added {
		log.Printf("user %v not registered cause already exists\n", user.Name)
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
