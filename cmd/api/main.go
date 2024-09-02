package main

import (
	"binp/scheduler"
	"binp/server"
	"binp/storage"
	"binp/util"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	util.InitLogger()
	logger := util.GetLogger()

	err := util.GenerateChromaCSS()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate chroma.css")
	}

	store, err := storage.NewStore()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create new store")
	}

	if err := store.Init(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initalize store")
	}

	runner := scheduler.NewScheduler()
	runner.Init(store)
	logger.Info().Msg("Starting scheduler...")
	runner.Start()
	logger.Info().Msg("Scheduler started!")

	serv := server.NewServer(store)
	logger.Info().Str("port", os.Getenv("PORT")).Msg("Server created. Starting...")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := serv.Start(os.Getenv("PORT")); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info().Msg("Shutting down...")

	// Stop the scheduler
	<-runner.Stop().Done()
	logger.Info().Msg("Scheduler stopped!")

	wg.Wait()

	logger.Info().Msg("Server stopped. Exiting.")
}
