package node

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronprotocol/tron-docker/utils"
)

// runSingleCmd represents the dataSource command
var runSingleCmd = &cobra.Command{
	Use:   "run-single",
	Short: "Run single java-tron node for different networks.",
	Long: `You need to make sure the local environment is ready before running the node.
Use this command "./trond node env" to check the environment before starting the node.

The following files are required:
  - Database directory: ./output-directory
  - Configuration file(by default, these exist in the current repository directory)
      main network: ../../conf/main_net_config.conf
      nile network: ../../conf/nile_net_config.conf
      private network: ../../conf/private_net_config_*.conf
  - Docker compose file(by default, these exist in the current repository directory)
        main network: ../../single_node/docker-compose.fullnode.main.yaml
        nile network: ../../single_node/docker-compose.fullnode.nile.yaml
        private network: ../../single_node/docker-compose.witness.private.yaml
  - Log directory: ./logs`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
		if len(service) == 0 {
			service = "tron-node"
		}

		nType, _ := cmd.Flags().GetString("type")

		dockerComposeFile := ""
		switch nType {
		case "full-main":
			dockerComposeFile = "../../single_node/docker-compose.fullnode.main.yaml"
		case "full-nil":
			dockerComposeFile = "../../single_node/docker-compose.fullnode.nile.yaml"
		case "witness-private":
			dockerComposeFile = "../../single_node/docker-compose.witness.private.yaml"
		default:
			fmt.Println("Error: type not supported", nType)
		}
		if yes, isDir := utils.PathExists(dockerComposeFile); !yes || isDir {
			fmt.Println("Error: file not exists or not a file:", dockerComposeFile)
		}
		if err := utils.RunComposeServiceOnce(dockerComposeFile, service); err != nil {
			fmt.Println("Error: ", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(runSingleCmd)

	runSingleCmd.Flags().StringP(
		"type", "t", "",
		"Node type you want to deploy (required, available: full-main, full-nile, witness-private)")

	runSingleCmd.Flags().StringP(
		"service", "s", "",
		"Service name you want to start in docker-compose.*.yaml (optional, default: tron)")

	if err := runSingleCmd.MarkFlagRequired("type"); err != nil {
		log.Fatalf("Error marking type flag as required: %v", err)
	}
}
