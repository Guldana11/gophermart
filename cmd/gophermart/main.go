package main

import (
	"log"

	"github.com/Guldana11/gophermart/database"
	"github.com/Guldana11/gophermart/handlers"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Init()
	database.Migrate(db)

	userRepo := database.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)

	r := gin.Default()

	r.POST("/api/user/register", handlers.RegisterHandler(userSvc))
	r.POST("/api/user/login", handlers.LoginHandler(userSvc))

	auth := r.Group("/api/user")
	auth.Use(handlers.AuthMiddleware())
	{
		auth.POST("/orders", handlers.UploadOrderHandler)
		auth.GET("/orders", handlers.GetOrdersHandler)
	}

	log.Println("server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
