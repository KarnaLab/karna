package cmd

import (
	"os"

	internal "github.com/karnalab/karna/internal"
	"github.com/spf13/cobra"
)

var logger *internal.KarnaLogger
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
		functionName, _ := cmd.Flags().GetString("function-name")
		alias, _ := cmd.Flags().GetString("alias")

		if elapsed, err := internal.Deploy(&functionName, &alias); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		} else {
			logger.Log("Completed in " + elapsed)
		}
	},
}

var cmdAlias = &cobra.Command{
	Use:   "remove-alias",
	Short: "Use Karna Deployment to deploy your Lambda application.",
	Long: `Karna Deployment will build and deploy your Lambda function 
	on top of your config file.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Log("Alias removing in progress...")
	},
	Run: func(cmd *cobra.Command, args []string) {
		functionName, _ := cmd.Flags().GetString("function-name")
		alias, _ := cmd.Flags().GetString("alias")

		if elapsed, err := internal.RemoveAlias(&functionName, &alias); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		} else {
			logger.Log("Completed in " + elapsed)
		}
	},
}

func init() {
	var functionName string
	var alias string

	cmdDeploy.Flags().StringVarP(&functionName, "function-name", "f", "", "Function to deploy (JSON key into your config file)")
	cmdDeploy.Flags().StringVarP(&alias, "alias", "a", "", "Alias to publish")

	cmdDeploy.MarkFlagRequired("function-name")
	cmdDeploy.MarkFlagRequired("alias")

	cmdAlias.Flags().StringVarP(&functionName, "function-name", "f", "", "Function to deploy (JSON key into your config file)")
	cmdAlias.Flags().StringVarP(&alias, "alias", "a", "", "Alias to publish")

	cmdAlias.MarkFlagRequired("function-name")
	cmdAlias.MarkFlagRequired("alias")

	rootCmd.AddCommand(cmdDeploy)
	rootCmd.AddCommand(cmdAlias)
}

//Execute => Will register commands && execute the right one.
func Execute() {
	rootCmd.Execute()
}
