package main

import (
	"context"
	"github.com/AfoninaOlga/xkcd/webserver/internal/handler"
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

	log.Println("xkcd server will listening on", port)

	ttl := cfg.TokenDuration
	if ttl == 0 {
		ttl = 10
	}
	h := handler.NewHandler(cfg.Api, time.Duration(ttl)*time.Minute)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/comics", http.StatusSeeOther)
	})
	router.HandleFunc("/login", h.Login)
	router.HandleFunc("/comics", h.Search)

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
