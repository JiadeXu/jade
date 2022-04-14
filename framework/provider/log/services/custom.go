package services

import (
	"a-projects/geekbang/framework"
	"a-projects/geekbang/framework/contract"
	"io"
)

type JadeCustomLog struct {
	JadeLog
}

func NewJadeCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &JadeConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.c = c

	return log, nil
}
