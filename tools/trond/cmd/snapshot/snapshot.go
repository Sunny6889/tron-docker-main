/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"
)

// SnapshotCmd represents the snapshot command
var SnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Command sets for managing Java-Tron node snapshots.",
	Long: `Commands used for downloading node's snapshot, such as:

  1. show data source
  2. list available snapshots of target data source
  3. download target snapshot

Please refer to the available commands below.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshot called")
	},
}

func init() {
}
