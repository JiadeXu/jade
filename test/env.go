package test

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/provider/app"
	"github.com/JiadeXu/jade/framework/provider/env"
)

const (
	BasePath = "/Users/xujiade/go/src/github.com/JiadeXu/jade"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewJadeContainer()
	// 绑定App服务提供者
	container.Bind(&app.JadeAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.JadeTestingEnvProvider{})
	return container
}
