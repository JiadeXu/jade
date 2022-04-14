package distributed

import (
	"a-projects/geekbang/framework"
	"a-projects/geekbang/framework/contract"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type LocalDistributedService struct {
	container framework.Container
}

// NewLocalDistributedService 初始化本地分布式服务
func NewLocalDistributedService(params ...interface{}) (interface{}, error) {
	if len(params) == 0 {
		return nil, errors.New("param error")
	}

	// 两个参数 容器 baseFolder
	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container}, nil
}

// Select 为分布式选择器
func (s LocalDistributedService) Select(serviceName string, appID string, holdTime time.Duration) (selectAppID string, err error) {
	appService := s.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)

	// 打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// 尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	// 抢不到
	if err != nil {
		// 读取被选择的appid
		selectAppIDByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt), err
	}

	// 在一段时间内选举有效 其他节点这段时间不能再抢占
	go func() {
		defer func() {
			// 释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 释放文件
			lock.Close()
			// 删除文件锁对应的文件
			os.Remove(lockFile)
		}()
		timer := time.NewTimer(holdTime)
		// 等待计时器结束
		<-timer.C
	}()

	// 这里已经是抢占了
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}
	return appID, nil
}
