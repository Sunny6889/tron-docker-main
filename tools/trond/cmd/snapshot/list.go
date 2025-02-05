package snapshot

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronprotocol/tron-docker/utils"
)

// listCmd represents the listFull command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available snapshots of target source.",
	Long: `Refer to the snapshot source domain you input, the available backup snapshots will be showen below.
Note: different domain may have different snapshots that can be downloaded.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get the flag value
		domain, _ := cmd.Flags().GetString("domain")

		// Validate flag values
		if domain == "" {
			fmt.Println("Error: --domain flag cannot be empty")
			return
		}

		if !utils.CheckDomain(domain) {
			fmt.Println("Error: domain value not supported\nRun \"./trond snapshot source\" to check available items")
			return
		}

		if err := utils.ShowSnapshotList(domain, utils.IsHttps(domain)); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {

	listCmd.Flags().StringP("domain", "d", "", "Domain for target snapshot source (required)\nPlease run command \"./trond snapshot source\" to get the available snapshot source domains")

	if err := listCmd.MarkFlagRequired("domain"); err != nil {
		log.Fatalf("Error marking domain flag as required: %v", err)
	}

	SnapshotCmd.AddCommand(listCmd)
}
