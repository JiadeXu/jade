package config

import (
	"a-projects/geekbang/framework"
	"a-projects/geekbang/framework/contract"
	"path/filepath"
)

type JadeConfigProvider struct {
	framework.ServiceProvider
}

// Register registe a new function for mak a service instance
func (provider *JadeConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewJadeConfig
}

// Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这个里面。

// Boot will called when the service instantiate
func (provider *JadeConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *JadeConfigProvider) IsDefer() bool {
	return false
}

// Params params 定义传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数
func (provider *JadeConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	// 配置文件夹地址
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

// Name 代表了这个服务提供者的凭证
func (provider *JadeConfigProvider) Name() string {
	return contract.ConfigKey
}
