package middleware

import (
	"a-projects/geekbang/framework/contract"
	"a-projects/geekbang/framework/gin"
)

// recovery 机制将协程中的函数异常进行捕获
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始
		tracer := c.MustMake(contract.TraceKey).(contract.Trace)
		tracerCtx := tracer.ExtractHTTP(c.Request)

		tracer.WithTrace(c, tracerCtx)

		//
		c.Next()
	}
}
