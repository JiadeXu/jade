package ssh

import (
	"context"
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"golang.org/x/crypto/ssh"
	"sync"
)

type JadeSSH struct {
	contract.SSHService

	container framework.Container
	clients map[string]*ssh.Client

	lock *sync.RWMutex
}

// NewJadeSSH 代表实例化Client
func NewJadeSSH(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*ssh.Client)
	lock := &sync.RWMutex{}
	return &JadeSSH{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (app *JadeSSH) GetClient(option ...contract.SSHOption) (*ssh.Client, error) {
	logService := app.container.MustMake(contract.LogKey).(contract.Log)
	// 读取默认配置
	config := GetBaseConfig(app.container)

	//
	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	key := config.UniqKey()

	app.lock.RLock()
	if db, ok := app.clients[key]; ok {
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	// 没有实例化过
	app.lock.Lock()
	defer app.lock.Unlock()

	// 实例化
	addr := config.Host + ":" + config.Port

	client, err := ssh.Dial(config.NetWork, addr, config.ClientConfig)
	if err != nil {
		logService.Error(context.Background(), "ssh dial error", map[string]interface{}{
			"err":  err,
			"addr": addr,
		})
	}

	// 挂载
	app.clients[key] = client
	return client, nil
}