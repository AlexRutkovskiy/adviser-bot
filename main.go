package main

import (
	tgClient "AdviserBot/clients/telegram"
	"AdviserBot/config"
	event_consumer "AdviserBot/consumer/event-consumer"
	"AdviserBot/events/telegram"
	"AdviserBot/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "file_storage"
	batchSize   = 100
)

func main() {
	cfg := config.MustLoad()

	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken),
		files.New(storagePath),
	)

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}
