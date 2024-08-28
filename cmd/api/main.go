package main

import (
	"binp/logger"
	"binp/server"
	"binp/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	logger.InitLogger()
	log := logger.GetLogger()

	store, err := storage.NewStore()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new store")
	}

	if err := store.Init(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initalize store")
	}

	serv := server.NewServer(store)
	log.Info().Str("port", os.Getenv("PORT")).Msg("Server created. Starting...")
	if err := serv.Start(os.Getenv("PORT")); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
