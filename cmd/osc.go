package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var oscFlags = struct {
	drive  int
	intDiv bool
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

	oscCmd.Flags().IntVar(&oscFlags.drive, "drive", 2, "the output drive strength in mA (2, 4, 6, 8)")
	oscCmd.Flags().BoolVar(&oscFlags.intDiv, "intDiv", false, "use a fractional mutliplier with an integer divider (works only with output Clk0!)")
	oscCmd.Flags().BoolVar(&oscFlags.noInit, "noInit", false, "do not initialize the Si5351")
}

func runOsc(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	refFrequency := device.Crystal.Frequency()
	log.Printf("Crystal @ %.2fHz", refFrequency)

	drive := toOutputDrive(oscFlags.drive)

	if !oscFlags.noInit {
		device.StartSetup()
	}

	if oscFlags.intDiv {
		if len(args) != 1 {
			log.Fatal("intDiv works only with one output")
		}
		frequency, err := ParseFrequency(args[0])
		if err != nil {
			log.Fatal(err)
		}

		multiplier, divider := si5351.FindFractionalMultiplierWithIntegerDivider(refFrequency, frequency)
		device.PLLA().SetupMultiplier(multiplier)
		pllFrequency := multiplier.Multiply(refFrequency)
		log.Printf("PLLA @ %.2fHz: %v", pllFrequency, multiplier)

		device.PrepareOutputs(si5351.PLLA, false, si5351.ClockInputMultisynth, drive, si5351.Clk0)
		device.Clk0().SetupDivider(divider)
		device.Clk0().SetIntegerMode(true)
		outputFrequency := divider.Divide(pllFrequency)
		log.Printf("Clk0 @ %.2fHz: %v", outputFrequency, divider)
	} else {
		f, _ := device.SetupPLL(si5351.PLLA, 900*si5351.MHz)
		log.Printf("PLLA @ %.2fHz: %v", f, device.PLLA().Multiplier)

		for i, arg := range args {
			output := si5351.OutputIndex(i)
			if output > si5351.Clk5 {
				break
			}

			frequency, err := ParseFrequency(arg)
			if err != nil {
				log.Fatal(err)
			}

			device.PrepareOutputs(si5351.PLLA, false, si5351.ClockInputMultisynth, drive, output)
			f, _ := device.SetOutputFrequency(output, frequency)

			log.Printf("Clk%d @ %.2fHz: %v", i, f, device.Clk0().FrequencyDivider)
		}
	}

	if !oscFlags.noInit {
		device.FinishSetup()
	}
}
