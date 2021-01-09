package cmd

import (
	"os"

	"github.com/karnalab/karna/core"
	"github.com/karnalab/karna/internal/create"
	"github.com/karnalab/karna/internal/deploy"

	"github.com/spf13/cobra"
)

var logger *core.KarnaLogger
var rootCmd = &cobra.Command{Use: "karna"}

var deployCmd = &cobra.Command{
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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Karna stack",
	Long:  "Create a Karna stack who will include a Lambda function, an REST API on APIGateway and a API Gateway resource",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Log("Creation in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		functionName, _ := cmd.Flags().GetString("function-name")
		APIName, _ := cmd.Flags().GetString("api-name")
		APIEndpoint, _ := cmd.Flags().GetString("api-endpoint")
		resource, _ := cmd.Flags().GetString("resource")
		verb, _ := cmd.Flags().GetString("verb")

		if elapsed, err := create.Run(&functionName, &APIName, &APIEndpoint, &resource, &verb); err != nil {
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
	var functionName string
	var APIEndpoint, APIName string
	var resource string
	var verb string

	deployCmd.Flags().StringVarP(&target, "target", "t", "", "Function to deploy (JSON key into your config file)")
	deployCmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias to publish")

	deployCmd.MarkFlagRequired("target")
	deployCmd.MarkFlagRequired("alias")

	createCmd.Flags().StringVarP(&functionName, "function-name", "f", "", "The name of the function to create")
	createCmd.Flags().StringVarP(&APIName, "api-name", "a", "", "The name of the REST API to create")
	createCmd.Flags().StringVarP(&APIEndpoint, "api-endpoint", "e", "EDGE", "The endpoint of the REST API. Can be EDGE, REGIONAL or PRIVATE.")
	createCmd.Flags().StringVarP(&resource, "resource", "r", "", "Path og the resource to create in the REST API")
	createCmd.Flags().StringVarP(&verb, "verb", "v", "", "The HTTP verb to create in the resource")

	createCmd.MarkFlagRequired("function-name")
	createCmd.MarkFlagRequired("api-name")
	createCmd.MarkFlagRequired("api-endpoint")
	createCmd.MarkFlagRequired("resource")
	createCmd.MarkFlagRequired("verb")

	rootCmd.AddCommand(deployCmd, createCmd)
}

//Execute => Will register commands && execute the right one.
func Execute() {
	rootCmd.Execute()
}
