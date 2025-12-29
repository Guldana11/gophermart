package handlers

import (
	"net/http"
	"strings"

	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) GetOrderAccrual(c *gin.Context) {
	orderNumber := strings.TrimSpace(c.Param("number"))
	if orderNumber == "" {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	resp, err := h.loyaltyService.GetOrderAccrual(c.Request.Context(), orderNumber)
	if err != nil {
		if err == service.ErrTooManyReq {
			c.Header("Retry-After", "60")
			c.Status(http.StatusTooManyRequests)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}
