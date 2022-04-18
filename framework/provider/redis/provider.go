package redis

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
)

// RedisProvider 提供App的具体实现方法
type RedisProvider struct {
}

// Register 注册方法
func (j *RedisProvider) Register(container framework.Container) framework.NewInstance {
	return NewJadeRedis
}

// Boot 启动调用
func (j *RedisProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (j *RedisProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (j *RedisProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (j *RedisProvider) Name() string {
	return contract.RedisKey
}
