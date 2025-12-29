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

	err := h.BalanceService.Withdraw(
		c.Request.Context(),
		userID,
		req.Order,
		req.Sum,
	)
	if err != nil {
		switch err {
		case service.ErrInvalidOrder:
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		case service.ErrInsufficientFunds:
			c.AbortWithStatus(http.StatusPaymentRequired) // 402
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *UserHandler) GetWithdrawals(c *gin.Context) {
	userID := c.GetString("userID")
	if strings.TrimSpace(userID) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.BalanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		log.Printf("GetWithdrawals failed for user %s: %v", userID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if withdrawals == nil {
		withdrawals = []models.Withdrawal{}
	}

	c.JSON(http.StatusOK, withdrawals)
}
