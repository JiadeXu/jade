package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"strings"
	"sync"
	"time"
)

type MemoryData struct {
	val        interface{}
	createTime time.Time
	ttl        time.Duration
}

type MemoryCache struct {
	contract.CacheService
	container framework.Container
	data      map[string]*MemoryData
	lock      sync.RWMutex
}

func NewMemoryCache(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &MemoryCache{
		container: container,
		data:      map[string]*MemoryData{},
		lock:      sync.RWMutex{},
	}, nil
}

func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := m.GetObj(ctx, key, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (m *MemoryCache) GetObj(ctx context.Context, key string, obj interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if md, ok := m.data[key]; ok {
		if md.ttl != NoneDuration {
			if time.Now().Sub(md.createTime) > md.ttl {
				delete(m.data, key)
				return ErrKeyNotFound
			}
		}

		bt, _ := json.Marshal(md.val)
		err := json.Unmarshal(bt, obj)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrKeyNotFound
}

// GetMany 获取某些key对应的值
func (m *MemoryCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)
	for _, key := range keys {
		val, err := m.Get(ctx, key)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rets[key] = val
	}
	if len(errs) == 0 {
		return rets, nil
	}
	return rets, errors.New(strings.Join(errs, "||"))
}

// Set 设置某个key和值到缓存，带超时时间
func (m *MemoryCache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return m.SetObj(ctx, key, val, timeout)
}

// SetObj 设置某个key和对象到缓存, 对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (m *MemoryCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	md := &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        timeout,
	}
	m.data[key] = md
	return nil
}

// SetMany 设置多个key和值到缓存
func (m *MemoryCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	errs := []string{}
	for k, v := range data {
		err := m.Set(ctx, k, v, timeout)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "||"))
	}
	return nil
}

// SetForever 设置某个key和值到缓存，不带超时时间
func (m *MemoryCache) SetForever(ctx context.Context, key string, val string) error {
	return m.Set(ctx, key, val, NoneDuration)
}

// SetForeverObj 设置某个key和对象到缓存，不带超时时间，对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (m *MemoryCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return m.SetObj(ctx, key, val, NoneDuration)
}

func (m *MemoryCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if md, ok := m.data[key]; ok {
		md.ttl = timeout
		return nil
	}
	return ErrKeyNotFound
}

func (m *MemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if md, ok := m.data[key]; ok {
		return md.ttl, nil
	}
	return NoneDuration, ErrKeyNotFound
}

// Remember 实现缓存的Cache-Aside模式, 先去缓存中根据key获取对象，如果有的话，返回，如果没有，调用RememberFunc 生成
func (m *MemoryCache) Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc contract.RememberFunc, obj interface{}) error {
	err := m.GetObj(ctx, key, obj)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	// key not found
	objNew, err := rememberFunc(ctx, m.container)
	if err != nil {
		return err
	}

	if err := m.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}

	if err := m.GetObj(ctx, key, &obj); err != nil {
		return err
	}
	return nil
}

// Calc 往key对应的值中增加step计数
func (m *MemoryCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	var val int64
	err := m.GetObj(ctx, key, &val)
	val = val + step
	if err == nil {
		m.data[key].val = val
		return val, nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return 0, err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        NoneDuration,
	}

	return val, nil
}

func (m *MemoryCache) Increment(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, 1)
}

func (m *MemoryCache) Decrement(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, -1)
}

// Del 删除某个key
func (m *MemoryCache) Del(ctx context.Context, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.data, key)
	return nil
}

// DelMany 删除某些key
func (m *MemoryCache) DelMany(ctx context.Context, keys []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}
