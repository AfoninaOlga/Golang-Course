package port

import (
	"context"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"net/http"
)

type Client interface {
	GetComic(int) (domain.UrlComic, error)
}

type ComicRepository interface {
	GetAll() map[int]domain.Comic
	GetIndex() map[string][]int
	GetMaxId() int
	AddComic(id int, c domain.Comic) error
	Exists(int) bool
	Size() int
	Flush() error
}

type ComicService interface {
	LoadComics(ctx context.Context) int
	Search(string) []domain.FoundComic
}

type Stemmer interface {
	Stem(string) ([]string, error)
}

type ComicHandler interface {
	Search(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
}
