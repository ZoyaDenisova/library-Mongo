// @title       Auth Service API
// @version     1.0
// @description Сервис авторизации с JWT и Swagger UI.
// @host        localhost:8080
// @BasePath    /
package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "library-Mongo/cmd/app/docs"
	"library-Mongo/internal/config"
	"library-Mongo/internal/handler"
	"library-Mongo/internal/repo/mongo"
	"library-Mongo/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//swag init -g ./cmd/app/main.go --parseInternal --output ./cmd/app/docs
//генерирует сваггер (из корня проекта)

func main() {
	// Контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Завершение по Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Завершение приложения по сигналу")
		cancel()
	}()

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к Mongo
	db, err := mongo.Connect(ctx, cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к Mongo:", err)
	}

	// Инициализация репозиториев
	userRepo := mongo.NewUserRepo(db)
	bookRepo := mongo.NewBookRepo(db)
	borrowRepo := mongo.NewBorrowRepo(db)

	// Инициализация usecase
	BorrowUC := usecase.NewBorrowUsecase(borrowRepo, bookRepo, userRepo)
	BookUC := usecase.NewBookUsecase(bookRepo)
	UserUC := usecase.NewUserUsecase(userRepo)

	// Инициализация хендлеров
	borrowHandler := handler.NewBorrowHandler(BorrowUC)
	bookHandler := handler.NewBookHandler(BookUC)
	userHandler := handler.NewUserHandler(UserUC)

	// HTTP сервер на Gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Регистрация Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/borrow/history/:userID", borrowHandler.GetBorrowHistory)
	r.POST("/borrow", borrowHandler.BorrowBook)
	r.POST("/borrow/return", borrowHandler.ReturnBook)
	r.GET("/borrow/overdue", borrowHandler.GetOverdueBorrows)
	r.GET("/borrow/stats", borrowHandler.GetDailyBorrowStats)
	r.GET("/borrow/active-count", borrowHandler.CountActiveBorrows)

	r.POST("/books", bookHandler.CreateBook)
	r.PUT("/books", bookHandler.UpdateBook)
	r.GET("/books/search", bookHandler.SearchBooks)
	r.DELETE("/books/:id", bookHandler.DeleteBook)
	r.GET("/books/:id", bookHandler.GetBookByID)
	r.GET("/books/count", bookHandler.CountBooks)

	r.POST("/users/login", userHandler.Login)
	r.GET("/users/search", userHandler.SearchUsers)
	r.PUT("/users", userHandler.UpdateUser)
	r.GET("/users/:id", userHandler.GetUserByID)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	// Запуск HTTP-сервера
	go func() {
		log.Println("HTTP сервер запущен на порту:", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание завершения
	<-ctx.Done()
	log.Println("Контекст завершён, выключаем HTTP сервер...")

	// Таймаут на завершение сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Приложение завершено корректно")
}
