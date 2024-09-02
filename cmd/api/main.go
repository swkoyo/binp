package main

import (
	"binp/scheduler"
	"binp/server"
	"binp/storage"
	"binp/util"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	util.InitLogger()
	logger := util.GetLogger()

	store, err := storage.NewStore()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create new store")
	}

	if err := store.Init(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initalize store")
	}

	runner := scheduler.NewScheduler()
	runner.AddFunc("@hourly", func() {
		logger.Info().Msg("Checking for expired snippets...")
		count, err := store.DeleteExpiredSnippets()
		if err != nil {
			logger.Error().Err(err).Int("count", count).Msg("Failed to delete expired snippets")
		}
		logger.Info().Int("count", count).Msg("Expired snippets deleted")
	})
	logger.Info().Msg("Starting scheduler...")
	runner.Start()
	logger.Info().Msg("Scheduler started!")

	serv := server.NewServer(store)
	logger.Info().Str("port", os.Getenv("PORT")).Msg("Server created. Starting...")
	if err := serv.Start(os.Getenv("PORT")); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
