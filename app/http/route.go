package http

import (
	"github.com/JiadeXu/jade/app/http/module/demo"
	"github.com/JiadeXu/jade/framework/gin"
	"github.com/JiadeXu/jade/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	// 动态路由定义
	demo.Register(r)
}
