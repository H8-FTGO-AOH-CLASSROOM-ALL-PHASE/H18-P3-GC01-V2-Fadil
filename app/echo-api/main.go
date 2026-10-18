package main

import (
	"context"
	"log"
	"os"
	"time"

	"H18-P3-GC01-V2-Fadil/database"
	"H18-P3-GC01-V2-Fadil/handler"
	appmw "H18-P3-GC01-V2-Fadil/middleware"
	"H18-P3-GC01-V2-Fadil/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	database.LoadEnv()

	cfg := &database.Config{DBMongoName: getenv("MONGODB_NAME", "shopping")}

	userRepo := repository.NewUserRepository(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := userRepo.EnsureIndexes(ctx); err != nil {
		log.Fatalf("failed to create user indexes: %v", err)
	}

	authHandler := handler.NewAuthHandler(userRepo)
	productHandler := handler.NewProductHandler(repository.NewProductRepository(cfg))
	transactionHandler := handler.NewTransactionHandler(repository.NewTransactionRepository(cfg))

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	api := e.Group("/api")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	protected := api.Group("", appmw.JWT(func(ctx context.Context, userID int64) (string, error) {
		u, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			return "", err
		}
		return u.JWTSecret, nil
	}))

	protected.POST("/products", productHandler.Create)
	protected.GET("/products", productHandler.GetAll)
	protected.GET("/products/:id", productHandler.GetByID)
	protected.PUT("/products/:id", productHandler.Update)
	protected.DELETE("/products/:id", productHandler.Delete)

	protected.POST("/transactions", transactionHandler.Create)
	protected.GET("/transactions", transactionHandler.GetAll)
	protected.GET("/transactions/:id", transactionHandler.GetByID)
	protected.PUT("/transactions/:id", transactionHandler.Update)
	protected.DELETE("/transactions/:id", transactionHandler.Delete)

	e.Logger.Fatal(e.Start(":" + getenv("PORT", "8080")))
}
