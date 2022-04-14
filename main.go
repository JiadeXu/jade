package main

import (
	"a-projects/geekbang/app/console"
	"a-projects/geekbang/app/http"
	"a-projects/geekbang/framework"
	"a-projects/geekbang/framework/provider/app"
	"a-projects/geekbang/framework/provider/config"
	"a-projects/geekbang/framework/provider/distributed"
	"a-projects/geekbang/framework/provider/env"
	"a-projects/geekbang/framework/provider/id"
	"a-projects/geekbang/framework/provider/kernel"
	"a-projects/geekbang/framework/provider/log"
	"a-projects/geekbang/framework/provider/trace"
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

	// 后续初始化需要绑定的服务提供者
	// 将 HTTP 引擎初始化，并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.JadeKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}
