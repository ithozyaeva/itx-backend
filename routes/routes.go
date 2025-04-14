package routes

import (
	"ithozyeva/internal/handler"
	"ithozyeva/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Инициализация сервисов и репозиториев
	telegramAuthHandler := handler.NewTelegramAuthHandler()
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Маршруты для авторизации через Telegram
	auth := app.Group("/api/auth")
	auth.Post("/telegram", telegramAuthHandler.Authenticate)
	auth.Post("/telegram-from-bot", telegramAuthHandler.HandleBotMessage)

	// Маршруты для аутентификации в админ панели
	// TODO: Рассмотреть вариант авторизации в админке через тг + роль admin у member (подумать о его переименование в user, или же оставить тотальное разделение между public и admin зонами)
	userHandler := handler.NewUserHandler()
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)

	// Маршруты для менторов
	mentorHandler := handler.NewMentorHandler()
	mentors := app.Group("/api/mentors")
	mentors.Get("/", mentorHandler.GetAllWithRelations)
	mentors.Get("/:id", mentorHandler.GetById)
	mentors.Post("/", mentorHandler.Create)
	mentors.Put("/:id", mentorHandler.Update)
	mentors.Delete("/:id", mentorHandler.Delete)
	mentors.Post("/findByTag", mentorHandler.FindByTag)
	mentors.Post("/review", mentorHandler.AddReviewToService)
	mentors.Get("/:id/services", mentorHandler.GetServices)

	// Защищенные маршруты
	protected := app.Group("/api", authMiddleware.RequireAuth)
	// Здесь будут защищенные маршруты

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	profTags := protected.Group("/profTags")
	profTags.Get("/", profTagHandler.Search)
	profTags.Get("/:id", profTagHandler.GetById)
	profTags.Post("/", profTagHandler.Create)
	profTags.Put("/", profTagHandler.Update)
	profTags.Delete("/:id", profTagHandler.Delete)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	members := protected.Group("/members")
	members.Get("/", memberHandler.Search)
	members.Get("/:id", memberHandler.GetById)
	members.Post("/", memberHandler.Create)
	members.Put("/", memberHandler.Update)
	members.Delete("/:id", memberHandler.Delete)

	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	reviews := protected.Group("/reviews")
	reviews.Get("/", reviewHandler.GetAllWithAuthor)
	reviews.Post("/", reviewHandler.AddReview)
	reviews.Get("/:id", reviewHandler.GetById)
	reviews.Put("/:id", reviewHandler.Update)
	reviews.Delete("/:id", reviewHandler.Delete)

	// Маршруты для отзывов на услуги
	reviewOnServiceHandler := handler.NewReviewOnServiceHandler()
	reviewsOnService := protected.Group("/reviews-on-service")
	reviewsOnService.Get("/", reviewOnServiceHandler.GetReviewsWithMentorInfo)
	reviewsOnService.Get("/:id", reviewOnServiceHandler.GetById)
	reviewsOnService.Post("/", reviewOnServiceHandler.CreateReview)
	reviewsOnService.Put("/:id", reviewOnServiceHandler.Update)
	reviewsOnService.Delete("/:id", reviewOnServiceHandler.Delete)
}
