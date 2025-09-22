package main

import (
	"log"
	"os"

	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/handler"

	"github.com/joho/godotenv"
)

var repo database.ContentRepositoryInterface

func main() {
	godotenv.Load()

	botUrl := os.Getenv("TELEGRAM_API") + os.Getenv("TELEGRAM_TOKEN")

	log.Println("Бот запущен")

	if err := database.InitDB("rva_bot.db"); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	repo = database.NewContentRepository()

	handler.BotLoop(botUrl, repo)

}
