package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var testFlags = struct {
	test string
}{}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A simple frame to test things",
	Run:   runSi5351(runTest),
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVar(&testFlags.test, "test", "", "a test string parameter")
}

func runTest(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	fmt.Printf("testing Si5351 @ 0x%x on I2C bus #%d %s\n", rootFlags.address, rootFlags.bus, testFlags.test)

	device.StartSetup()

	device.SetupOutputRaw(si5351.Clk0, si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA)
	device.SetupOutputRaw(si5351.Clk1, si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA)
	fpll, fout, _ := device.SetupQuadratureOutput(si5351.PLLA, si5351.Clk0, si5351.Clk1, 30*si5351.MHz)

	device.FinishSetup()

	log.Printf("PLLA @ %dHz", fpll)
	log.Printf("Clk0 @ %dHz", fout)
}
