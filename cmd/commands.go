package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdDeploy = &cobra.Command{
	Use:   "deploy [string to print]",
	Short: "Use Karna Deployment to deploy your Lambda application.",
	Long: `Karna Deployment will build and deploy your Lambda function 
	on top of your config file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Karna deploy")
	},
}

var cmdViz = &cobra.Command{
	Use:   "viz [string to echo]",
	Short: "Use Karna Viz to build a Lambda tree on Neo4J.",
	Long: `Karna Viz will build a graph on top of your AWS resources and build
	this tree into Neo4J.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Karna Viz")
	},
}
var apiEcho = &cobra.Command{
	Use:   "api [string to echo]",
	Short: "Use Karna API to build a GUI.",
	Long: `Karna API will start a WebServer which exposes a collection of
	endpoints to build, interact, vizualize your Lambda architecture.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Karna API")
	},
}

func Execute() {
	var rootCmd = &cobra.Command{Use: "karna"}
	rootCmd.AddCommand(cmdDeploy, apiEcho, cmdViz)
	rootCmd.Execute()
}
