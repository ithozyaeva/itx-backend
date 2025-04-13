package main

import (
	"fmt"
	"ithozyeva/config"
	"ithozyeva/internal/bot"
	"ithozyeva/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	// Формируем DSN для подключения к базе данных
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.CFG.Database.Host,
		config.CFG.Database.Port,
		config.CFG.Database.User,
		config.CFG.Database.Password,
		config.CFG.Database.Name,
	)

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создаем экземпляр Fiber
	app := fiber.New(fiber.Config{
		AppName: "ITX API",
	})

	// Добавляем middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Telegram-User-ID",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Настраиваем маршруты
	routes.SetupRoutes(app, db)

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
