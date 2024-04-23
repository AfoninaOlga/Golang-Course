package app

import (
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"testing"
)

var keywords = []string{"account", "zip", "zero", "know", "question", "complain", "overlap", "live", "guess", "truth", "save", "bug"}

func BenchmarkDB(b *testing.B) {
	db, _ := database.New("../../database.json")
	app := App{db: db}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.dbSearch(keywords)
	}
}

func BenchmarkIndex(b *testing.B) {
	db, _ := database.New("../../database.json")
	app := App{db: db}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.indexSearch(keywords)
	}
}
