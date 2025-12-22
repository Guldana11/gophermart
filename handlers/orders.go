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
	loyaltyService service.LoyaltyServiceType
}

func NewOrderHandler(orderSvc service.OrderService, loyaltySvc service.LoyaltyServiceType) *OrderHandler {
	return &OrderHandler{
		orderService:   orderSvc,
		loyaltyService: loyaltySvc,
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
	if strings.TrimSpace(userID) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetOrders(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	var result []map[string]interface{}
	for _, order := range orders {
		var status string
		var accrual float64

		if h.loyaltyService != nil {
			acc, err := h.loyaltyService.GetOrderAccrual(c.Request.Context(), order.Number)
			if err == nil && acc != nil {
				status = acc.Status
				accrual = acc.Accrual
			} else {
				status = "REGISTERED"
			}
		} else {
			status = "REGISTERED"
		}

		orderMap := map[string]interface{}{
			"number": order.Number,
			"status": status,
		}
		if accrual > 0 {
			orderMap["accrual"] = accrual
		}
		result = append(result, orderMap)
	}

	c.JSON(http.StatusOK, result)
}
