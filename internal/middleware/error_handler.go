package middleware

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			_ = c.Errors.Last()
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
		}
	}
}
