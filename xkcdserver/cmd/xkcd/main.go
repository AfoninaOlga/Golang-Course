package main

import (
	"context"
	"database/sql"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/adapter/client"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/adapter/handler"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/adapter/repository/sqlite"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/service"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	log.Println("web server will listening on", port)

	goCnt := cfg.GoroutineCount
	// if there's no field "parallel" in config
	if goCnt == 0 {
		goCnt = 1
		log.Println("Didn't find \"parallel\" in config file, setting number of goroutines to 1")
	}

	rl := cfg.RateLimit
	if rl <= 0 {
		rl = 10
	}
	rateLimiter := handler.NewRateLimiter(ctx, rl, 5*time.Minute)
	cl := cfg.ConcurrencyLimit
	if cl <= 0 {
		cl = 10
	}
	concurrencyLimiter := handler.NewConcurrencyLimiter(cl)

	xkcdClient := client.NewClient(cfg.Url, 10*time.Second, goCnt)

	// Trying to connect to database
	db, err := sql.Open("sqlite3", cfg.Database)
	if err != nil {
		log.Fatalln("error opening database:", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalln("error connecting to database:", err)
	}

	// Trying to run migrations
	comicDB := sqlite.NewComicDB(db)
	userDB := sqlite.NewUserDB(db)
	if err = comicDB.RunMigrationUp(); err != nil {
		log.Fatalln("error running migration:", err)
	}

	//Filling database before server start
	xkcdService := service.NewXkcdService(comicDB, xkcdClient, 10, goCnt)
	ttl := cfg.TokenDuration
	if ttl == 0 {
		ttl = 10
	}
	authService := service.NewAuthService(userDB, "quokka", time.Duration(ttl)*time.Minute)
	if cnt := xkcdService.LoadComics(ctx); cnt == 0 {
		log.Println("Nothing to load, database is up to date")
	} else {
		log.Printf("Loaded %v comics, database is up to date", cnt)
	}
	xkcdService.SetUpdateTime(ctx, cfg.Time)

	xkcdHandler := handler.NewXkcdHandler(xkcdService)
	authHandler := handler.NewAuthHandler(authService)
	router := http.NewServeMux()
	router.HandleFunc("POST /update", handler.Auth(true, authService, handler.RateLimiting(rateLimiter, xkcdHandler.Update)))
	router.HandleFunc("GET /pics", handler.Auth(false, authService, handler.RateLimiting(rateLimiter, xkcdHandler.Search)))
	router.HandleFunc("POST /login", handler.RateLimiting(rateLimiter, authHandler.Login))
	router.HandleFunc("POST /register", handler.RateLimiting(rateLimiter, authHandler.Register))

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      handler.ConcurrencyLimiting(concurrencyLimiter, router.ServeHTTP),
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
