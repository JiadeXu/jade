package cobra

import (
	"a-projects/geekbang/framework"
	"github.com/robfig/cron/v3"
	"log"
)

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

// CronSpec 保存Cron命令的信息 用于展示
type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetParantNull() {
	c.parent = nil
}

// AddCronCommand 是用来创建一个cron任务
func (c *Command) AddCronCommand(spec string, cmd *Command) {
	// cron 结构是挂载在根Command上的
	root := c.Root()
	if root.Cron == nil {
		// 初始化cron
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}
	// 增加说明信息
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})

	// 制作一个rootCommand
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParantNull()
	cronCmd.SetContainer(root.GetContainer())

	// 增加调用函数
	root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			// 打印出err信息
			log.Println(err)
		}
	})
}
