package main

import (
	"log"
	"os"

	"Blog_project_with_Go/internal/database"
	"Blog_project_with_Go/internal/handlers"
	"Blog_project_with_Go/internal/middleware"
	"Blog_project_with_Go/internal/repository"
	"Blog_project_with_Go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env topilmadi, system env ishlatiladi")
	}

	database.Connect()

	postRepo := repository.NewPostRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)

	postService := service.NewPostService(postRepo)
	authService := service.NewAuthService(userRepo)

	postHandler := handlers.NewPostHandler(postService)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPost)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/posts", postHandler.CreatePost)
		protected.PATCH("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}