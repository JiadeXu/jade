package command

import "a-projects/geekbang/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initEnvCommand())
	root.AddCommand(initConfigCommand())
	root.AddCommand(initBuildCommand())

	root.AddCommand(goCommand)
	root.AddCommand(npmCommand)

	root.AddCommand(initDevCommand())

	root.AddCommand(initProviderCommand())

	root.AddCommand(initCmdCommand())
	root.AddCommand(initMiddlewareCommand())
}
