package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dCatherinee/plant-care-bot/internal/app"
	"github.com/dCatherinee/plant-care-bot/internal/config"
	"github.com/dCatherinee/plant-care-bot/internal/storage/postgres"
	"github.com/dCatherinee/plant-care-bot/internal/transport/telegram"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.Default()

	if err := godotenv.Load(); err != nil {
		logger.Info("skip .env loading", "err", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("close db", "err", err)
		}
	}()

	application := app.New(db)
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	tgBot, err := telegram.New(token, logger, application.PlantService, application.UserService, application.CareEventService)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := tgBot.Run(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
