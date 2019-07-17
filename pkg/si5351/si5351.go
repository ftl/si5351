package si5351

import (
	"errors"
	"io"
)

// DefaultI2CAddress is the default address of the Si5351 on the I2C bus.
const DefaultI2CAddress uint8 = 0x60

// Si5351 represents the chip.
type Si5351 struct {
	Crystal      Crystal
	InputDivider ClockDivider

	PLL              []*PLL
	FractionalOutput []*FractionalOutput
	IntegerOutput    []*IntegerOutput

	bus Bus
}

// Bus on which to communicate with the Si5351.
type Bus interface {
	ReadReg(reg uint8, p []byte) (int, error)
	WriteReg(reg uint8, values ...byte) (int, error)
	RegWriter(reg uint8) io.Writer
	Err() error
	Close() error
}

// New returns a new Si5351 instance.
func New(crystal Crystal, bus Bus) *Si5351 {
	return &Si5351{
		Crystal:          crystal,
		PLL:              loadPLLs(bus),
		FractionalOutput: loadFractionalOutputs(bus),
		IntegerOutput:    loadIntegerOutputs(bus),
		bus:              bus,
	}
}

// StartSetup starts the setup sequence of the Si5351:
// * disable all outputs
// * power down all output drivers
// * set the CLKIN input divider
// * set the internal load capacitance of the crystal
// After these steps the individual setup of PLLs and Clocks should take place.
// As last setup step, don't forget to call FinishSetup.
func (s *Si5351) StartSetup() error {
	s.Shutdown()
	s.bus.WriteReg(RegCrystalInternalLoadCapacitance, byte(s.Crystal.Load))
	return s.bus.Err()
}

// FinishSetup finishes the setup sequence:
// * reset the PLLs
// * enable all outputs
func (s *Si5351) FinishSetup() error {
	s.resetAllPLLs()
	s.enableAllOutputs(true)
	return s.bus.Err()
}

// SetupPLLInputSource writes the input source configuration to the Si5351's register.
func (s *Si5351) SetupPLLInputSource(clkinInputDivider ClockDivider, pllASource, pllBSource PLLInputSource) error {
	value := byte((clkinInputDivider&0xF)<<4) |
		byte((pllASource&1)<<s.PLL[PLLA].Register.InputSourceOffset) |
		byte((pllBSource&1)<<s.PLL[PLLB].Register.InputSourceOffset)

	s.bus.WriteReg(RegPLLInputSource, value)

	if s.bus.Err() == nil {
		s.InputDivider = clkinInputDivider
		s.PLL[PLLA].InputSource = pllASource
		s.PLL[PLLB].InputSource = pllBSource
	}
	return s.bus.Err()
}

// SetupPLLRaw directly sets the frequency multiplier parameters for the given PLL and resets it.
func (s *Si5351) SetupPLLRaw(pll PLLIndex, a, b, c int) error {
	s.PLL[pll].SetupMultiplier(FractionalRatio{A: a, B: b, C: c})
	s.PLL[pll].Reset()
	return s.bus.Err()
}

// SetupOutputRaw directly sets the control parameters of the given output.
func (s *Si5351) SetupOutputRaw(output OutputIndex, pll PLLIndex, invert bool, inputSource ClockInputSource, drive OutputDrive) error {
	if int(output) < len(s.FractionalOutput) {
		s.FractionalOutput[output].SetupControl(false, false, pll, invert, inputSource, drive)
	} else {
		s.IntegerOutput[int(output)-len(s.FractionalOutput)].SetupControl(false, false, pll, invert, inputSource, drive)
	}

	return s.bus.Err()
}

// SetupMultisynthRaw directly sets the frequency divider and RDiv parameters for the Multisynth of the given output.
func (s *Si5351) SetupMultisynthRaw(output OutputIndex, a, b, c int, RDiv ClockDivider) error {
	if int(output) >= len(s.FractionalOutput) {
		return errors.New("only CLK0-CLK5 are currently supported")
	}

	s.FractionalOutput[output].SetupDivider(FractionalRatio{A: a, B: b, C: c})

	return s.bus.Err()
}

// SetupPLL sets the given PLL to the closest possible value of the given frequency and resets it.
func (s *Si5351) SetupPLL(pll PLLIndex, frequency Frequency) (Frequency, error) {
	multiplier := FindFractionalMultiplier(s.Crystal.Frequency(), frequency)

	s.PLL[pll].SetupMultiplier(multiplier)
	s.PLL[pll].Reset()

	return multiplier.Multiply(s.Crystal.Frequency()), s.bus.Err()
}

// SetOutputFrequency sets the given output to the closest possible value of the given frequency that can be
// generated with the PLL the output is associated with.
// The method returns the effective output frequency.
func (s *Si5351) SetOutputFrequency(output OutputIndex, frequency Frequency) (Frequency, error) {
	if int(output) >= len(s.FractionalOutput) {
		return 0, errors.New("only CLK0-CLK5 are currently supported")
	}

	o := s.FractionalOutput[output]
	pllFrequency := s.PLL[o.PLL].Multiplier.Multiply(s.Crystal.Frequency())
	divider := FindFractionalDivider(pllFrequency, frequency)
	o.SetupDivider(divider)

	return divider.Divide(pllFrequency), s.bus.Err()
}

// SetOutputDivider sets the divider of the given output.
// The method returns the effective output frequency.
func (s *Si5351) SetOutputDivider(output OutputIndex, a, b, c int) (Frequency, error) {
	if int(output) >= len(s.FractionalOutput) {
		return 0, errors.New("only CLK0-CLK5 are currently supported")
	}

	o := s.FractionalOutput[output]
	pllFrequency := s.PLL[o.PLL].Multiplier.Multiply(s.Crystal.Frequency())
	divider := FractionalRatio{A: a, B: b, C: c}
	o.SetupDivider(divider)

	return divider.Divide(pllFrequency), s.bus.Err()
}

// SetupQuadratureOutput sets up the given PLL and the given outputs to generate the closest possible value
// of the given frequency with a quadrature signal (90° phase shifted) on the second output.
// The method returns the effective PLL frequency and the effective output frequency.
func (s *Si5351) SetupQuadratureOutput(pll PLLIndex, phase, quadrature OutputIndex, frequency Frequency) (Frequency, Frequency, error) {
	if int(phase) >= len(s.FractionalOutput) || int(quadrature) >= len(s.FractionalOutput) {
		return 0, 0, errors.New("only CLK0-CLK5 are currently supported")
	}

	p := s.PLL[pll]
	i := s.FractionalOutput[phase]
	q := s.FractionalOutput[quadrature]

	// Find the multiplier and an integer divider.
	multiplier, divider := FindFractionalMultiplierWithIntegerDivider(s.Crystal.Frequency(), frequency)
	pllFrequency := multiplier.Multiply(s.Crystal.Frequency())
	outputFrequency := divider.Divide(pllFrequency)

	i.SetPLL(pll)
	q.SetPLL(pll)
	p.SetupMultiplier(multiplier)
	i.SetupDivider(divider)
	q.SetupDivider(divider)

	shift := uint8(divider.A & 0xFF)
	i.SetupPhaseShift(0)
	q.SetupPhaseShift(shift)

	p.Reset()

	return pllFrequency, outputFrequency, s.bus.Err()
}

// Shutdown the Si5351: disable all outputs, power down all output drivers.
func (s *Si5351) Shutdown() error {
	s.enableAllOutputs(false)
	s.powerDownAllOutputDrivers()
	return s.bus.Err()
}

func (s *Si5351) enableAllOutputs(enabled bool) error {
	var value byte
	if enabled {
		value = 0x00
	} else {
		value = 0xFF
	}
	_, err := s.bus.WriteReg(RegOutputEnableControl, value)
	return err
}

func (s *Si5351) powerDownAllOutputDrivers() error {
	// for all clocks: power down, fractional division mode, PLLA, not inverted, Multisynth, 2mA
	_, err := s.bus.WriteReg(RegClk0Control,
		0x80,
		0x80,
		0x80,
		0x80,
		0x80,
		0x80,
		0x80,
		0x80,
	)
	return err
}

func (s *Si5351) resetAllPLLs() error {
	value := byte((1 << 7) | (1 << 5))
	s.bus.WriteReg(RegPLLReset, value)
	return s.bus.Err()
}
