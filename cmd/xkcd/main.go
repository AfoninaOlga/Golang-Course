package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/pkg/config"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"github.com/schollz/progressbar/v3"
	"time"
)

func main() {
	cnt, output, configPath := config.ParseInput()

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		fmt.Printf("Could not read config file. Error: %v\n", err)
		return
	}

	if cfg.Url != "https://xkcd.com" {
		fmt.Printf("Unsuppotrted url %v\n", cfg.Url)
		return
	}

	xkcdClient := xkcd.NewClient(cfg.Url, 10*time.Second)

	maxCnt, err := xkcdClient.GetComicsCount()
	if err != nil {
		fmt.Printf("Error getting comics count: %v\n", err)
	}

	//setting comics count limit
	if cnt > maxCnt {
		fmt.Printf(
			"Entered number %v is bigger than existing comics count, -n is set to %v\n",
			cnt,
			maxCnt,
		)
		cnt = maxCnt
	}

	// reading DB if exists
	comicDB, err := database.New(cfg.DB)
	if err != nil {
		fmt.Println(err)
	}

	bar := progressbar.Default(int64(cnt))
	maxId := comicDB.GetMaxId()
	//return if all wanted comics are in DB and no output is needed
	if int(cnt) < maxId && !output {
		err = bar.Add(int(cnt))
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err = bar.Add(maxId)
	if err != nil {
		fmt.Println(err)
	}

	acc := 0
	// adding comics if cnt is bigger than maxId in DB
	for i := maxId + 1; i <= int(cnt); i++ {
		comic, err := xkcdClient.GetComic(i)
		if err != nil {
			fmt.Println(err)
		}

		err = bar.Add(1)
		if err != nil {
			fmt.Println(err)
		}
		keywords, err := words.StemInput(comic.Alt + " " + comic.Transcript)
		if err != nil {
			fmt.Printf("Error in comic #%v: %v", i, err)
		}
		comicDB.AddComic(i, database.Comic{Url: comic.Url, Keywords: keywords})

		//intermediate DB writing
		acc++
		if acc%50 == 0 {
			err = comicDB.Flush()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if output {
		cm := comicDB.GetAll()
		displayComicMap(cm, int(cnt))
	}

	//if intermediate writing didn't write all comics
	if acc%50 != 0 {
		err = comicDB.Flush()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func displayComicMap(cm map[int]database.Comic, cnt int) {
	for i := 1; i <= cnt; i++ {
		value := cm[i]
		fmt.Printf("Comic #%v:\n", i)
		fmt.Println("\turl:", value.Url)
		fmt.Println("\tkeywords:", value.Keywords)
	}
}
