package demo2

import (
	"github.com/JiadeXu/jade/framework/gin"
	"time"
)

func UserLoginController(c *gin.Context) {
	// 打印控制器名字
	foo, _ := c.DefaultQueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second) // 输出结果
	c.ISetOkStatus().IJson("ok, UserLoginController " + foo)
}
