package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stuton/xm-golang-exercise/utils"
)

type Header struct {
	BearerToken string `header:"Authorization" binding:"required"`
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &Header{}

		if err := c.ShouldBindHeader(header); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if err := utils.TokenValid(header.BearerToken); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
