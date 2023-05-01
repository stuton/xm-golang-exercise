package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stuton/xm-golang-exercise/utils/jwt"
)

type Header struct {
	BearerToken string `header:"Authorization" binding:"required"`
}

type JwtAuthMiddleware struct {
	token jwt.JWT
}

func NewJwtAuthMiddleware(token jwt.JWT) JwtAuthMiddleware {
	return JwtAuthMiddleware{token: token}
}

func (m JwtAuthMiddleware) Do() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &Header{}

		if err := c.ShouldBindHeader(header); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if err := m.token.TokenValid(header.BearerToken); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
