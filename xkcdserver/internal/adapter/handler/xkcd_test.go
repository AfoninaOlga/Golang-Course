package handler

import (
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXkcdHandler_Search_EmptyQuery(t *testing.T) {
	svc := new(mocks.ComicService)
	xh := NewXkcdHandler(svc)
	req, err := http.NewRequest("GET", "/pics", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	xh.Search(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestXkcdHandler_Search_ServiceError(t *testing.T) {
	svc := new(mocks.ComicService)
	svc.On("Search", mock.Anything, mock.Anything).Return(nil)
	xh := NewXkcdHandler(svc)
	req, err := http.NewRequest("GET", "/pics?search=apple", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	xh.Search(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestXkcdHandler_Update(t *testing.T) {
	svc := new(mocks.ComicService)
	svc.On("LoadComics", mock.Anything).Return(0)
	xh := NewXkcdHandler(svc)
	req, err := http.NewRequest("POST", "/update", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	xh.Update(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestXkcdHandler_UpdateGoroutine(t *testing.T) {
	svc := new(mocks.ComicService)
	svc.On("LoadComics", mock.Anything).Return(0)
	xh := NewXkcdHandler(svc)
	xh.mtx.Lock()
	defer xh.mtx.Unlock()
	req, err := http.NewRequest("POST", "/update", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	xh.Update(recorder, req)
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}
