package demo2

import (
	"github.com/JiadeXu/jade/framework/gin"
	"github.com/JiadeXu/jade/framework/middleware"
)

func registerRouter(core *gin.Engine) {

	// core中使用use注册中间件
	core.Use(
		middleware.Test1(),
		middleware.Test2())
	// 需求1+2:HTTP方法+静态路由匹配
	core.GET("/user/login", UserLoginController)

	// 需求3 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Test3())
		// 需求4:动态路由
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)
		t2 := subjectApi.Group("/pp")
		{
			t2.GET("/p1", func(c *gin.Context) {
				c.ISetOkStatus().IJson("233456")
			})
		}
	}
}
