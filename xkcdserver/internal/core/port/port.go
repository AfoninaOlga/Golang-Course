package port

import (
	"context"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/domain"
	"net/http"
)

//go:generate mockery --all

type Client interface {
	GetComic(int) (domain.UrlComic, error)
}

type ComicRepository interface {
	GetAll(context.Context) (map[int]domain.Comic, error)
	GetIndex(context.Context, string) ([]int, error)
	GetMaxId(context.Context) (int, error)
	AddComic(context.Context, int, domain.Comic) error
	Exists(context.Context, int) (bool, error)
	Size(context.Context) (int, error)
	GetUrls(context.Context) (map[int]string, error)
	RunMigrationUp() error
}

type UserRepository interface {
	Add(context.Context, domain.User) error
	GetByName(context.Context, string) (*domain.User, error)
}

type ComicService interface {
	LoadComics(ctx context.Context) int
	Search(context.Context, string) []domain.FoundComic
}

type AuthService interface {
	Login(context.Context, domain.User) (string, error)
	Register(context.Context, domain.User) (bool, error)
	GetUserByToken(context.Context, string) (*domain.User, error)
}

type ComicHandler interface {
	Search(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
}
