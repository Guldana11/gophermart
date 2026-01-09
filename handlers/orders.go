package handlers

import (
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService   service.OrderService
	loyaltyService service.LoyaltyService
}

func NewOrderHandler(orderSvc service.OrderService, loyaltySvc service.LoyaltyService) *OrderHandler {
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

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.After(orders[j].UploadedAt)
	})

	var result []map[string]interface{}
	for _, order := range orders {
		status := "NEW"
		accrual := 0.0

		if h.loyaltyService != nil {
			resp, err := h.loyaltyService.GetOrderAccrual(c.Request.Context(), order.Number)
			if err == nil && resp != nil {
				status = resp.Status
				accrual = resp.Accrual
			}
		}

		orderMap := map[string]interface{}{
			"number":      order.Number,
			"status":      status,
			"uploaded_at": order.UploadedAt.Format(time.RFC3339),
		}
		if accrual > 0 {
			orderMap["accrual"] = accrual
		}
		result = append(result, orderMap)
	}

	c.JSON(http.StatusOK, result)
}
