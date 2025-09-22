package main

import (
	"log"

	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/config"
	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/handler"
	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/repository"
	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/service"
	"github.com/gin-gonic/gin"
)

func main() {

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	// init repo
	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// init service
	txService := service.NewTransactionService(userRepo, txRepo)

	// init handler
	txHandler := handler.NewTransactionHandler(txService)
	r := gin.Default()

	api := r.Group("/api/transactions")
	{
		api.POST("/credit", txHandler.Credit)
		api.POST("/debit", txHandler.Debit)
	}

	log.Println("server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server:", err)
	}

}
