package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dCatherinee/plant-care-bot/internal/transport/telegram"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.Default()

	if err := godotenv.Load(); err != nil {
		logger.Info("skip .env loading", "err", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	tgBot, err := telegram.New(token, logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := tgBot.Run(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
