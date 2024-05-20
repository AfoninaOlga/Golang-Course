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
	GetAll(context.Context) (map[int]domain.Comic, error)
	GetIndex(context.Context) (map[string][]int, error)
	GetMaxId(context.Context) (int, error)
	AddComic(context.Context, int, domain.Comic) error
	Exists(context.Context, int) (bool, error)
	Size(ctx context.Context) (int, error)
}

type ComicService interface {
	LoadComics(ctx context.Context) int
	Search(context.Context, string) []domain.FoundComic
}

type ComicHandler interface {
	Search(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
}
