package handlers

import (
	"net/http"
	"strings"

	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) GetOrderAccrual(c *gin.Context) {
	orderNumber := c.Param("number")
	if strings.TrimSpace(orderNumber) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid order number"})
		return
	}

	resp, err := h.LoyaltyService.GetOrderAccrual(c.Request.Context(), orderNumber)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			c.Status(http.StatusNoContent)
		case service.ErrTooManyReq:
			c.Header("Retry-After", "60")
			c.String(http.StatusTooManyRequests, "No more than N requests per minute allowed")
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
