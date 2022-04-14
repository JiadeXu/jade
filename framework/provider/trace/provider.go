package trace

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
)

type JadeTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *JadeTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewJadeTraceService
}

// Boot will called when the service instantiate
func (provider *JadeTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *JadeTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *JadeTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

/// Name define the name for this service
func (provider *JadeTraceProvider) Name() string {
	return contract.TraceKey
}
