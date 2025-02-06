package docker

import (
	"github.com/spf13/cobra"
)

// DockerCmd represents the docker command
var DockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Commands for operating java-tron docker image.",
	Long: `Commands used for operating docker image, such as:

  1. build java-tron docker image locally
  2. test the built image
  3. pull the latest image from officical docker hub

Please refer to the available commands below.`,
}

func init() {
}
