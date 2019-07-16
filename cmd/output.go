package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var outputFlags = struct {
	pll       string
	drive     int
	frequency string
	divider   string
	invert    bool
	off       bool
	on        bool
}{}

var outputCmd = &cobra.Command{
	Use:   "output [0-5]",
	Short: "Set parameters of the given output without initializing the Si5351",
	Run:   runSi5351(runOutput),
}

func init() {
	rootCmd.AddCommand(outputCmd)
	outputCmd.Flags().StringVar(&outputFlags.pll, "pll", "", "the PLL to associate the output with")
	outputCmd.Flags().IntVar(&outputFlags.drive, "drive", 0, "the output drive strengs in mA (2, 4, 6, 8)")
	outputCmd.Flags().StringVar(&outputFlags.frequency, "freq", "", "the frequency of the output in Hz")
	outputCmd.Flags().StringVar(&outputFlags.divider, "div", "", "the divider ratio of the output: a,b,c")
	outputCmd.Flags().BoolVar(&outputFlags.invert, "invert", false, "invert the output")
	outputCmd.Flags().BoolVar(&outputFlags.off, "off", false, "power down the output")
	outputCmd.Flags().BoolVar(&outputFlags.on, "on", false, "power up the output")
}

func runOutput(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	if len(args) != 1 {
		log.Fatal("wrong parameters, try output -help")
	}
	output, err := parseOutput(args[0])
	if err != nil {
		log.Fatal(err)
	}

	device.FractionalOutput[output].SetInputSource(si5351.ClockInputMultisynth)

	if outputFlags.pll != "" {
		pll, err := parsePLL(outputFlags.pll)
		if err != nil {
			log.Fatal(err)
		}
		device.FractionalOutput[output].SetPLL(pll)
	}

	if outputFlags.drive != 0 {
		drive := toOutputDrive(outputFlags.drive)
		device.FractionalOutput[output].SetDrive(drive)
	}

	if outputFlags.invert {
		device.FractionalOutput[output].SetInvert(true)
	}

	if outputFlags.off {
		device.FractionalOutput[output].SetPowerDown(true)
	} else if outputFlags.on {
		device.FractionalOutput[output].SetPowerDown(false)
	}

	switch {
	case outputFlags.frequency != "":
		f, err := parseFrequency(outputFlags.frequency)
		if err != nil {
			log.Fatal(err)
		}
		device.SetOutputFrequency(output, f)
	case outputFlags.divider != "":
		a, b, c, err := parseRatio(outputFlags.divider)
		if err != nil {
			log.Fatal(err)
		}
		device.SetOutputDivider(output, a, b, c)
	}
}
