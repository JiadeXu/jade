package env

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
)

type JadeEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *JadeEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewJadeEnv
}

// Boot will called when the service instantiate
func (provider *JadeEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *JadeEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *JadeEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

/// Name define the name for this service
func (provider *JadeEnvProvider) Name() string {
	return contract.EnvKey
}
