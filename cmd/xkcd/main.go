package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/internal/adapter/client"
	"github.com/AfoninaOlga/xkcd/internal/adapter/repository/json"
	"github.com/AfoninaOlga/xkcd/internal/core/service"
	"log"
	"time"
)

func main() {
	configPath, sQuery, useIndex := ParseFlag()

	if sQuery == "" {
		return
	}

	cfg, err := GetConfig(configPath)
	if err != nil {
		log.Fatalf("Could not read config file. Error: %v\n", err)
	}

	goCnt := cfg.GoroutineCount
	if goCnt == 0 {
		goCnt = 1
		log.Println("Didn't find \"parallel\" in config file, setting number of goroutines to 1")
	}

	xkcdClient := client.NewClient(cfg.Url, 10*time.Second, goCnt)

	// reading DB if exists
	comicDB, err := json.New(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	a := service.New(comicDB, xkcdClient)

	a.LoadComics(goCnt)

	stemmed, err := client.StemInput(sQuery)
	if err != nil {
		log.Println(err)
	}
	for id, comic := range a.GetTopN(stemmed, 10, useIndex) {
		fmt.Printf("#%v relevant (%v overlap): %v\n", id+1, comic.Count, comic.Url)
	}
}
