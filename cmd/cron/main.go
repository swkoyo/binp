package main

import (
	"binp/scheduler"
	"binp/storage"
	"binp/util"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	util.InitLogger()
	logger := util.GetLogger()

	store, err := storage.NewStore()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create store")
	}

	if err := store.Init(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize store")
	}

	runner := scheduler.NewScheduler()
	runner.Init(store)

	logger.Info().Msg("Starting scheduler...")
	runner.Start()
	logger.Info().Msg("Scheduler started!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info().Msg("Shutting down scheduler...")
	<-runner.Stop().Done()
	logger.Info().Msg("Scheduler stopped!")
}
