package main

import (
	"log"
	"os"

	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/handler"
	"x.localhost/rvabot/internal/logger"

	"github.com/joho/godotenv"
)

var repo database.ContentRepositoryInterface

func main() {
	godotenv.Load()

	// Устанавливаем уровень логирования
	logger.SetLevel(logger.INFO)

	botUrl := os.Getenv("TELEGRAM_API") + os.Getenv("TELEGRAM_TOKEN")

	logger.BotInfo("Запуск бота...")
	logger.BotInfo("URL бота: %s", botUrl)

	if err := database.InitDB("rva_bot.db"); err != nil {
		logger.BotError("Ошибка инициализации БД: %v", err)
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	logger.DatabaseInfo("База данных инициализирована успешно")

	repo = database.NewContentRepository()

	logger.BotInfo("Запуск основного цикла бота...")
	handler.BotLoop(botUrl, repo)
}
