package main

import (
	serv "binp/server"
	"binp/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := serv.NewServer(store)
	server.Start(os.Getenv("PORT"))
}
