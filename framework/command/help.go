package command

import (
	"github.com/JiadeXu/jade/framework/cobra"
	"github.com/JiadeXu/jade/framework/contract"
	"fmt"
)

// helpCommand show current envionment
var DemoCommand = &cobra.Command{
	Use:   "demo.md",
	Short: "demo.md for framework",
	Run: func(c *cobra.Command, args []string) {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println("app base folder:", appService.BaseFolder())
	},
}
