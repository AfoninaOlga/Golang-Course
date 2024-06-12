package handler

import (
	"bytes"
	"fmt"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var resErr = fmt.Errorf("error")

func TestAuthHandler_Register_BodeError(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(""))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Register(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	svc.On("Register", mock.Anything, mock.Anything).Return(false, resErr)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString("{\"user\": \"\", \"password\": \"\"}"))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Register(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAuthHandler_Register_UserExists(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	svc.On("Register", mock.Anything, mock.Anything).Return(false, nil)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString("{\"user\": \"\", \"password\": \"\"}"))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Register(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthHandler_Register(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	svc.On("Register", mock.Anything, mock.Anything).Return(true, nil)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString("{\"user\": \"\", \"password\": \"\"}"))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Register(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestAuthHandler_Login_ErrorBody(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(""))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Login(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthHandler_Login_ServiceError(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	svc.On("Login", mock.Anything, mock.Anything).Return("", resErr)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString("{\"user\": \"\", \"password\": \"\"}"))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Login(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthHandler_Login(t *testing.T) {
	log.SetOutput(io.Discard)
	svc := new(mocks.AuthService)
	svc.On("Login", mock.Anything, mock.Anything).Return("token", nil)
	ah := NewAuthHandler(svc)
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString("{\"user\": \"\", \"password\": \"\"}"))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	ah.Login(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
