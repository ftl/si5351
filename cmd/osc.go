package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var oscFlags = struct {
	noInit bool
}{}

var oscCmd = &cobra.Command{
	Use:   "osc [freq0] [freq1] [freq2] [freq3] [freq4] [freq5]",
	Short: "Output the given frequencies on the outputs CLK0-CLK5 using PLL A",
	Long: `Output the given frequencies on the outputs CLK0-CLK5 using PLL A.
If the list of given frequencies is shorter than six entries, only the outputs with given frequencies are setup.

Example: osc 10M 5M 3500k 3400k # output 10MHz, 5MHz, 3500kHz, and 3400kHz on the outputs CLK0-CLK4
`,
	Run: runSi5351(runOsc),
}

func init() {
	rootCmd.AddCommand(oscCmd)

	oscCmd.Flags().BoolVar(&oscFlags.noInit, "noInit", false, "do not initialize the Si5351")
}

func runOsc(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	if !oscFlags.noInit {
		device.StartSetup()
	}

	f, _ := device.SetupPLL(si5351.PLLA, 900*si5351.MHz)
	log.Printf("PLLA @ %dHz", f)

	for i, arg := range args {
		output := si5351.OutputIndex(i)
		if output > si5351.Clk5 {
			break
		}

		frequency, err := parseFrequency(arg)
		if err != nil {
			log.Fatal(err)
		}

		device.PrepareOutputs(si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA, output)
		f, _ := device.SetOutputFrequency(output, frequency)

		log.Printf("Clk%d @ %dHz", i, f)
	}

	if !oscFlags.noInit {
		device.FinishSetup()
	}
}
