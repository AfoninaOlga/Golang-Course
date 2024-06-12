package service

import (
	"context"
	"fmt"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userRepo      port.UserRepository
	secret        []byte
	tokenDuration time.Duration
}

func NewAuthService(userRepo port.UserRepository, secret string, duration time.Duration) *AuthService {
	return &AuthService{userRepo: userRepo, secret: []byte(secret), tokenDuration: duration}
}

func (as *AuthService) Login(ctx context.Context, u domain.User) (string, error) {
	dbUser, err := as.userRepo.GetByName(ctx, u.Name)
	if err != nil {
		return "", fmt.Errorf("error getting user %v: %v", u.Name, err)
	}

	if dbUser == nil {
		return "", fmt.Errorf("no user called %v", u.Name)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password)); err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.Name, "exp": time.Now().Add(as.tokenDuration).Unix()})

	t, err := token.SignedString(as.secret)
	if err != nil {
		return "", fmt.Errorf("errror signing token: %v", err)
	}
	return t, nil
}

func (as *AuthService) Register(ctx context.Context, u domain.User) (bool, error) {
	dbUser, err := as.userRepo.GetByName(ctx, u.Name)
	if err != nil {
		return false, fmt.Errorf("error checking user %v existance: %v", u.Name, err)
	}

	//if user already exists
	if dbUser != nil {
		return false, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("error hashing password: %v", err)
	}

	u.Password = string(hashedPassword)
	err = as.userRepo.Add(ctx, u)

	if err != nil {
		return false, fmt.Errorf("error adding user: %v", err)
	}
	return true, nil
}

func (as *AuthService) GetUserByToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return as.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	name, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in claims")
	}
	user, err := as.userRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return user, nil
}
