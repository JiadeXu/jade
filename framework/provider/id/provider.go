package id

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
)

type JadeIDProvider struct {
}

// Register registe a new function for make a service instance
func (provider *JadeIDProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeIDService
}

// Boot will called when the service instantiate
func (provider *JadeIDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *JadeIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *JadeIDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this service
func (provider *JadeIDProvider) Name() string {
	return contract.IDKey
}
