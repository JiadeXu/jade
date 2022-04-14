package middleware

import (
	"github.com/JiadeXu/jade/framework/gin"
	"context"
	"fmt"
	"log"
	"time"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 执行业务逻辑前操作: 初始化超时 context
		durationContext, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 执行具体业务逻辑
			c.Next()

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			c.ISetStatus(500).IJson("time out")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationContext.Done():
			c.ISetStatus(500).IJson("time out")
		}
	}
}
