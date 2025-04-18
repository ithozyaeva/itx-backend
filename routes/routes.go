package routes

import (
	"ithozyeva/internal/handler"
	"ithozyeva/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	SetupPublicRoutes(app, db)
	SetupPrivateRoutes(app, db)
}
func SetupPublicRoutes(app *fiber.App, db *gorm.DB) {
	// Инициализация сервисов и репозиториев
	telegramAuthHandler := handler.NewTelegramAuthHandler()

	api := app.Group("/api")

	// Маршруты для авторизации через Telegram
	auth := api.Group("/auth")
	auth.Post("/telegram", telegramAuthHandler.Authenticate)
	auth.Post("/telegram-from-bot", telegramAuthHandler.HandleBotMessage)

	userHandler := handler.NewUserHandler()
	// Маршруты для аутентификации в админ панели
	// TODO: Рассмотреть вариант авторизации в админке через тг + роль admin у member (подумать о его переименование в user, или же оставить тотальное разделение между public и admin зонами)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)

	mentorHandler := handler.NewMentorHandler()
	api.Get("/mentors", mentorHandler.GetAllWithRelations)

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	api.Get("/profTags", profTagHandler.Search)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	api.Get("/members", memberHandler.Search)

	// Маршруты для отзывов на услуги
	reviewOnServiceHandler := handler.NewReviewOnServiceHandler()
	api.Get("/reviews-on-service", reviewOnServiceHandler.GetReviewsWithMentorInfo)

	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	api.Get("/reviews", reviewHandler.GetAllWithAuthor)
}

func SetupPrivateRoutes(app *fiber.App, db *gorm.DB) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Защищенные маршруты
	protected := app.Group("/api", authMiddleware.RequireAuth)

	// Маршруты для менторов
	mentorHandler := handler.NewMentorHandler()
	mentors := protected.Group("/mentors")
	mentors.Get("/:id", mentorHandler.GetById)
	mentors.Post("/", mentorHandler.Create)
	mentors.Put("/:id", mentorHandler.Update)
	mentors.Delete("/:id", mentorHandler.Delete)
	mentors.Post("/findByTag", mentorHandler.FindByTag)
	mentors.Post("/review", mentorHandler.AddReviewToService)
	mentors.Get("/:id/services", mentorHandler.GetServices)

	// Здесь будут защищенные маршруты

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	profTags := protected.Group("/profTags")
	profTags.Get("/:id", profTagHandler.GetById)
	profTags.Post("/", profTagHandler.Create)
	profTags.Put("/", profTagHandler.Update)
	profTags.Delete("/:id", profTagHandler.Delete)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	members := protected.Group("/members")
	members.Get("/me", memberHandler.Me)
	members.Patch("/me", memberHandler.Update)
	members.Post("/", memberHandler.Create)
	members.Get("/:id", memberHandler.GetById)
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
	reviewsOnService.Get("/:id", reviewOnServiceHandler.GetById)
	reviewsOnService.Post("/", reviewOnServiceHandler.CreateReview)
	reviewsOnService.Put("/:id", reviewOnServiceHandler.Update)
	reviewsOnService.Delete("/:id", reviewOnServiceHandler.Delete)
}
