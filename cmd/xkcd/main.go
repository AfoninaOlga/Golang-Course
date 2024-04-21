package main

import (
	"context"
	"github.com/AfoninaOlga/xkcd/pkg/config"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	configPath := config.ParseFlag()

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Could not read config file. Error: %v\n", err)
	}

	goCnt := cfg.GoroutineCount
	if goCnt == 0 {
		goCnt = 1
		log.Println("Didn't find \"parallel\" in config file, setting number of goroutines to 1")
	}

	xkcdClient := xkcd.NewClient(cfg.Url, 10*time.Second, goCnt)

	// reading DB if exists
	comicDB, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := comicDB.Flush(); err != nil {
			log.Println(err)
		}
	}()

	getComics(&comicDB, &xkcdClient, goCnt)
}

func worker(xkcdClient *xkcd.Client, db *database.JsonDatabase, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jobs {
		if db.Exists(id) {
			continue
		}
		log.Printf("Getting Comic №%v", id)
		comic, err := xkcdClient.GetComic(id)

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
		if err := db.AddComic(id, database.Comic{Url: comic.Url, Keywords: keywords}); err != nil {
			log.Println(err)
		}
	}
}

func getComics(comicDB *database.JsonDatabase, client *xkcd.Client, goCnt int) {
	curId := 1
	defer func() {
		if err := comicDB.Flush(); err != nil {
			log.Println(err)
		}
	}()

	jobs := make(chan int, goCnt)
	var wg sync.WaitGroup
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)

	for w := 1; w <= goCnt; w++ {
		wg.Add(1)
		go worker(client, comicDB, jobs, &wg)
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

	if err := comicDB.Flush(); err != nil {
		log.Println(err)
	}
}
