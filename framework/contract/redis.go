package contract

import (
	"fmt"
	"github.com/JiadeXu/jade/framework"
	"github.com/go-redis/redis/v8"
)

const RedisKey = "jade:redis.yaml"

// RedisOption 代表初始化的时候的选项
type RedisOption func(container framework.Container, config *RedisConfig) error

type RedisService interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}

type RedisConfig struct {
	*redis.Options
}

// redis的唯一标识
func (config *RedisConfig) UniqueKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}