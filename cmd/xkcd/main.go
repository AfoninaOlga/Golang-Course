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
		port = cfg.port
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

	xkcdService := service.New(comicDB, xkcdClient, 10, goCnt)
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
