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

func (xs *XkcdService) LoadComics() int {
	size := xs.db.Size()

	curId := 1
	if xs.db.GetMaxId()-1 == xs.db.Size() {
		curId = xs.db.GetMaxId() + 1
	}

	jobs := make(chan int, xs.goCnt)
	var wg sync.WaitGroup
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)

	for w := 1; w <= xs.goCnt; w++ {
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
	return xs.db.Size() - size
}

func (xs *XkcdService) SetUpdateTime(uTime string) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	timeFormat := "15:04"
	updateTime, err := time.Parse(timeFormat, uTime)
	now := time.Now().Format(timeFormat)
	curTime, _ := time.Parse(timeFormat, now)
	if err != nil {
		log.Println("Error parsing time from config file:", err)
		updateTime = curTime
	}

	waitTime := updateTime.Sub(curTime)
	if waitTime < 0 {
		waitTime = updateTime.Sub(curTime.Add(-24 * time.Hour))
	}
	log.Println("Scheduled update at", updateTime.Format(timeFormat), "wait time:", waitTime)

	go func() {
		<-time.After(waitTime)
		for ; ; <-ticker.C {
			log.Println("Completed scheduled comics update")
			xs.LoadComics()
		}
	}()
}

func (xs *XkcdService) Search(text string) []domain.FoundComic {
	keywords, err := stemmer.Stem(text)
	if err != nil {
		log.Println("Error stemming search query:", err)
	}
	return xs.GetTopN(keywords, xs.searchLimit)
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

func (xs *XkcdService) GetTopN(keywords []string, n int) []domain.FoundComic {
	found := make([]domain.FoundComic, 0, xs.db.Size())
	comics := xs.db.GetAll()
	counts := xs.indexSearch(keywords)

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

func (xs *XkcdService) indexSearch(keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for _, id := range xs.db.GetIndex()[k] {
			counts[id]++
		}
	}
	return counts
}

func (xs *XkcdService) dbSearch(keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for id, c := range xs.db.GetAll() {
			if _, contains := slices.BinarySearch(c.Keywords, k); contains {
				counts[id]++
			}
		}
	}
	return counts
}
