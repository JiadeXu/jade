package http

import "github.com/JiadeXu/jade/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	// 设置为 release 为的是默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	// 默认启动一个web引擎
	r := gin.Default()
	// 业务绑定路由操作
	Routes(r)
	return r, nil
}
