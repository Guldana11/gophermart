package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService   service.OrderService
	LoyaltyService service.LoyaltyServiceType
}

func NewOrderHandler(orderSvc service.OrderService, loyaltySvc service.LoyaltyServiceType) *OrderHandler {
	return &OrderHandler{
		orderService:   orderSvc,
		LoyaltyService: loyaltySvc,
	}
}

func (h *OrderHandler) UploadOrderHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if strings.TrimSpace(userID) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !service.CheckLuhn(orderNumber) {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	err = h.orderService.UploadOrder(c.Request.Context(), userID, orderNumber)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrder):
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, service.ErrAlreadyUploadedSelf):
			c.Status(http.StatusOK)
			return
		case errors.Is(err, service.ErrAlreadyUploadedOther):
			c.Status(http.StatusConflict)
			return
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusAccepted)
}

func (h *OrderHandler) GetOrdersHandler(c *gin.Context) {
	userID := c.GetString("userID")

	orders, err := h.orderService.GetOrders(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	c.JSON(http.StatusOK, orders)
}
