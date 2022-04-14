package command

import (
	"github.com/JiadeXu/jade/framework/cobra"
	"github.com/JiadeXu/jade/framework/contract"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// app 启动地址
var appAddress = ""

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().StringVar(&appAddress, "address", ":8888", "设置app启动的地址，默认为8888")

	appCommand.AddCommand(appStartCommand)
	return appCommand
}

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 打印帮助文档
		cmd.Help()
		return nil
	},
}

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 启动
		container := cmd.GetContainer()
		// 获取kernel
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()

		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}

		//
		go func() {
			server.ListenAndServe()
		}()

		// 当前的goroutine等待信号量
		quit := make(chan os.Signal)
		// 监控信号：SIGINT, SIGTERM, SIGQUIT
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 这里会阻塞当前goroutine等待信号
		<-quit

		// 调用Server.Shutdown graceful结束
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		return nil
	},
}
