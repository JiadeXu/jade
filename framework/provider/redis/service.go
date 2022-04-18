package redis

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"github.com/go-redis/redis/v8"
	"sync"
)

type JadeRedis struct {
	contract.RedisService
	container framework.Container
	clients   map[string]*redis.Client

	lock *sync.RWMutex
}

func NewJadeRedis(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}

	return &JadeRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (app *JadeRedis) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetBaseConfig(app.container)

	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	key := config.UniqueKey()

	app.lock.RLock()
	if db, ok := app.clients[key]; ok {
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	app.lock.Lock()
	defer app.lock.Unlock()

	client := redis.NewClient(config.Options)

	// 挂载
	app.clients[key] = client

	return client, nil
}
