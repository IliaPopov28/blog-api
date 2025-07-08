package main

import (
	"blog-api/config"
	"blog-api/controllers"
	"blog-api/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.ConnectDatabase()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/")
	auth.Use(middleware.Auth())
	{
		auth.GET("/posts", controllers.GetPosts)
		auth.GET("/posts/:id", controllers.GetPostByID)
		auth.POST("/posts", controllers.CreatePost)
		auth.PUT("/posts/:id", controllers.UpdatePost)
		auth.DELETE("/posts/:id", controllers.DeletePost)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
