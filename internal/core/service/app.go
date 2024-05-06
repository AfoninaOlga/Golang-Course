package service

import (
	"context"
	"github.com/AfoninaOlga/xkcd/internal/adapter/client"
	"github.com/AfoninaOlga/xkcd/internal/adapter/repository/json"
	"log"
	"os"
	"os/signal"
	"slices"
	"sync"
)

type App struct {
	client *client.Client
	db     *json.JsonDatabase
}

func New(db *json.JsonDatabase, c *client.Client) *App {
	return &App{client: c, db: db}
}

func (a *App) LoadComics(goCnt int) {
	curId := 1
	if a.db.GetMaxId()-1 == a.db.Size() {
		curId = a.db.GetMaxId() + 1
	}

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

func worker(client *client.Client, db *json.JsonDatabase, jobs <-chan int, wg *sync.WaitGroup) {
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

		keywords, err := client.StemInput(comic.Alt + " " + comic.Transcript + " " + comic.Title)

		if err != nil {
			log.Printf("Stemming error in comic #%v: %v", id, err)
		}

		// sorting to use binary search in DBSearch
		slices.Sort(keywords)
		if err := db.AddComic(id, json.Comic{Url: comic.Url, Keywords: keywords}); err != nil {
			log.Println(err)
		}
	}
}
