package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown the outputs of the Si5351",
	Run:   runSi5351(runShutdown),
}

func init() {
	rootCmd.AddCommand(shutdownCmd)
}

func runShutdown(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	device.Shutdown()
}
