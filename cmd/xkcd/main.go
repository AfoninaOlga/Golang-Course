package main

import (
	"fmt"
	"github.com/AfoninaOlga/xkcd/pkg/database"
	"github.com/AfoninaOlga/xkcd/pkg/utils"
	"github.com/AfoninaOlga/xkcd/pkg/words"
	"github.com/AfoninaOlga/xkcd/pkg/xkcd"
	"strconv"
)

func main() {
	cntIsSet, cnt, output := utils.ParseInput()

	c, err := utils.GetConfig("config.yaml")
	if err != nil {
		fmt.Printf("Could not read config file. Error: %v\n", err)
		return
	}

	url := c.Url

	if url != "https://xkcd.com" {
		fmt.Printf("Unsuppotrted url %v\n", c.Url)
		return
	}

	xkcdClient := xkcd.NewClient(url)

	maxCnt, err := xkcdClient.GetComicsCount()
	if err != nil {
		fmt.Printf("Error getting comics count: %v\n", err)
	}

	if cnt > maxCnt {
		fmt.Printf(
			"Entered number %v is bigger than existing comics count, -n is set to %v\n",
			cnt,
			maxCnt,
		)
		cnt = maxCnt
	}

	if !cntIsSet {
		cnt = maxCnt
	}

	cm := database.ComicMap{}

	for i := 1; i <= int(cnt); i++ {
		id := strconv.Itoa(i)
		comic, err := xkcdClient.GetComicResponse(i)
		if err != nil {
			fmt.Println(err)
		}
		keywords, err := words.StemInput(comic.Alt + " " + comic.Transcript)
		if err != nil {
			fmt.Println(err)
		}
		cm[id] = database.Comic{Url: comic.Url, Keywords: keywords}
	}

	if output {
		database.DisplayComicMap(cm)
	}

	err = database.WriteFile(c.DB, cm)
	if err != nil {
		fmt.Println(err)
	}
}
