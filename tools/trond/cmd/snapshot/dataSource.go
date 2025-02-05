package snapshot

import (
	"github.com/spf13/cobra"
	"github.com/tronprotocol/tron-docker/utils"
)

// sourceCmd represents the snapshot source command
var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Show available snapshot source.",
	Long: `Available snapshot sources will be shown.
Support different types of snapshot (Fullnode, Lite Fullnode), in different regions (Singapore, America).
You can choose the one you need.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.ShowSnapshotDataSourceList()
	},
}

func init() {
	SnapshotCmd.AddCommand(sourceCmd)
}
