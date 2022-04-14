package app

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
)

type JadeAppProvider struct {
	BaseFolder string
}

// Register 注册
func (j *JadeAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewJadeApp
}

// Boot 启动调用
func (j *JadeAppProvider) Boot(container framework.Container) error {
	return nil
}

// 是否延迟初始化
func (j *JadeAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (j *JadeAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, j.BaseFolder}
}

func (j *JadeAppProvider) Name() string {
	return contract.AppKey
}
