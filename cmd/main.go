package main

import (
	"ithozyeva/config"
	"ithozyeva/database"
	"ithozyeva/internal/bot"
	"ithozyeva/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	// Подключаемся к базе данных
	database.SetupDatabase()

	// Создаем экземпляр Fiber
	app := fiber.New(fiber.Config{
		AppName: "ITX API",
	})

	// Добавляем middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Telegram-User-Token",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Настраиваем маршруты
	routes.SetupRoutes(app, database.DB)

	// Запускаем Telegram бота в отдельной горутине
	go func() {
		telegramBot, err := bot.NewTelegramBot()
		if err != nil {
			log.Printf("Error creating bot: %v", err)
			return
		}

		log.Println("Telegram bot started successfully")
		telegramBot.Start()
	}()

	// Запускаем сервер
	log.Printf("Server starting on port %s", config.CFG.Port)
	if err := app.Listen(":" + config.CFG.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
