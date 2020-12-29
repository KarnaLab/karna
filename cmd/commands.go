package cmd

import (
	"os"

	"github.com/karnalab/karna/core"
	"github.com/karnalab/karna/internal/deploy"
	"github.com/karnalab/karna/internal/viz"

	"github.com/spf13/cobra"
)

var logger *core.KarnaLogger
var rootCmd = &cobra.Command{Use: "karna"}

var cmdDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Use Karna Deployment to deploy your Lambda application.",
	Long: `Karna Deployment will build and deploy your Lambda function 
	on top of your config file.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Log("Deployment in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		target, _ := cmd.Flags().GetString("target")
		alias, _ := cmd.Flags().GetString("alias")

		if elapsed, err := deploy.Run(&target, &alias); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		} else {
			logger.Log("Completed in " + elapsed)
		}
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
		logger.Log("Create Neo4J trees in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		credentials, _ := cmd.Flags().GetString("credentials")
		host, _ := cmd.Flags().GetString("host")

		if elapsed, err := viz.Run(&port, &credentials, &host); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		} else {
			logger.Log("Completed in " + elapsed)
		}
	},
}

var cmdVizCleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean Neo4J database.",
	Long:  "This subcommand will remove all Neo4J nodes.",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Log("Cleaning Neo4J in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		credentials, _ := cmd.Flags().GetString("credentials")
		host, _ := cmd.Flags().GetString("host")

		if elapsed, err := viz.Cleanup(&port, &credentials, &host); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		} else {
			logger.Log("Completed in " + elapsed)
		}
	},
}

func init() {
	var target string
	var alias string
	var port string
	var credentials string
	var host string

	cmdViz.AddCommand(cmdVizShow, cmdVizCleanup)

	cmdDeploy.Flags().StringVarP(&target, "target", "t", "", "Function to deploy (JSON key into your config file)")
	cmdDeploy.Flags().StringVarP(&alias, "alias", "a", "", "Alias to publish")

	cmdDeploy.MarkFlagRequired("target")
	cmdDeploy.MarkFlagRequired("alias")

	cmdVizShow.Flags().StringVarP(&port, "port", "p", "", "Database port")
	cmdVizShow.Flags().StringVarP(&credentials, "credentials", "c", "", "Credentials for Neo4J database")
	cmdVizShow.Flags().StringVarP(&host, "host", "", "", "Host for Neo4J database")

	cmdVizCleanup.Flags().StringVarP(&port, "port", "p", "", "Database port")
	cmdVizCleanup.Flags().StringVarP(&credentials, "credentials", "c", "", "Credentials for Neo4J database")
	cmdVizCleanup.Flags().StringVarP(&host, "host", "", "", "Host for Neo4J database")

	rootCmd.AddCommand(cmdDeploy, cmdViz)
}

//Execute => Will register commands && execute the right one.
func Execute() {
	rootCmd.Execute()
}
