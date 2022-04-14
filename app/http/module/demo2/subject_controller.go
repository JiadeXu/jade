package demo2

import (
	"a-projects/geekbang/framework/gin"
	demo2 "a-projects/geekbang/provider/demo"
	"fmt"
)

func SubjectAddController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectAddController")
}

func SubjectListController(c *gin.Context) {
	// 获取demo服务实例
	demoService := c.MustMake(demo2.Key).(demo2.Service)

	foo := demoService.GetFoo()

	c.ISetOkStatus().IJson(foo)

	//c.ISetOkStatus().IJson("ok, SubjectListController")
}

func SubjectDelController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectDelController")
}

func SubjectUpdateController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectUpdateController")
}

func SubjectGetController(c *gin.Context) {
	subjectId, _ := c.DefaultParamInt("id", 0)
	c.ISetOkStatus().IJson("ok, SubjectGetController:" + fmt.Sprint(subjectId))
}

func SubjectNameController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectNameController")
}
