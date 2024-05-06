package service

import (
	"context"
	"github.com/AfoninaOlga/xkcd/internal/adapter/stemmer"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/internal/core/port"
	"log"
	"os"
	"os/signal"
	"slices"
	"sync"
)

type XkcdService struct {
	client port.Client
	db     port.ComicRepository
}

func New(db port.ComicRepository, c port.Client) *XkcdService {
	return &XkcdService{client: c, db: db}
}

func (xs *XkcdService) LoadComics(goCnt int) {
	curId := 1
	if xs.db.GetMaxId()-1 == xs.db.Size() {
		curId = xs.db.GetMaxId() + 1
	}

	jobs := make(chan int, goCnt)
	var wg sync.WaitGroup
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)

	for w := 1; w <= goCnt; w++ {
		wg.Add(1)
		go worker(xs.client, xs.db, jobs, &wg)
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

	if err := xs.db.Flush(); err != nil {
		log.Println(err)
	}
}

func worker(client port.Client, db port.ComicRepository, jobs <-chan int, wg *sync.WaitGroup) {
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

		keywords, err := stemmer.Stem(comic.Alt + " " + comic.Transcript + " " + comic.Title)

		if err != nil {
			log.Printf("Stemming error in comic #%v: %v", id, err)
		}

		// sorting to use binary search in DBSearch
		slices.Sort(keywords)

		if err := db.AddComic(id, domain.Comic{Url: comic.Url, Keywords: keywords}); err != nil {
			log.Println(err)
		}
	}
}
