package config

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type JadeConfig struct {
	contract.Config
	c        framework.Container
	folder   string // 文件夹
	keyDelim string // 路径分隔符 默认为点
	lock     sync.RWMutex
	envMaps  map[string]string
	confMaps map[string]interface{} // 配置文件结构
	confRows map[string][]byte      // 配置文件的原始信息
}

// 读取某个配置文件
func (conf *JadeConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	// 判断文件是否已yaml或者yml作为后缀
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yml" || s[1] == "yaml") {
		name := s[0]

		// 读取文件内容
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		// 直接针对文本做环境变量的替换
		bf = replace(bf, conf.envMaps)
		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		conf.confMaps[name] = c
		conf.confRows[name] = bf

		// 读取app.path中的信息 更新app对应的folder
		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := conf.c.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

func (conf *JadeConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRows, name)
		delete(conf.confMaps, name)
	}
	return nil
}

// NewJadeConfig 初始化Config方法
func NewJadeConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	// 检查文件夹是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}

	// 实例化
	jadeConf := &JadeConfig{
		c:        container,
		folder:   envFolder,
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRows: map[string][]byte{},
		keyDelim: ".",
		lock:     sync.RWMutex{},
	}

	// 读取每个文件
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, file := range files {
		fileName := file.Name()
		err := jadeConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					// 判断事件发生的类型
					// Create 创建
					// Write 写入
					// Remove 删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						jadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						jadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						jadeConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	return jadeConf, nil
}

func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		rekey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(rekey), []byte(val))
	}

	return content
}

// 查找某个路径配置项
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}
	// 判断是否有下个路径
	next, ok := source[path[0]]
	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}

		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果interface的map，使用cast进行下value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// 直接循环调用
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}
	return nil
}

// 通过path获取某个配置项
func (conf *JadeConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func (conf *JadeConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

func (conf *JadeConfig) Get(key string) interface{} {
	return conf.find(key)
}

func (conf *JadeConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

func (conf *JadeConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *JadeConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *JadeConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *JadeConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *JadeConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *JadeConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *JadeConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *JadeConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *JadeConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

func (conf *JadeConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(conf.find(key))
}
