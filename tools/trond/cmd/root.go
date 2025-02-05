package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tronprotocol/tron-docker/cmd/node"
	"github.com/tronprotocol/tron-docker/cmd/snapshot"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "trond",
	Short: "Docker automation for TRON nodes",
	Long:  `This tool bundles multiple commands into one, enabling the community to quickly get started with TRON network interaction and development.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(snapshot.SnapshotCmd)
	rootCmd.AddCommand(node.NodeCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tron-docker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
