package main

import (
	"context"
	"flag"
	tgClient "github.com/nanmenkaimak/read-advisor-bot/clients/telegram"
	"github.com/nanmenkaimak/read-advisor-bot/consumer/event-consumer"
	"github.com/nanmenkaimak/read-advisor-bot/events/telegram"
	"github.com/nanmenkaimak/read-advisor-bot/storage/sqlite"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "./storage.db"
	batchSize         = 100
)

func main() {
	// s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
