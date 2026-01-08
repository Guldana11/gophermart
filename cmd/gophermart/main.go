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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	middleware.SetJWTKey([]byte(jwtSecret))

	dbURL := os.Getenv("DATABASE_URI")
	if dbURL == "" {
		log.Fatal("DATABASE_URI is not set")
	}

	accrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if accrualAddr == "" {
		log.Fatal("ACCRUAL_SYSTEM_ADDRESS is empty")
	}

	if err := database.Migrate(dbURL); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	dbPool, err := database.InitDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to init db pool: %v", err)
	}
	defer dbPool.Close()

	userRepo := database.NewUserRepo(dbPool)
	orderRepo := database.NewOrderRepo(dbPool)

	userSvc := service.NewUserService(userRepo)
	orderSvc := service.NewOrderService(orderRepo)
	loyaltySvc := service.NewLoyaltyService(accrualAddr)
	balanceSvc := service.NewBalanceService(userRepo)

	orderHandler := handlers.NewOrderHandler(orderSvc, loyaltySvc)
	userHandler := handlers.NewUserHandler(balanceSvc)

	r := gin.Default()
	r.POST("/api/user/register", handlers.RegisterHandler(userSvc))
	r.POST("/api/user/login", handlers.LoginHandler(userSvc))

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddlewareJWT())
	{
		auth.POST("/user/orders", orderHandler.UploadOrderHandler)
		auth.GET("/user/orders", orderHandler.GetOrdersHandler)

		auth.GET("/user/balance", userHandler.GetBalance)
		auth.POST("/user/balance/withdraw", userHandler.Withdraw)
		auth.GET("/user/withdrawals", userHandler.GetWithdrawals)
	}

	log.Println("server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
