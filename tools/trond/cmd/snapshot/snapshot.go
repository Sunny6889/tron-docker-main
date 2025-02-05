package snapshot

import (
	"github.com/spf13/cobra"
)

// SnapshotCmd represents the snapshot command
var SnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Commands for getting java-tron node snapshots.",
	Long: `Commands used for downloading node's snapshot, such as:

  1. show available snapshot source
  2. list available snapshots in target source
  3. download target snapshot`,
}

func init() {
}
