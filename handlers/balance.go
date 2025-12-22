// handlers/user.go
package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	BalanceService service.BalanceServiceType
}

func NewUserHandler(balanceSvc service.BalanceServiceType) *UserHandler {
	return &UserHandler{BalanceService: balanceSvc}
}

func (h *UserHandler) GetBalance(c *gin.Context) {
	userID := c.GetString("userID")
	if strings.TrimSpace(userID) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := h.BalanceService.GetUserBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	})
}

func (h *UserHandler) Withdraw(c *gin.Context) {
	userID := c.GetString("userID")
	if strings.TrimSpace(userID) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var req models.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	newBalance, err := h.BalanceService.Withdraw(c.Request.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch err {
		case service.ErrInvalidOrder:
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid order"})
		case service.ErrInsufficientFunds:
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "insufficient funds"})
		default:
			log.Printf("Withdraw error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, models.WithdrawResponse{
		Current: newBalance,
	})
}

func (h *UserHandler) GetWithdrawals(c *gin.Context) {
	userID := c.GetString("userID")
	if strings.TrimSpace(userID) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.BalanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
