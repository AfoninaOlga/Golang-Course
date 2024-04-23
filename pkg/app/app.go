package app

import (
	"context"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"log"
	"os"
	"os/signal"
	"slices"
	"sync"
)

type App struct {
	client *xkcd.Client
	db     *database.JsonDatabase
}

func New(db *database.JsonDatabase, c *xkcd.Client) *App {
	return &App{client: c, db: db}
}

func (a *App) LoadComics(goCnt int) {
	curId := 1
	if a.db.GetMaxId()-1 == a.db.Size() {
		curId = a.db.GetMaxId() + 1
	}
	defer func() {
		if err := a.db.Flush(); err != nil {
			log.Println(err)
		}
	}()

	jobs := make(chan int, goCnt)
	var wg sync.WaitGroup
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)

	for w := 1; w <= goCnt; w++ {
		wg.Add(1)
		go worker(a.client, a.db, jobs, &wg)
	}

	go func() {
		wg.Wait()
		cancelFunc()
	}()

LOOP:
	for {
		select {
		case <-ctx.Done():
			close(jobs)
			break LOOP
		case jobs <- curId:
			curId++
		}
	}

	if err := a.db.Flush(); err != nil {
		log.Println(err)
	}
}

func worker(client *xkcd.Client, db *database.JsonDatabase, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jobs {
		if db.Exists(id) {
			continue
		}

		comic, err := client.GetComic(id)
		log.Printf("Got comic #%v\n", id)

		if err != nil {
			log.Println(err)
			//no more comics
			if id != 404 {
				return
			}
			continue
		}

		keywords, err := words.StemInput(comic.Alt + " " + comic.Transcript)

		if err != nil {
			log.Printf("Stemming error in comic #%v: %v", id, err)
		}

		// sorting to use binary search in DBSearch
		slices.Sort(keywords)
		if err := db.AddComic(id, database.Comic{Url: comic.Url, Keywords: keywords}); err != nil {
			log.Println(err)
		}
	}
}
