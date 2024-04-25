package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/pkg/app"
	"github.com/AfoninaOlga/xkcd/pkg/config"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"log"
	"time"
)

func main() {
	configPath, sQuery, useIndex := config.ParseFlag()

	if sQuery == "" {
		return
	}

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

	a := app.New(comicDB, xkcdClient)

	a.LoadComics(goCnt)

	stemmed, err := words.StemInput(sQuery)
	if err != nil {
		log.Println(err)
	}
	for id, comic := range a.GetTopN(stemmed, 10, useIndex) {
		fmt.Printf("#%v relevant (%v overlap): %v\n", id+1, comic.Count, comic.Url)
	}
}
