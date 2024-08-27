package main

import (
	"binp/scheduler"
	"binp/storage"
	"context"
	"log"
	"os"
	"os/signal"
	"time"

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

	runner := scheduler.NewScheduler()

	runner.AddFunc("@every 1m", func() {
		err := store.DeleteExpiredSnippets()
		if err != nil {
			log.Println(err)
		}
	})

	runner.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	<-runner.Stop().Done()
}
