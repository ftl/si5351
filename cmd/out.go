package cmd

import (
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var outFlags = struct {
	out string
}{}

var outCmd = &cobra.Command{
	Use:   "out [freq0] [freq1] [freq2] [freq3] [freq4] [freq5]",
	Short: "Output the given frequencies on the outputs CLK0-CLK5 using PLL A",
	Long: `Output the given frequencies on the outputs CLK0-CLK5 using PLL A.
If the list of given frequencies is shorter than six entries, only the outputs with given frequencies are setup.

Example: out 10M 5M 3500k 3400k # output 10MHz, 5MHz, 3500kHz, and 3400kHz on the outputs CLK0-CLK4
`,
	Run: runSi5351(runOut),
}

func init() {
	rootCmd.AddCommand(outCmd)
}

func runOut(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	device.StartSetup()

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

		device.SetupOutputRaw(output, si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA)
		f, _ := device.SetOutputFrequency(output, frequency)

		log.Printf("Clk%d @ %dHz", i, f)
	}

	device.FinishSetup()
}

func parseFrequency(f string) (si5351.Frequency, error) {
	input := strings.ToLower(strings.TrimSpace(f))
	var magnitude si5351.Frequency
	switch {
	case strings.HasSuffix(input, "m"):
		magnitude = si5351.MHz
		input = input[:len(input)-1]
	case strings.HasSuffix(input, "k"):
		magnitude = si5351.KHz
		input = input[:len(input)-1]
	default:
		magnitude = si5351.Hz
	}
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return si5351.Frequency(value) * magnitude, nil
}
