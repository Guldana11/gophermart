package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(svc service.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type"})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1MB

		var req models.RegisterRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
			return
		}

		if req.Login == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "login and password cannot be empty"})
			return
		}

		user, err := svc.Register(c.Request.Context(), req.Login, req.Password)
		if err != nil {
			if err.Error() == "login already exists" {
				c.JSON(http.StatusConflict, gin.H{"error": "login already taken"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		c.SetCookie("session_id", user.ID, 3600*24, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"id": user.ID, "login": user.Login})
	}
}

func LoginHandler(svc service.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		user, err := svc.Login(c.Request.Context(), req.Login, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
			return
		}

		c.SetCookie("session_id", user.ID, 3600*24, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"id": user.ID, "login": user.Login})
	}
}
