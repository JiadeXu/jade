package kernel

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"github.com/JiadeXu/jade/framework/gin"
)

// JadeKernelProvider 提供web引擎
type JadeKernelProvider struct {
	HttpEngine *gin.Engine
}

// REgister 注册服务提供者
func (provider *JadeKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewJadeKernelService
}

// Boot 启动的时候判断是否由外界注入了Engine，如果注入的话用注入，如果没有 重新实例化一个
func (provider *JadeKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}
	provider.HttpEngine.SetContainer(c)
	return nil
}

// IsDefer 引擎的初始化我们希望开始就进行初始化
func (provider *JadeKernelProvider) IsDefer() bool {
	return false
}

// Params 参数就是一个HttpEngine
func (provider *JadeKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

// Name 提供凭证
func (provider *JadeKernelProvider) Name() string {
	return contract.KernelKey
}
