package middleware

import (
	"github.com/JiadeXu/jade/framework/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				c.ISetStatus(500).IJson(p)
			}
		}()
		c.Next()
	}
}
