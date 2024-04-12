package main

import (
	"github.com/AfoninaOlga/xkcd/pkg/config"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"log"
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

	xkcdClient := xkcd.NewClient(cfg.Url, 10*time.Second)

	// reading DB if exists
	comicDB, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	curId := comicDB.GetMaxId() + 1

	jobs := make(chan int, goCnt)
	done := make(chan bool, goCnt)

	for w := 1; w <= goCnt; w++ {
		go getParallel(&xkcdClient, &comicDB, jobs, done)
	}

	for i, id := range comicDB.GetMissingIds() {
		jobs <- id

		if i%50 == 0 {
			err = comicDB.FlushParallel()
			if err != nil {
				log.Println(err)
			}
		}
	}
LOOP:
	for {
		select {
		case <-done:
			break LOOP
		case jobs <- curId:
			curId++
			if curId%50 == 0 {
				err = comicDB.FlushParallel()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	close(jobs)

	//waiting for workers to finish
	for w := 1; w < goCnt; w++ {
		<-done
	}

	err = comicDB.Flush()
	if err != nil {
		log.Println(err)
	}
}

func getParallel(xkcdClient *xkcd.Client, db *database.JsonDatabase, jobs <-chan int, done chan<- bool) {
	for id := range jobs {
		comic, err := xkcdClient.GetComic(id)

		if err != nil {
			log.Println(err)
			//no more comics
			if id != 404 {

				done <- true
				return
			}
			continue
		}

		keywords, err := words.StemInput(comic.Alt + " " + comic.Transcript)
		if err != nil {
			log.Printf("Stemming error in comic #%v: %v", id, err)
		}
		db.AddComicParallel(id, database.Comic{Url: comic.Url, Keywords: keywords})
	}
	done <- true
}
