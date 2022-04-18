package env

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
)

type JadeTestingEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *JadeTestingEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewJadeTestingEnv
}

// Boot will called when the service instantiate
func (provider *JadeTestingEnvProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *JadeTestingEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *JadeTestingEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this service
func (provider *JadeTestingEnvProvider) Name() string {
	return contract.EnvKey
}
