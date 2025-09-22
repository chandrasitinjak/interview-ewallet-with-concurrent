package handler

import (
	"net/http"

	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Credit(c *gin.Context) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid request"})
		return
	}

	tx, newBalance, err := h.service.Credit(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"transaction_id": tx.ID,
		"new_balance":    newBalance,
	})
}

func (h *TransactionHandler) Debit(c *gin.Context) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid request"})
		return
	}

	tx, newBalance, err := h.service.Debit(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"transaction_id": tx.ID,
		"new_balance":    newBalance,
	})
}
