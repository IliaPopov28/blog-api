//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

type AuthResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

func setupTestRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)

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

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))
	{
		auth.GET("/posts", postHandler.GetAll)
		auth.GET("/posts/:id", postHandler.GetByID)
		auth.POST("/posts", postHandler.Create)
		auth.PUT("/posts/:id", postHandler.Update)
		auth.DELETE("/posts/:id", postHandler.Delete)
	}

	return r
}

func TestIntegration_API(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Skip("No .env file, skipping integration test")
	}

	cfg, err := config.Init()
	if err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}
	defer cfg.Close()

	router := setupTestRouter(cfg)
	var token string

	t.Run("Health check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Register user", func(t *testing.T) {
		body := `{"username":"testuser","password":"password123"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Login user", func(t *testing.T) {
		body := `{"username":"testuser","password":"password123"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp AuthResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		token = resp.Token
		if token == "" {
			t.Error("Expected token, got empty string")
		}
	})

	if token == "" {
		t.Skip("No token, skipping authenticated tests")
	}

	t.Run("Create post", func(t *testing.T) {
		body := `{"title":"Test Title","content":"Test Content"}`
		req, _ := http.NewRequest("POST", "/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Get all posts", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/posts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}
