package snapshot

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronprotocol/tron-docker/utils"
)

// listCmd represents the listFull command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download target backup snapshot to local current directory",
	Long: `Refer to the snapshot source domain and backup name you input, the available backup snapshot will be downloaded to the local directory.

Note:
 - because some snapshot sources have multiple snapshot types, you need to specify the type(full, lite) of snapshot you want to download.
 - the snapshot is large, it may need a long time to finish the download, depends on your network performance.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the flag value
		domain, _ := cmd.Flags().GetString("domain")

		// Get the flag value
		backup, _ := cmd.Flags().GetString("backup")

		// Get the flag value
		nType, _ := cmd.Flags().GetString("type")

		if !utils.CheckDomain(domain) {
			fmt.Println("Error: domain value not supported\nRun \"./trond snapshot source\" to check available items")
			return
		}

		downloadURLMD5 := utils.GenerateSnapshotMD5DownloadURL(domain, backup, nType)
		if downloadURLMD5 == "" {
			fmt.Println("Error: --type value not supported, available: full, lite")
			return
		}
		downloadMD5File, err := utils.DownloadFileWithProgress(downloadURLMD5, "")
		if err != nil {
			fmt.Println("Error:", err)
		}

		downloadURL := utils.GenerateSnapshotDownloadURL(domain, backup, nType)
		if downloadURL == "" {
			fmt.Println("Error: --type value not supported, available: full, lite")
			return
		}
		if _, err := utils.DownloadFileWithProgress(downloadURL, downloadMD5File); err != nil {
			fmt.Println("Error:", err)
		}

	},
}

func init() {
	SnapshotCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringP(
		"domain", "d", "",
		"Domain for target snapshot source(required).\nPlease run command \"./trond snapshot source\" to get the available snapshot source domains.") // -d or --domain
	downloadCmd.Flags().StringP(
		"backup", "b", "",
		"Backup name(required).\nPlease run command \"./trond snapshot list\" to get the available backup name under target source domains.") // -b or --backup
	downloadCmd.Flags().StringP(
		"type", "t", "",
		"Node type of the snapshot(required, available: full, lite).")

	// Mark source and destination as required flags
	if err := downloadCmd.MarkFlagRequired("domain"); err != nil {
		log.Fatalf("Error marking domain flag as required: %v", err)
	}
	if err := downloadCmd.MarkFlagRequired("backup"); err != nil {
		log.Fatalf("Error marking backup flag as required: %v", err)
	}
	if err := downloadCmd.MarkFlagRequired("type"); err != nil {
		log.Fatalf("Error marking type flag as required: %v", err)
	}

}
