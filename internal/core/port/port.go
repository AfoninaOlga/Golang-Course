package port

import "github.com/AfoninaOlga/xkcd/internal/core/domain"

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

type Stemmer interface {
	Stem(string) []string
}
