package framework

import (
	"errors"
	"fmt"
	"sync"
)

type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，不返回 error
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会 panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

type JadeContainer struct {
	Container // 强制要求 JadeContainer 实现Container接口
	// providers 存储注册的服务提供者, key 为字符串凭证
	providers map[string]ServiceProvider
	// instances 存储具体的实例, key为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

// NewJadeContainer 创建一个服务容器
func NewJadeContainer() *JadeContainer {
	return &JadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (jade *JadeContainer) PrintProviders() []string {
	ret := []string{}
	for _, provider := range jade.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind 将服务容器和关键字做了绑定
func (jade *JadeContainer) Bind(provider ServiceProvider) error {
	jade.lock.Lock()
	key := provider.Name()

	jade.providers[key] = provider
	jade.lock.Unlock()

	// if provider is not defer
	if !provider.IsDefer() {
		if err := provider.Boot(jade); err != nil {
			return err
		}
		// 实例化方法
		params := provider.Params(jade)
		method := provider.Register(jade)
		instance, err := method(params...)
		if err != nil {
			return err
		}
		jade.instances[key] = instance
	}
	return nil
}

func (jade *JadeContainer) IsBind(key string) bool {
	return jade.findServiceProvider(key) != nil
}

func (jade *JadeContainer) findServiceProvider(key string) ServiceProvider {
	jade.lock.RLock()
	defer jade.lock.RUnlock()
	if sp, ok := jade.providers[key]; ok {
		return sp
	}
	return nil
}

func (jade *JadeContainer) Make(key string) (interface{}, error) {
	return jade.make(key, nil, false)
}

func (jade *JadeContainer) MustMake(key string) interface{} {
	serv, err := jade.make(key, nil, true)
	if err != nil {
		panic(err)
	}
	return serv
}

func (jade *JadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return jade.make(key, params, true)
}

func (jade *JadeContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	// force new a
	if err := sp.Boot(jade); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(jade)
	}
	method := sp.Register(jade)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// 真正实例化一个服务
func (jade *JadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	jade.lock.RLock()
	defer jade.lock.RUnlock()
	// 查询是否已注册了这个服务提供者，如果没有注册，则返回错误
	sp := jade.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return jade.newInstance(sp, params)
	}

	// 不需要强制创新实例化，如果容器中已经实例化了，那么久直接使用容器中的实例
	if ins, ok := jade.instances[key]; ok {
		return ins, nil
	}

	// 容器还未实例化，则进行一次实例化
	inst, err := jade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	jade.instances[key] = inst
	return inst, nil
}

func (jade *JadeContainer) NameList() []string {
	ret := make([]string, 0)
	for _, provider := range jade.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}