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
	"time"
)

type XkcdService struct {
	client      port.Client
	db          port.ComicRepository
	searchLimit int
	goCnt       int
}

func New(db port.ComicRepository, c port.Client, searchLimit int, goCnt int) *XkcdService {
	return &XkcdService{client: c, db: db, searchLimit: searchLimit, goCnt: goCnt}
}

func (xs *XkcdService) LoadComics(ctx context.Context) int {
	size := xs.db.Size(ctx)

	curId := 1
	if xs.db.GetMaxId(ctx)-1 == xs.db.Size(ctx) {
		curId = xs.db.GetMaxId(ctx) + 1
	}

	jobs := make(chan int, xs.goCnt)
	var wg sync.WaitGroup
	loadCtx, cancelFunc := signal.NotifyContext(ctx, os.Interrupt)

	for w := 1; w <= xs.goCnt; w++ {
		wg.Add(1)
		go worker(loadCtx, xs.client, xs.db, jobs, &wg)
	}

	go func() {
		wg.Wait()
		cancelFunc()
	}()

LOOP:
	for {
		select {
		case <-loadCtx.Done():
			close(jobs)
			break LOOP
		case jobs <- curId:
			curId++
		}
	}
	return xs.db.Size(ctx) - size
}

func (xs *XkcdService) SetUpdateTime(ctx context.Context, uTime string) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	timeFormat := "15:04"
	updateTime, err := time.Parse(timeFormat, uTime)
	curTime, _ := time.Parse(timeFormat, time.Now().Format(timeFormat))
	if err != nil {
		log.Println("Error parsing update time:", err)
		updateTime = curTime
	}

	if updateTime.Before(curTime) {
		updateTime = updateTime.Add(24 * time.Hour)
	}
	waitTime := updateTime.Sub(curTime)
	log.Println("Scheduled update at", updateTime.Format(timeFormat), "wait time:", waitTime)

	go func() {
		<-time.After(waitTime)
		for ; ; <-ticker.C {
			log.Println("Completed scheduled comics update")
			xs.LoadComics(ctx)
		}
	}()
}

func (xs *XkcdService) Search(ctx context.Context, text string) []domain.FoundComic {
	keywords, err := stemmer.Stem(text)
	if err != nil {
		log.Println("Error stemming search query:", err)
	}
	return xs.GetTopN(ctx, keywords, xs.searchLimit)
}

func worker(ctx context.Context, client port.Client, db port.ComicRepository, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jobs {
		if db.Exists(ctx, id) {
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

		if err := db.AddComic(ctx, id, domain.Comic{Url: comic.Url, Keywords: keywords}); err != nil {
			log.Println(err)
		}
	}
}

func (xs *XkcdService) GetTopN(ctx context.Context, keywords []string, n int) []domain.FoundComic {
	found := make([]domain.FoundComic, 0, xs.db.Size(ctx))
	comics := xs.db.GetAll(ctx)
	counts := xs.indexSearch(ctx, keywords)

	for id, cnt := range counts {
		found = append(found, domain.FoundComic{Id: id, Count: cnt, Url: comics[id].Url})
	}

	slices.SortFunc(found, func(a, b domain.FoundComic) int {
		return b.Count - a.Count
	})
	if len(found) < n {
		return found
	}
	return found[:n]
}

func (xs *XkcdService) indexSearch(ctx context.Context, keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for _, id := range xs.db.GetIndex(ctx)[k] {
			counts[id]++
		}
	}
	return counts
}

func (xs *XkcdService) dbSearch(ctx context.Context, keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for id, c := range xs.db.GetAll(ctx) {
			if _, contains := slices.BinarySearch(c.Keywords, k); contains {
				counts[id]++
			}
		}
	}
	return counts
}
