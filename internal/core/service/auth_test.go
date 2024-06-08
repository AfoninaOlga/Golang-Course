package service

import (
	"context"
	"fmt"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/internal/core/port/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAuthService_Register_NewUser(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(nil, nil).Once()
	uRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
	as := NewAuthService(uRepo, "quokka", time.Minute)
	added, err := as.Register(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.Equal(t, true, added)
}

func TestAuthService_Register_ErrorAdd(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	retErr := fmt.Errorf("error")
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(nil, nil).Once()
	uRepo.On("Add", mock.Anything, mock.Anything).Return(retErr)
	as := NewAuthService(uRepo, "quokka", time.Minute)
	added, err := as.Register(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.Equal(t, fmt.Errorf("error adding user: %v", retErr), err)
	require.Equal(t, false, added)
}

func TestAuthService_Register_ExistingUser(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(&u, nil).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	added, err := as.Register(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.Equal(t, false, added)
}

func TestAuthService_Register_RepositoryErrorGet(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uRepo := mocks.NewUserRepository(t)
	err := fmt.Errorf("error")
	uRepo.On("GetByName", mock.Anything, mock.Anything).Return(nil, err).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	added, err := as.Register(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.Error(t, err)
	require.Equal(t, false, added)
}

func TestAuthService_Login_ErrorAdding(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	retErr := fmt.Errorf("error")
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(nil, retErr).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := as.Login(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.Equal(t, fmt.Errorf("error getting user %v: %v", u.Name, retErr), err)
	require.Equal(t, "", token)
}

func TestAuthService_Login_ExistingUser(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uHashed := domain.User{Name: "user", Password: "$2a$10$vwU67W22FtqKB6xWHGmfButAtZgSbbJDrLm9brwv/phXc5CTo/i7y"}
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(&uHashed, nil).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := as.Login(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.NotEqual(t, "", token)
}

func TestAuthService_Login_NotExistingUser(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(nil, nil).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := as.Login(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.Equal(t, fmt.Errorf("no user called %v", u.Name), err)
	require.Equal(t, "", token)
}

func TestAuthService_Login_IncorrectPassword(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(&u, nil).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := as.Login(context.Background(), u)
	uRepo.AssertExpectations(t)
	require.Equal(t, fmt.Errorf("incorrect password"), err)
	require.Equal(t, "", token)
}

func TestAuthService_GetUserByToken_Invalid(t *testing.T) {
	uRepo := mocks.NewUserRepository(t)
	as := NewAuthService(uRepo, "quokka", time.Minute)
	user, err := as.GetUserByToken(context.Background(), "")
	uRepo.AssertExpectations(t)
	require.Equal(t, fmt.Errorf("invalid token"), err)
	require.Nil(t, user)
}

func TestAuthService_GetUserByToken(t *testing.T) {
	u := domain.User{Name: "user", Password: "pass123"}
	uRepo := mocks.NewUserRepository(t)
	uRepo.On("GetByName", mock.Anything, u.Name).Return(&u, nil).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.Name, "exp": time.Now().Add(as.tokenDuration).Unix()}).SignedString(as.secret)
	require.NoError(t, err)
	user, err := as.GetUserByToken(context.Background(), token)
	uRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.Equal(t, u.Name, user.Name)
}

func TestAuthService_GetUserByToken_RepositoryError(t *testing.T) {
	uRepo := mocks.NewUserRepository(t)
	retErr := fmt.Errorf("error")
	uRepo.On("GetByName", mock.Anything, mock.Anything).Return(nil, retErr).Once()
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user", "exp": time.Now().Add(as.tokenDuration).Unix()}).SignedString(as.secret)
	require.NoError(t, err)
	user, err := as.GetUserByToken(context.Background(), token)
	uRepo.AssertExpectations(t)
	require.Equal(t, retErr, err)
	require.Nil(t, user)
}

func TestAuthService_GetUserByToken_InvalidUser(t *testing.T) {
	uRepo := mocks.NewUserRepository(t)
	as := NewAuthService(uRepo, "quokka", time.Minute)
	token, err := jwt.New(jwt.SigningMethodHS256).SignedString(as.secret)
	require.NoError(t, err)
	user, err := as.GetUserByToken(context.Background(), token)
	uRepo.AssertExpectations(t)
	require.Error(t, fmt.Errorf("invalid username in claims"), err)
	require.Nil(t, user)
}
