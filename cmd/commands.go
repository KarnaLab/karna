package cmd

import (
	"karna/core"
	"karna/internal/api"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		core.LogSuccessMessage("Deployment in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		target, _ := cmd.Flags().GetString("target")
		alias, _ := cmd.Flags().GetString("alias")

		elapsed := deploy.Run(&target, &alias)
		core.LogSuccessMessage("Completed in " + elapsed)
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
	PreRun: func(cmd *cobra.Command, args []string) {
		core.LogSuccessMessage("Create Neo4J trees in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		credentials, _ := cmd.Flags().GetString("credentials")
		host, _ := cmd.Flags().GetString("host")

		elapsed := viz.Run(&port, &credentials, &host)
		core.LogSuccessMessage("Completed in " + elapsed)
	},
}

var cmdVizCleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean Neo4J database.",
	Long:  "This subcommand will remove all Neo4J nodes.",
	PreRun: func(cmd *cobra.Command, args []string) {
		core.LogSuccessMessage("Cleaning Neo4J in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		elapsed := viz.Cleanup()
		core.LogSuccessMessage("Completed in " + elapsed)
	},
}

var cmdAPI = &cobra.Command{
	Use:   "api [sub]",
	Short: "Use Karna API to build a GUI.",
	Long: `Karna API will start a WebServer which exposes a collection of
	endpoints to build, interact, vizualize your Lambda architecture.`,
	Args: cobra.MinimumNArgs(1),
}

var cmdAPIStart = &cobra.Command{
	Use:   "start",
	Short: "Use Karna API to build a GUI.",
	Long: `Karna API will start a WebServer which exposes a collection of
	endpoints to build, interact and vizualize your Lambda architecture.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		core.LogSuccessMessage("Starting API in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}

func init() {
	var target string
	var alias string
	var port string
	var credentials string
	var host string

	cmdAPI.AddCommand(cmdAPIStart)
	cmdViz.AddCommand(cmdVizShow, cmdVizCleanup)

	cmdDeploy.Flags().StringVarP(&target, "target", "t", "", "Function to deploy (JSON key into your config file)")
	cmdDeploy.Flags().StringVarP(&alias, "alias", "a", "", "Alias to publish")

	cmdDeploy.MarkFlagRequired("target")
	cmdDeploy.MarkFlagRequired("alias")

	cmdVizShow.Flags().StringVarP(&port, "port", "p", "", "Database port")
	cmdVizShow.Flags().StringVarP(&credentials, "credentials", "c", "", "Credentials for Neo4J database")
	cmdVizShow.Flags().StringVarP(&host, "host", "", "", "Host for Neo4J database")

	rootCmd.AddCommand(cmdDeploy, cmdAPI, cmdViz)
}

//Execute => Will register commands && execute the right one.
func Execute() {
	rootCmd.Execute()
}
