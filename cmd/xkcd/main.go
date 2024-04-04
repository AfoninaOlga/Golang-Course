package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/utils"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"github.com/schollz/progressbar/v3"
)

type ComicsBase interface {
	Flush() error
	GetAll() map[int]database.Comic
	AddComic(id int, c database.Comic)
	GetMaxId() int
}

func main() {
	cnt, output, configPath := utils.ParseInput()

	config, err := utils.GetConfig(configPath)
	if err != nil {
		fmt.Printf("Could not read config file. Error: %v\n", err)
		return
	}

	if config.Url != "https://xkcd.com" {
		fmt.Printf("Unsuppotrted url %v\n", config.Url)
		return
	}

	xkcdClient := xkcd.NewClient(config.Url)

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

	var jsonDb database.JsonDatabase
	// reading DB if exists
	err = jsonDb.Init(config.DB)
	if err != nil {
		fmt.Println(err)
	}
	var comicDb ComicsBase = &jsonDb

	bar := progressbar.Default(int64(cnt))
	maxId := comicDb.GetMaxId()
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
		comic, err := xkcdClient.GetComicResponse(i)
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
		comicDb.AddComic(i, database.Comic{Url: comic.Url, Keywords: keywords})

		//intermediate DB writing
		acc++
		if acc%50 == 0 {
			err = comicDb.Flush()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if output {
		cm := comicDb.GetAll()
		displayComicMap(cm, int(cnt))
	}

	//if intermediate writing didn't write all comics
	if acc%50 != 0 {
		err = comicDb.Flush()
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
