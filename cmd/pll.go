package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ftl/si5351/pkg/si5351"
)

var pllFlags = struct {
	frequency  string
	multiplier string
	reset      bool
}{}

var pllCmd = &cobra.Command{
	Use:   "pll [A|B]",
	Short: "Set the frequency of the given PLL without initializing the Si5351",
	Run:   runSi5351(runPLL),
}

func init() {
	rootCmd.AddCommand(pllCmd)
	pllCmd.Flags().StringVar(&pllFlags.frequency, "freq", "", "the frequency of the PLL in Hz")
	pllCmd.Flags().StringVar(&pllFlags.multiplier, "multi", "", "the multiplier ratio of the PLL: a,b,c")
	pllCmd.Flags().BoolVar(&pllFlags.reset, "reset", false, "reset the PLL")
}

func runPLL(cmd *cobra.Command, args []string, device *si5351.Si5351) {
	if len(args) != 1 {
		log.Fatal("wrong parameters, try pll -help")
	}
	pll, err := parsePLL(args[0])
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case pllFlags.frequency != "":
		f, err := parseFrequency(pllFlags.frequency)
		if err != nil {
			log.Fatal(err)
		}
		device.SetupPLL(pll, f)
	case pllFlags.multiplier != "":
		a, b, c, err := parseRatio(pllFlags.multiplier)
		if err != nil {
			log.Fatal(err)
		}
		device.SetupPLLRaw(pll, a, b, c)
	}

	if pllFlags.reset {
		device.PLL[pll].Reset()
	}
}
