package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/allegro/bigcache"
	"github.com/dotkom/image-server/api"
	"github.com/dotkom/image-server/models"
	"github.com/dotkom/image-server/storage/adapters"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type config struct {
	ListenAddr string
}

func main() {
	c := loadConfig()
	storageAdapter, err := adapters.New("images.dotkom")
	if err != nil {
		log.Fatal("Failed to create storage adapter", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to establish database connection", err)
	}

	log.Info("Migrating database models")
	models.Migrate(db)

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Fatal("Failed to create cache", err)
	}
	router := mux.NewRouter()
	api := api.New(storageAdapter, db, router, cache)

	server := &http.Server{
		Addr:         c.ListenAddr,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      api,
	}

	go func() {
		log.Infof("Server listening to %s", c.ListenAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	sig := <-channel
	log.Infof("Recieved signal: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	server.Shutdown(ctx)

	log.Info("Shutting down")
	os.Exit(0)
}

func loadConfig() *config {
	c := &config{}
	log.Info("Loading config")
	flag.StringVar(&c.ListenAddr, "listen-addr", "0.0.0.0:8080", "Address for server to listen to in the form ip:port")
	flag.Parse()
	return c
}
