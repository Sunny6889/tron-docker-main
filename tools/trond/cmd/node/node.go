package node

import (
	"github.com/spf13/cobra"
)

// NodeCmd represents the node command
var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Commands for operating java-tron docker node.",
	Long: `Commands used for operating node, such as:

  1. check and configure node local environment
  2. deploy single Fullnode/SR for different environment(main network, nile network, private network)
  3. stop single node

  ** coming soon **
  4. deploy multiple nodes in local single machine
  5. stop multiple nodes in local single machine
  6. deploy multiple nodes in different local machines with ssh(one node on one machine)
  7. deploy wallet-cli

Please refer to the available commands below.`,
}

func init() {
}
