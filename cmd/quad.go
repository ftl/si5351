package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var quadFlags = struct {
	drive  int
	noInit bool
}{}

var quadCmd = &cobra.Command{
	Use:   "quad [pll] [i output] [q output] [frequency]",
	Short: "Output the given frequency on the two given outputs with a phase shift of 90°, using the given PLL.",
	Run:   runSi5351(runQuad),
}

func init() {
	rootCmd.AddCommand(quadCmd)

	quadCmd.Flags().IntVar(&quadFlags.drive, "drive", 2, "the output drive strength in mA (2, 4, 6, 8)")
	quadCmd.Flags().BoolVar(&quadFlags.noInit, "noInit", false, "do not initialize the Si5351")
}

func runQuad(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	if len(args) != 4 {
		log.Fatal("wrong number of arguments, try quad --help")
	}

	var err error
	pll, err := parsePLL(args[0])
	iOutput, err := parseOutput(args[1])
	qOutput, err := parseOutput(args[2])
	frequency, err := ParseFrequency(args[3])
	drive := toOutputDrive(quadFlags.drive)

	if err != nil {
		log.Fatal(err)
	}

	if !quadFlags.noInit {
		device.StartSetup()
	}

	device.PrepareOutputs(pll, false, si5351.ClockInputMultisynth, drive, iOutput, qOutput)
	device.SetupQuadratureOutput(pll, iOutput, qOutput, frequency)

	if !quadFlags.noInit {
		device.FinishSetup()
	}
}
