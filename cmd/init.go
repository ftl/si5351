package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var initFlags = struct {
}{}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Si5351",
	Run:   runSi5351(runInit),
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	device.StartSetup()
	device.FinishSetup()
}
