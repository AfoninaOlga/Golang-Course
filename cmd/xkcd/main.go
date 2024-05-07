package main

import (
	"github.com/AfoninaOlga/xkcd/internal/adapter/client"
	"github.com/AfoninaOlga/xkcd/internal/adapter/handler"
	"github.com/AfoninaOlga/xkcd/internal/adapter/repository/json"
	"github.com/AfoninaOlga/xkcd/internal/core/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	configPath, port := ParseFlag()

	cfg, err := GetConfig(configPath)
	if err != nil {
		log.Fatalf("Could not read config file. Error: %v\n", err)
	}

	// if port flag wasn't set
	if port == -1 {
		port = cfg.Port
		// if there's no field "port" in config
		if port == 0 {
			port = 8080
		}
	}

	goCnt := cfg.GoroutineCount
	// if there's no field "parallel" in config
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

	//Filling datbase before server start
	xkcdService := service.New(comicDB, xkcdClient, 10, goCnt)
	if cnt := xkcdService.LoadComics(); cnt == 0 {
		log.Println("Nothing to load, database is up to date")
	} else {
		log.Printf("Loaded %v comics, database is up to date", cnt)
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	timeFormat := "15:04"
	updateTime, err := time.Parse(timeFormat, cfg.Time)
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
			xkcdService.LoadComics()
		}
	}()

	xkcdHandler := handler.NewXkcdHandler(xkcdService)
	router := http.NewServeMux()
	router.HandleFunc("POST /update", xkcdHandler.Update)
	router.HandleFunc("GET /pics", xkcdHandler.Search)

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      router,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
