package client

import (
	"encoding/json"
	"fmt"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newTestServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res := strings.Split(fmt.Sprintf("%v", req.URL), "/")
		if len(res) == 3 && res[2] == "info.0.json" {
			id, err := strconv.Atoi(res[1])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if id < 1 {
				w.WriteHeader(http.StatusBadRequest)
			}
			if id == 404 {
				w.WriteHeader(http.StatusNotFound)
			}
			w.Header().Set("Content-Type", "application/json")
			comic := domain.UrlComic{Id: id, Transcript: "apple"}
			_ = json.NewEncoder(w).Encode(comic)

		}
		if len(res) == 2 && res[1] == "info.0.json" {
			w.Header().Set("Content-Type", "application/json")
			comic := domain.UrlComic{Id: 3000}
			_ = json.NewEncoder(w).Encode(comic)
		}
	}))
	return server
}

func newTestServerError() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "error", http.StatusInternalServerError)
	}))
	return server
}

func TestClient_GetComic(t *testing.T) {
	server := newTestServer()
	client := NewClient(server.URL, time.Minute, 10)
	testTable := []struct {
		id          int
		expectedRes domain.UrlComic
		expectedErr error
	}{
		{
			id:          0,
			expectedRes: domain.UrlComic{},
			expectedErr: fmt.Errorf("Error getting %v/0/info.0.json, StatusCode=%v", server.URL, http.StatusBadRequest),
		},
		{
			id:          1,
			expectedRes: domain.UrlComic{Id: 1, Transcript: "apple"},
			expectedErr: nil,
		},
		{
			id:          404,
			expectedRes: domain.UrlComic{},
			expectedErr: fmt.Errorf("Error getting %v/404/info.0.json, StatusCode=%v", server.URL, http.StatusNotFound),
		},
	}
	for _, testCase := range testTable {
		c, err := client.GetComic(testCase.id)
		assert.Equal(t, testCase.expectedRes, c)
		assert.Equal(t, testCase.expectedErr, err)
	}
}

func TestClient_GetComic_Error(t *testing.T) {
	client := NewClient("", time.Minute, 10)
	c, err := client.GetComic(1)
	assert.Equal(t, domain.UrlComic{}, c)
	assert.Error(t, err)
}

func TestClient_GetComicsCount(t *testing.T) {
	server := newTestServer()
	client := NewClient(server.URL, time.Minute, 10)
	res, err := client.GetComicsCount()
	assert.NoError(t, err)
	assert.Equal(t, 3000, res)
}

func TestClient_GetComicsCount_Error(t *testing.T) {
	server := newTestServerError()
	client := NewClient(server.URL, time.Minute, 10)
	res, err := client.GetComicsCount()
	assert.Error(t, err)
	assert.Equal(t, 0, res)
}
