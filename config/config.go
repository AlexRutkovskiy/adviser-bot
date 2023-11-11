package config

import (
	"flag"
	"log"
)

type Config struct {
	TgBotToken string
}

func MustLoad() Config {
	tgFlag := flag.String(
		"tg-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *tgFlag == "" {
		log.Fatal("token is not specified")
	}

	return Config{
		TgBotToken: *tgFlag,
	}
}
