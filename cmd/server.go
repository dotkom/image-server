package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dotkom/image-server/internal/api"

	"github.com/allegro/bigcache"
	gorm_adapter "github.com/dotkom/image-server/internal/storage/gorm"
	s3_adapter "github.com/dotkom/image-server/internal/storage/s3"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  `Start the server`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	fs, err := s3_adapter.New(viper.GetString(s3Bucket))
	if err != nil {
		log.Fatal("Failed to create storage adapter", err)
	}

	ms := gorm_adapter.New(gorm_adapter.DBDriver(viper.GetString(dbDriver)), viper.GetString(dbDSN))

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Fatal("Failed to create cache", err)
	}

	router := mux.NewRouter()
	api := api.New(fs, ms, router, cache)

	server := &http.Server{
		Addr:         viper.GetString(listenAddr),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      api,
	}

	go func() {
		log.Infof("Server listening to %s", viper.GetString(listenAddr))
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
