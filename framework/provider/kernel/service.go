package kernel

import (
	"github.com/JiadeXu/jade/framework/gin"
	"net/http"
)

type JadeKernelService struct {
	engine *gin.Engine
}

// 初始化
func NewJadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &JadeKernelService{engine: httpEngine}, nil
}

// 返回 web 引擎
func (j *JadeKernelService) HttpEngine() http.Handler {
	return j.engine
}
