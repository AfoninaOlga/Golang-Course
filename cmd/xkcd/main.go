package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/utils"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"github.com/schollz/progressbar/v3"
)

func main() {
	cnt, output := utils.ParseInput()

	c, err := utils.GetConfig("config.yaml")
	if err != nil {
		fmt.Printf("Could not read config file. Error: %v\n", err)
		return
	}

	url := c.Url
	db := c.DB

	if url != "https://xkcd.com" {
		fmt.Printf("Unsuppotrted url %v\n", c.Url)
		return
	}

	xkcdClient := xkcd.NewClient(url)

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

	maxId := database.GetMaxIdFromDB(db)
	//return if all wanted comics are in DB and no output is needed
	if int(cnt) < maxId && !output {
		return
	}

	// reading existing comics
	cm, err := database.ReadFile(db)
	if err != nil {
		fmt.Printf("Error reading DB: %v\n", err)
	}

	// adding comics if cnt is bigger than maxId in DB
	bar := progressbar.Default(int64(cnt))
	err = bar.Add(maxId)
	if err != nil {
		fmt.Println(err)
	}

	acc := 0

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
		cm[i] = database.Comic{Url: comic.Url, Keywords: keywords}

		//intermediate DB writing
		acc++
		if acc%50 == 0 {
			err = database.WriteFile(db, cm, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if output {
		database.DisplayComicMap(cm, int(cnt))
	}

	//if intermediate writing didn't write all comics
	if acc%50 != 0 {
		err = database.WriteFile(db, cm, int(cnt))
		if err != nil {
			fmt.Println(err)
		}
	}
}
