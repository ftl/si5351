package main

import (
	"log"

	"github.com/ftl/si5351/pkg/i2c"
	"github.com/ftl/si5351/pkg/si5351"
)

func main() {
	const defaultBus = 1

	crystal := si5351.Crystal{BaseFrequency: si5351.Crystal25MHz, Load: si5351.CrystalLoad10PF, CorrectionPPM: 0}
	bus, err := i2c.Open(si5351.DefaultI2CAddress, defaultBus)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	i2c.Debug = true

	s := si5351.New(crystal, bus)

	s.StartSetup()

	s.SetupOutputRaw(si5351.Clk0, si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA)
	s.SetupOutputRaw(si5351.Clk1, si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA)
	fpll, fout, _ := s.SetupQuadratureOutput(si5351.PLLA, si5351.Clk0, si5351.Clk1, 30*si5351.MHz)

	s.FinishSetup()

	log.Printf("PLLA @ %dHz", fpll)
	log.Printf("Clk0 @ %dHz", fout)

	if bus.Err() != nil {
		log.Fatal(err)
	}
}
