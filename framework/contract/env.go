package contract

const (
	// EnvProduction 生产
	EnvProduction = "production"
	//
	EnvTesting = "testing"
	//
	EnvDevelopment = "development"
	//
	EnvKey = "jade:env"
)

// env 定义环境变量服务
type Env interface {
	// AppEnv 获取当前的环境变量
	AppEnv() string
	// IsExist 判断一个环境变量是否有背设置
	IsExist(string) bool
	// Get 获取某个环境变量
	Get(string) string
	//
	All() map[string]string
}
