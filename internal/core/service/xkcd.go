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

func NewXkcdService(db port.ComicRepository, c port.Client, searchLimit int, goCnt int) *XkcdService {
	return &XkcdService{client: c, db: db, searchLimit: searchLimit, goCnt: goCnt}
}

func (xs *XkcdService) LoadComics(ctx context.Context) int {
	size, err := xs.db.Size(ctx)
	if err != nil {
		log.Println("Error getting comic table size:", err)
	}
	maxId, err := xs.db.GetMaxId(ctx)
	if err != nil {
		log.Println("Error getting max id from comic table:", err)
	}

	curId := 1
	if maxId-1 == size {
		curId = maxId + 1
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
	newSize, err := xs.db.Size(ctx)
	if err != nil {
		log.Println("Error getting comic table new size:", err)
		newSize = size
	}
	return newSize - size
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
	res, err := xs.GetTopN(ctx, keywords, xs.searchLimit)
	if err != nil {
		log.Println("Error getting top N URL's:", err)
		return nil
	}
	return res
}

func worker(ctx context.Context, client port.Client, db port.ComicRepository, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jobs {
		exists, err := db.Exists(ctx, id)
		if err != nil {
			log.Println("Error defining comic existence in database:", err)
		}
		if exists {
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

func (xs *XkcdService) GetTopN(ctx context.Context, keywords []string, n int) ([]domain.FoundComic, error) {
	size, err := xs.db.Size(ctx)
	if err != nil {
		return nil, err
	}
	found := make([]domain.FoundComic, 0, size)
	urls, err := xs.db.GetUrls(ctx)
	if err != nil {
		return nil, err
	}
	counts, err := xs.indexSearch(ctx, keywords)
	if err != nil {
		return nil, err
	}

	for id, cnt := range counts {
		found = append(found, domain.FoundComic{Id: id, Count: cnt, Url: urls[id]})
	}

	slices.SortFunc(found, func(a, b domain.FoundComic) int {
		return b.Count - a.Count
	})
	if len(found) < n {
		return found, nil
	}
	return found[:n], nil
}

func (xs *XkcdService) indexSearch(ctx context.Context, keywords []string) (map[int]int, error) {
	counts := make(map[int]int)
	for _, k := range keywords {
		index, err := xs.db.GetIndex(ctx, k)
		if err != nil {
			return nil, err
		}
		for _, id := range index {
			counts[id]++
		}
	}
	return counts, nil
}
