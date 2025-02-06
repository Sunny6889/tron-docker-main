package node

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tronprotocol/tron-docker/utils"
)

// envCmd represents the dataSource command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Check and configure node local environment",
	Long: `Default environment configuration for node operation:

Current directory: tron-docker/tools/trond

The following files are required:
  - Database directory: ./output-directory (if not exists, will create it)
  - Configuration file(by default, these exist in the current repository directory)
      main network: ../../conf/main_net_config.conf
      nile network: ../../conf/nile_net_config.conf
      private network: ../../conf/private_net_config_*.conf
  - Docker compose file(by default, these exist in the current repository directory)
      single node
        main network: ../../single_node/docker-compose.fullnode.main.yaml
        nile network: ../../single_node/docker-compose.fullnode.nile.yaml
        private network: ../../single_node/docker-compose.witness.private.yaml
      multiple nodes
	    private network: ../../private_net/docker-compose.private.yaml
  - Log directory: ./logs (if not exists, will create it)`,
	Run: func(cmd *cobra.Command, args []string) {
		if yes, err := utils.PwdEndsWith("tron-docker/tools/trond"); err != nil || !yes {
			fmt.Println("Error: current directory is wrong, should be tron-docker/tools/trond")
			return
		}

		checkDirectory := map[string]bool{
			"../../conf":         false,
			"./output-directory": true,
			"./logs":             true,
			"../../single_node":  false,
			"../../private_net":  false,
		}
		checkFile := []string{
			"../../conf/main_net_config.conf",
			"../../conf/nile_net_config.conf",
			"../../conf/private_net_config_witness1.conf",
			"../../conf/private_net_config_witness2.conf",
			"../../conf/private_net_config_others.conf",
			"../../single_node/docker-compose.fullnode.main.yaml",
			"../../single_node/docker-compose.fullnode.nile.yaml",
			"../../single_node/docker-compose.witness.private.yaml",
			"../../private_net/docker-compose.private.yaml",
		}
		for k, v := range checkDirectory {
			if yes, isDir := utils.PathExists(k); !yes {
				if v {
					fmt.Println("Warning: directory not exists:", k)
					fmt.Println(" - Creating it")
					if err := utils.CreateDir(k, true); err != nil {
						fmt.Println(" - Error creating directory:", err)
					} else {
						fmt.Println(" - Directory created successfully:", k)

						if k == "./output-directory" {
							fmt.Println(" - Note: no history database, the node will start from 0 block")
						}
					}
				} else {
					fmt.Println("Error: directory not exists:", k)
				}
			} else if !isDir {
				fmt.Println("Error: target is not a directory:", k)
			} else {
				fmt.Println("Directory exists:", k)
			}
		}
		for _, v := range checkFile {
			if yes, isDir := utils.PathExists(v); !yes || isDir {
				fmt.Println("Error: file not exists or not a file:", v)
			}
		}
	},
}

func init() {
	NodeCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
