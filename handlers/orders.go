package handlers

import (
	"net/http"
	"strings"

	"github.com/Guldana11/gophermart/database"
	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

func UploadOrderHandler(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(string) // строка UUID

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

	for _, ch := range orderNumber {
		if ch < '0' || ch > '9' {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
	}

	if !service.CheckLuhn(orderNumber) {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx := c.Request.Context()

	existingUserID, exists, err := database.CheckOrderExists(ctx, orderNumber)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if exists {
		if existingUserID == userID {
			c.Status(http.StatusOK)
			return
		}
		c.Status(http.StatusConflict)
		return
	}

	order := models.Order{
		UserID: userID,
		Number: orderNumber,
	}

	if err := database.CreateOrder(ctx, order); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func GetOrdersHandler(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(string)

	orders, err := database.GetOrdersByUser(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}
