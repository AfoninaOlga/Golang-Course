package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/pkg/config"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"github.com/schollz/progressbar/v3"
	"log"
	"time"
)

func main() {
	cnt, output, configPath := config.ParseInput()
	if cnt <= 0 {
		fmt.Printf("Nothing to do with %v comics\n", cnt)
		return
	}

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Could not read config file. Error: %v\n", err)
	}

	xkcdClient := xkcd.NewClient(cfg.Url, 10*time.Second)

	maxCnt, err := xkcdClient.GetComicsCount()
	if err != nil {
		log.Fatalf("Error getting comics count: %v\n", err)
	}

	//setting comics count limit
	if cnt > maxCnt {
		if output {
			fmt.Printf(
				"Entered number %v is bigger than existing comics count, -n is set to %v\n",
				cnt,
				maxCnt,
			)
		}
		cnt = maxCnt
	}

	// reading DB if exists
	comicDB, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	maxId := comicDB.GetMaxId()
	//return if all wanted comics are in DB and no output is needed
	if cnt <= maxId && !output {
		fmt.Println("Comics are already in the database")
		return
	}

	bar := progressbar.NewOptions(cnt, progressbar.OptionShowCount(), progressbar.OptionSetDescription("Getting comics..."))
	//adding to progressbar count of existing comics
	err = bar.Set(maxId)
	if err != nil {
		log.Println(err)
	}

	acc := 0
	// adding comics if cnt is bigger than maxId in DB
	for i := maxId + 1; i <= cnt; i++ {
		acc++
		err = bar.Add(1)
		if err != nil {
			log.Println(err)
		}

		comic, err := xkcdClient.GetComic(i)
		if err != nil {
			log.Println(err)
			//do not write to DB
			continue
		}

		keywords, err := words.StemInput(comic.Alt + " " + comic.Transcript)
		if err != nil {
			log.Printf("Stemming error in comic #%v: %v", i, err)
		}
		comicDB.AddComic(i, database.Comic{Url: comic.Url, Keywords: keywords})

		//intermediate DB writing
		if acc%50 == 0 {
			err = comicDB.Flush()
			if err != nil {
				log.Println(err)
			}
		}
	}

	if output {
		cm := comicDB.GetAll()
		displayComicMap(cm, cnt)
	}

	//if intermediate writing didn't write all comics
	if acc%50 != 0 {
		err = comicDB.Flush()
		if err != nil {
			log.Println(err)
		}
	}
}

func displayComicMap(cm map[int]database.Comic, cnt int) {
	for i := 1; i <= cnt; i++ {
		value, ok := cm[i]
		if ok {
			fmt.Printf("Comic #%v:\n", i)
			fmt.Println("\turl:", value.Url)
			fmt.Println("\tkeywords:", value.Keywords)
		}
	}
}
