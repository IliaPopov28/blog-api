package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog-api/internal/config"
	"blog-api/internal/handlers"
	"blog-api/internal/middleware"
	"blog-api/internal/repository"
	"blog-api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	if err := config.RunMigrations(cfg.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(cfg.DB)
	postRepo := repository.NewPostRepository(cfg.DB)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	postService := service.NewPostService(postRepo)

	authHandler := handlers.NewAuthHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.RateLimiter())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))
	{
		auth.GET("/me", middleware.Timeout(5*time.Second), authHandler.GetMe)
		auth.GET("/posts", middleware.Timeout(5*time.Second), postHandler.GetAll)
		auth.GET("/posts/:id", middleware.Timeout(5*time.Second), postHandler.GetByID)
		auth.POST("/posts", middleware.Timeout(5*time.Second), postHandler.Create)
		auth.PUT("/posts/:id", middleware.Timeout(5*time.Second), postHandler.Update)
		auth.DELETE("/posts/:id", middleware.Timeout(5*time.Second), postHandler.Delete)
	}

	srv := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Server starting", slog.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	if err := cfg.Close(); err != nil {
		slog.Error("Database connection close error", slog.String("error", err.Error()))
	}

	slog.Info("Server exited")
}
