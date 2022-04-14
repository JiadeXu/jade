package middleware

import (
	"github.com/JiadeXu/jade/framework/gin"
	"log"
	"time"
)

func Cost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request.RequestURI, cost.Seconds())
	}
}
