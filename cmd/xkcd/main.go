package main

import (
	"context"
	"database/sql"
	"github.com/AfoninaOlga/xkcd/internal/adapter/client"
	"github.com/AfoninaOlga/xkcd/internal/adapter/handler"
	"github.com/AfoninaOlga/xkcd/internal/adapter/repository/sqlite"
	"github.com/AfoninaOlga/xkcd/internal/core/service"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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
	db, err := sql.Open("sqlite3", cfg.Database)
	if err != nil {
		log.Fatalln("error opening database:", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalln("error connecting to database:", err)
	}

	comicDB := sqlite.New(db)
	//Filling database before server start
	xkcdService := service.New(comicDB, xkcdClient, 10, goCnt)
	if cnt := xkcdService.LoadComics(ctx); cnt == 0 {
		log.Println("Nothing to load, database is up to date")
	} else {
		log.Printf("Loaded %v comics, database is up to date", cnt)
	}
	xkcdService.SetUpdateTime(ctx, cfg.Time)

	xkcdHandler := handler.NewXkcdHandler(xkcdService)
	router := http.NewServeMux()
	router.HandleFunc("POST /update", xkcdHandler.Update)
	router.HandleFunc("GET /pics", xkcdHandler.Search)

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      router,
		BaseContext:  func(net.Listener) context.Context { return ctx },
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}
}
