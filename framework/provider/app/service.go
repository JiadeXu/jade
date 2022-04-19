package app

import (
	"errors"
	"flag"
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"github.com/JiadeXu/jade/framework/util"
	"github.com/google/uuid"
	"path/filepath"
)

var baseFolder = ""

func init() {
	flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
	flag.Parse()
}

type JadeApp struct {
	contract.App
	container framework.Container // 服务容器

	baseFolder string            // 基础路径
	appId      string            // 表示当前这个app的唯一id, 可以用于分布式锁等
	configMap  map[string]string // 配置加载
}

func (j JadeApp) Version() string {
	return "0.0.3"
}

func (j JadeApp) BaseFolder() string {
	if j.baseFolder != "" {
		return j.baseFolder
	}
	// 如果参数也没有
	return util.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (j JadeApp) ConfigFolder() string {
	if val, ok := j.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (j JadeApp) LogFolder() string {
	if val, ok := j.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(j.StorageFolder(), "log")
}

func (j JadeApp) HttpFolder() string {
	if val, ok := j.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "app", "http")
}

func (j JadeApp) ConsoleFolder() string {
	if val, ok := j.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "app", "console")
}

func (j JadeApp) StorageFolder() string {
	if val, ok := j.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (j JadeApp) ProviderFolder() string {
	if val, ok := j.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "app", "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (j JadeApp) MiddlewareFolder() string {
	if val, ok := j.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(j.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (j JadeApp) CommandFolder() string {
	if val, ok := j.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(j.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (j JadeApp) RuntimeFolder() string {
	if val, ok := j.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(j.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (j JadeApp) TestFolder() string {
	if val, ok := j.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "test")
}

// DeployFolder 定义测试需要的信息
func (app JadeApp) DeployFolder() string {
	if val, ok := app.configMap["deploy_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "deploy")
}

// NewJadeApp 初始化 JadeApp
func NewJadeApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	// 有两个参数 一个是容器 一个是 baseFolder
	container := params[0].(framework.Container)
	baseFolderTmp := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolderTmp == "" {
		baseFolderTmp = baseFolder
	}
	appId := uuid.New().String()
	configMap := map[string]string{}
	return &JadeApp{baseFolder: baseFolderTmp, container: container, appId: appId, configMap: configMap}, nil
}

func (j JadeApp) AppID() string {
	return j.appId
}

func (j *JadeApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		j.configMap[key] = val
	}
}

func (j *JadeApp) AppFolder() string {
	if val, ok := j.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(j.BaseFolder(), "app")
}
