package env

import (
	"github.com/JiadeXu/jade/framework/contract"
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

type JadeEnv struct {
	//contract.Env
	folder string
	maps   map[string]string
}

// NewJadeEnv 有一个参数 .env 文件所在的目录
// NewJadeEnv("/envfolder/") 会读取 /envfolder/.env
// .env的文件格式 FOO_ENV=BAR
func NewJadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewJadeEnv param error")
	}

	// 读取folder文件
	folder := params[0].(string)

	// 实例化
	jadeEnv := &JadeEnv{
		folder: folder,
		// 实例化环境变量 APP_ENV默认设置为开发环境
		maps: map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// 解析folder/.env文件
	file := path.Join(folder, ".env")
	// 读取.env文件，不管任意失败，都不影响后续
	// 打开文件.env
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		// 读取文件
		br := bufio.NewReader(fi)
		for {
			// 按照行进行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			// 按照符号解析
			s := bytes.SplitN(line, []byte{'='}, 2)
			// 不符合过滤
			if len(s) < 2 {
				continue
			}
			// 保存
			key := string(s[0])
			val := string(s[1])
			jadeEnv.maps[key] = val
		}
	}

	// 获取当前程序的环境变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		jadeEnv.maps[pair[0]] = pair[1]
	}

	return jadeEnv, nil
}

// AppEnv 获取表示当前APP环境的变量APP_ENV
func (en *JadeEnv) AppEnv() string {
	return en.Get("APP_ENV")
}

// IsExist 判断一个环境变量是否有背设置
func (en *JadeEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// Get 获取某个环境变量
func (en *JadeEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

//
func (en *JadeEnv) All() map[string]string {
	return en.maps
}
