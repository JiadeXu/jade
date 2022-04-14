package http

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	// 设置为 release 为的是默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	// 默认启动一个web引擎
	r := gin.New()
	// 设置了Engine
	r.SetContainer(container)

	// 默认注册recovery中间件
	r.Use(gin.Recovery())

	// 业务绑定路由操作
	Routes(r)
	return r, nil
}
