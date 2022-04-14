package log

import (
	"github.com/JiadeXu/jade/framework"
	"github.com/JiadeXu/jade/framework/contract"
	"github.com/JiadeXu/jade/framework/provider/log/formatter"
	"github.com/JiadeXu/jade/framework/provider/log/services"
	"io"
	"strings"
)

type JadeLogServiceProvider struct {
	framework.ServiceProvider

	Driver string // Driver

	// 日志级别
	level contract.LogLevel
	// 日志输出个数方法
	Formatter contract.Formatter
	// 日志context上下文信息获取函数
	CtxFielder contract.CtxFielder
	// 日志输出函数
	Output io.Writer
}

func (l *JadeLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			// 默认使用console
			return services.NewJadeConsoleLog
		}

		cs := tcs.(contract.Config)
		l.Driver = strings.ToLower(cs.GetString("log.Driver"))
	}

	// 根据driver的配置项确定
	switch l.Driver {
	case "single":
		return services.NewJadeSingleLog
	case "rotate":
		return services.NewJadeRotateLog
	case "console":
		return services.NewJadeConsoleLog
	case "custom":
		return services.NewJadeCustomLog
	default:
		return services.NewJadeConsoleLog
	}
}

// Boot 启动的时候注入
func (l *JadeLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer 是否延迟加载
func (l *JadeLogServiceProvider) IsDefer() bool {
	return false
}

// Params 定义要传递给实例化方法的参数
func (l *JadeLogServiceProvider) Params(c framework.Container) []interface{} {
	// 获取configService
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	// 设置参数formatter
	if l.Formatter == nil {
		l.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				l.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				l.Formatter = formatter.TextFormatter
			}
		}
	}

	if l.level == contract.UnknownLevel {
		l.level = contract.InfoLevel
		if configService.IsExist("log.level") {
			l.level = logLevel(configService.GetString("log.level"))
		}
	}

	return []interface{}{c, l.level, l.CtxFielder, l.Formatter, l.Output}
}

func (l *JadeLogServiceProvider) Name() string {
	return contract.LogKey
}

// logLevel get level from string
func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
