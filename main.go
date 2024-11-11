package main

import (
	"context"
	"flag"
	"log"
	tgClient "tgBot/clients/telegram"
	eventConsumer "tgBot/consumer/event-consumer"
	"tgBot/events/telegram"
	"tgBot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: %v", err)
	}

	if err = s.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: %v", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(mustTokenAndHost()),
		s,
	)

	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err = consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustTokenAndHost() (string, string) {
	host := flag.String(
		"tg-bot-host",
		"",
		"host for telegram bot",
	)

	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *host == "" {
		log.Fatal("host is required")
	}

	if *token == "" {
		log.Fatal("token is required")
	}

	return *host, *token
}
