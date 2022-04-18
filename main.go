package main

import (
	"github.com/JiadeXu/jade/app/console"
	"github.com/JiadeXu/jade/app/http"
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/provider/app"
	"github.com/JiadeXu/jade/framework/provider/cache"
	"github.com/JiadeXu/jade/framework/provider/config"
	"github.com/JiadeXu/jade/framework/provider/distributed"
	"github.com/JiadeXu/jade/framework/provider/env"
	"github.com/JiadeXu/jade/framework/provider/id"
	"github.com/JiadeXu/jade/framework/provider/kernel"
	"github.com/JiadeXu/jade/framework/provider/log"
	"github.com/JiadeXu/jade/framework/provider/orm"
	"github.com/JiadeXu/jade/framework/provider/redis"
	"github.com/JiadeXu/jade/framework/provider/trace"
)

func main() {
	// 初始化服务器容器
	container := framework.NewJadeContainer()

	// 绑定 App 服务提供者
	container.Bind(&app.JadeAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.JadeEnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&config.JadeConfigProvider{})
	container.Bind(&id.JadeIDProvider{})
	container.Bind(&trace.JadeTraceProvider{})
	container.Bind(&log.JadeLogServiceProvider{})
	container.Bind(&orm.GormProvider{})

	container.Bind(&redis.RedisProvider{})
	container.Bind(&cache.JadeCacheProvider{})
	// 后续初始化需要绑定的服务提供者
	// 将 HTTP 引擎初始化，并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.JadeKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}
