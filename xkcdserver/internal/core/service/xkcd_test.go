package service

import (
	"context"
	"fmt"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/port/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"testing"
	"time"
)

var resErr = fmt.Errorf("error")

func TestXkcdService_LoadComics(t *testing.T) {
	log.SetOutput(io.Discard)
	searchLimit := 10
	goCnt := 5
	client := mocks.NewClient(t)
	//client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, nil).Times(1)
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, resErr)
	db := mocks.NewComicRepository(t)
	db.On("Size", mock.Anything).Return(399, nil).Once()
	db.On("Size", mock.Anything).Return(414, nil).Once()
	db.On("GetMaxId", mock.Anything).Return(399, nil).Once()
	db.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	cnt := xs.LoadComics(context.Background())
	require.Equal(t, 15, cnt)
}

func TestXkcdService_LoadComics_WithErrors(t *testing.T) {
	log.SetOutput(io.Discard)
	searchLimit := 10
	goCnt := 1
	client := mocks.NewClient(t)
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, nil).Once()
	err := resErr
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, err).Times(goCnt)
	db := mocks.NewComicRepository(t)
	db.On("Size", mock.Anything).Return(400, err).Once()
	db.On("Size", mock.Anything).Return(415, err).Once()
	db.On("GetMaxId", mock.Anything).Return(400, err).Once()
	db.On("Exists", mock.Anything, mock.Anything).Return(true, nil).Once()
	db.On("Exists", mock.Anything, mock.Anything).Return(false, err).Once()
	db.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	db.On("AddComic", mock.Anything, mock.Anything, mock.Anything).Return(err).Once()
	db.On("AddComic", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	cnt := xs.LoadComics(context.Background())
	require.Equal(t, 0, cnt)
}

func TestXkcdService_SetUpdateTime(t *testing.T) {
	log.SetOutput(io.Discard)
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{Alt: "letter: β"}, nil).Times(5)
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, resErr).Times(goCnt)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(500, nil).Once()
	db.On("Size", mock.Anything).Return(515, nil).Once()
	db.On("GetMaxId", mock.Anything).Return(501, nil).Once()
	db.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	db.On("AddComic", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	xs.SetUpdateTime(context.Background(), "")
}

func TestXkcdService_SetUpdateTime_Now(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, nil).Times(5)
	client.On("GetComic", mock.Anything).Return(domain.UrlComic{}, resErr).Times(goCnt)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(500, nil).Once()
	db.On("Size", mock.Anything).Return(515, nil).Once()
	db.On("GetMaxId", mock.Anything).Return(501, nil).Once()
	db.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	db.On("AddComic", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	xs.SetUpdateTime(context.Background(), time.Now().Add(-time.Minute).Format("15:04"))
}

func TestXkcdService_indexSearch(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("GetIndex", mock.Anything, mock.Anything).Return([]int{1, 2}, nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res, err := xs.indexSearch(context.Background(), []string{"appl", "hi"})
	require.Equal(t, map[int]int{1: 2, 2: 2}, res)
	require.NoError(t, err)
}

func TestXkcdService_indexSearch_Error(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("GetIndex", mock.Anything, mock.Anything).Return([]int{}, resErr)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res, err := xs.indexSearch(context.Background(), []string{"appl", "hi"})
	require.Nil(t, res)
	require.Error(t, err, "error")
}

func TestXkcdService_GetTopN(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(50, nil).Once()
	db.On("GetIndex", mock.Anything, "appl").Return([]int{2, 1}, nil)
	db.On("GetIndex", mock.Anything, "hi").Return([]int{1}, nil)
	db.On("GetUrls", mock.Anything).Return(map[int]string{}, nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res, err := xs.GetTopN(context.Background(), []string{"appl", "hi"}, 1)
	require.Equal(t, []domain.FoundComic{domain.FoundComic{Id: 1, Url: "", Count: 2}}, res)
	require.NoError(t, err)
}

func TestXkcdService_GetTopN_ErrorSize(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(0, resErr).Once()
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res, err := xs.GetTopN(context.Background(), []string{""}, 2)
	require.Nil(t, res)
	require.Error(t, err, "error")
}

func TestXkcdService_GetTopN_ErrorIndex(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(0, nil).Once()
	db.On("GetUrls", mock.Anything).Return(map[int]string{}, nil)
	db.On("GetIndex", mock.Anything, mock.Anything).Return(nil, resErr)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res, err := xs.GetTopN(context.Background(), []string{""}, 1)
	require.Nil(t, res)
	require.Error(t, err, "error")
}

func TestXkcdService_GetTopN_ErrorUrls(t *testing.T) {
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(0, nil).Once()
	db.On("GetUrls", mock.Anything).Return(map[int]string{}, resErr)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res, err := xs.GetTopN(context.Background(), []string{""}, 1)
	require.Nil(t, res)
	require.Error(t, err, "error")
}

func TestXkcdService_Search(t *testing.T) {
	log.SetOutput(io.Discard)
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(50, nil).Once()
	db.On("GetIndex", mock.Anything, mock.Anything).Return([]int{}, nil)
	db.On("GetUrls", mock.Anything).Return(map[int]string{}, nil)
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res := xs.Search(context.Background(), "letter: β")
	require.Equal(t, []domain.FoundComic{}, res)
}

func TestXkcdService_Search_Error(t *testing.T) {
	log.SetOutput(io.Discard)
	searchLimit := 10
	goCnt := 5
	client := new(mocks.Client)
	db := new(mocks.ComicRepository)
	db.On("Size", mock.Anything).Return(0, resErr).Once()
	xs := NewXkcdService(db, client, searchLimit, goCnt)
	res := xs.Search(context.Background(), "letter: β")
	require.Nil(t, res)
}
