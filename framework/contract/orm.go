package contract

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"net"
	"strconv"
	"time"
)

const ORMKey = "jade:orm"

type ORMService interface {
	GetDB(option ...DBOption) (*gorm.DB, error)
}

// DBOption 代表初始化的时候的选项
type DBOption func(container framework.Container, config *DBConfig) error

type DBConfig struct {
	// 以下配置关于dsn
	WriteTimeout string `yaml:"write_timeout"` // 写超时
	loc          string `yaml:"loc"`           // 时区
	Port         int    `yaml:"port"`          // 端口
	ReadTimeout  string `yaml:"read_timeout"`  // 读超时
	Charset      string `yaml:"charset"`       // 字符集
	ParseTime    bool   `yaml:"parse_time"`    // 是否解析时间
	Protocol     string `yaml:"protocol"`      // 传输协议
	Dsn          string `yaml:"dsn"`           // 直接传递dsn，如果传递了，其他关于dsn的配置均无效
	Database     string `yaml:"database"`      // 数据库
	Collation    string `yaml:"collation"`     // 字符序
	Timeout      string `yaml:"timeout"`       // 连接超时时间
	Username     string `yaml:"username"`      // 用户名
	Password     string `yaml:"password"`      // 密码
	Driver       string `yaml:"driver"`        // 驱动
	Host         string `yaml:"host"`          // 数据库地址

	// 以下配置关于连接池
	ConnMaxIdle     int    `yaml:"conn_max_idle"`     // 最大空闲连接数
	ConnMaxOpen     int    `yaml:"conn_max_open"`     // 最大连接数
	ConnMaxLifetime string `yaml:"conn_max_lifetime"` // 连接最大生命周期
	ConnMaxIdletime string `yaml:"conn_max_idletime"` // 空闲最大生命周期

	// 关于gorm
	*gorm.Config
}

// FormatDsn 生成dsn
func (conf *DBConfig) FormatDsn() (string, error) {
	port := strconv.Itoa(conf.Port)
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return "", err
	}
	readTimeout, err := time.ParseDuration(conf.ReadTimeout)
	if err != nil {
		return "", err
	}
	writeTimeout, err := time.ParseDuration(conf.WriteTimeout)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation(conf.loc)
	if err != nil {
		return "", err
	}
	driveConf := mysql.Config{
		User: conf.Username,
		Passwd: conf.Password,
		Net: conf.Protocol,
		Addr: net.JoinHostPort(conf.Host, port),
		DBName: conf.Database,
		Collation: conf.Collation,
		Loc: location,
		Timeout: timeout,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		ParseTime: conf.ParseTime,
		AllowNativePasswords: true,
	}
	return driveConf.FormatDSN(), nil
}