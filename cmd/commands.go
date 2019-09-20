package cmd

import (
	"fmt"
	"karna/internal/deploy"
	"karna/internal/viz"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "karna"}

var cmdDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Use Karna Deployment to deploy your Lambda application.",
	Long: `Karna Deployment will build and deploy your Lambda function 
	on top of your config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		target, _ := cmd.Flags().GetString("target")
		alias, _ := cmd.Flags().GetString("alias")

		deploy.Run(&target, &alias)
	},
}

var cmdViz = &cobra.Command{
	Use:   "viz [sub]",
	Short: "Use Karna Viz to build a Lambda tree on Neo4J.",
	Long: `Karna Viz will build a graph on top of your AWS resources and build
	this tree into Neo4J.`,
	Args: cobra.MinimumNArgs(1),
}

var cmdVizShow = &cobra.Command{
	Use:   "show",
	Short: "Feed Neo4J with Lambda tree.",
	Long: `This command will call AWS services with your IAM role to build the Lambda
	tree and its dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		viz.Run()
	},
}

var cmdVizCleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean Neo4J database.",
	Long:  "This subcommand will remove all Neo4J nodes.",
	Run: func(cmd *cobra.Command, args []string) {
		viz.Cleanup()
	},
}

var cmdAPI = &cobra.Command{
	Use:   "api [sub]",
	Short: "Use Karna API to build a GUI.",
	Long: `Karna API will start a WebServer which exposes a collection of
	endpoints to build, interact, vizualize your Lambda architecture.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Karna API")
	},
}

var cmdAPIStart = &cobra.Command{
	Use:   "start",
	Short: "Use Karna API to build a GUI.",
	Long: `Karna API will start a WebServer which exposes a collection of
	endpoints to build, interact, vizualize your Lambda architecture.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Karna start API")
	},
}

func init() {
	var Target string
	var Alias string

	cmdAPI.AddCommand(cmdAPIStart)
	cmdViz.AddCommand(cmdVizShow, cmdVizCleanup)

	cmdDeploy.Flags().StringVarP(&Target, "target", "t", "", "Function to deploy (JSON key into your config file)")
	cmdDeploy.Flags().StringVarP(&Alias, "alias", "a", "", "Alias to publish")

	cmdDeploy.MarkFlagRequired("target")
	cmdDeploy.MarkFlagRequired("alias")

	rootCmd.AddCommand(cmdDeploy, cmdAPI, cmdViz)
}

func Execute() {
	rootCmd.Execute()
}
