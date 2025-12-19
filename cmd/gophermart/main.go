package main

import (
	"log"
	"os"

	"github.com/Guldana11/gophermart/database"
	"github.com/Guldana11/gophermart/handlers"
	"github.com/Guldana11/gophermart/middleware"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URI")
	if dbURL == "" {
		log.Fatal("DATABASE_URI is not set")
	}

	if err := database.Migrate(dbURL); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	dbPool := database.InitDB()
	defer dbPool.Close()

	userRepo := database.NewUserRepo(dbPool)
	orderRepo := database.NewOrderRepo(dbPool)

	userSvc := service.NewUserService(userRepo)
	orderSvc := service.NewOrderService(orderRepo)

	orderHandler := handlers.NewOrderHandler(orderSvc)

	r := gin.Default()

	r.POST("/api/user/register", handlers.RegisterHandler(userSvc))
	r.POST("/api/user/login", handlers.LoginHandler(userSvc))

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddlewareJWT())
	{
		auth.POST("/orders", orderHandler.UploadOrderHandler)
		auth.GET("/orders", orderHandler.GetOrdersHandler)
	}

	log.Println("server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
