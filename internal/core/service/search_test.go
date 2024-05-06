package service

import (
	"github.com/AfoninaOlga/xkcd/internal/adapter/repository/json"
	"testing"
)

var keywords = []string{"account", "zip", "zero", "know", "question", "complain", "overlap", "live", "guess", "truth", "save", "bug"}

func BenchmarkDB(b *testing.B) {
	db, _ := json.New("../../database.json")
	app := XkcdService{db: db}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.dbSearch(keywords)
	}
}

func BenchmarkIndex(b *testing.B) {
	db, _ := json.New("../../database.json")
	app := XkcdService{db: db}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.indexSearch(keywords)
	}
}
