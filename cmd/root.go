package cmd

import (
	"log"

	"github.com/ftl/si5351/pkg/i2c"
	"github.com/ftl/si5351/pkg/si5351"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootFlags = struct {
	address     uint8
	bus         int
	debugI2C    bool
	crystalFreq int
	crystalLoad int
	ppm         int
}{}

var rootCmd = &cobra.Command{
	Use:   "si5351",
	Short: "Control the Si5351",
}

// Execute is called by main.main() as the entry point to the Cobra framework.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().Uint8Var(&rootFlags.address, "address", si5351.DefaultI2CAddress, "the I2C address of the Si5351")
	rootCmd.PersistentFlags().IntVar(&rootFlags.bus, "bus", 1, "the I2C bus number to which the Si5351 is attached to")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.debugI2C, "debugI2C", false, "enable debug output of the communication on the I2C bus")
	rootCmd.PersistentFlags().IntVar(&rootFlags.crystalFreq, "crystalFreq", 25, "the frequency of the crystal in MHz (25, 27)")
	rootCmd.PersistentFlags().IntVar(&rootFlags.crystalLoad, "crystalLoad", 10, "the internal capacitive load of the crystal in pF (6, 8, 10)")
	rootCmd.PersistentFlags().IntVar(&rootFlags.ppm, "ppm", 0, "the frequency correction of the crystal in PPM")
}

func runSi5351(f func(cmd *cobra.Command, args []string, device *si5351.Si5351)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		crystal := si5351.Crystal{BaseFrequency: toCrystalFrequency(rootFlags.crystalFreq), Load: toCrystalLoad(rootFlags.crystalLoad), CorrectionPPM: rootFlags.ppm}
		bus, err := i2c.Open(rootFlags.address, rootFlags.bus)
		if err != nil {
			log.Fatal(err)
		}
		defer bus.Close()
		i2c.Debug = rootFlags.debugI2C

		device := si5351.New(crystal, bus)

		f(cmd, args, device)

		if bus.Err() != nil {
			log.Fatal(err)
		}
	}
}
