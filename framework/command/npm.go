package command

import (
	"a-projects/geekbang/framework/cobra"
	"log"
	"os"
	"os/exec"
)

var npmCommand = &cobra.Command{
	Use:   "npm",
	Short: "运行path/npm程序",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Fatalln("jade npm: should install npm in your PATH")
			return nil
		}

		cmd := exec.Command(path, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return nil
	},
}
