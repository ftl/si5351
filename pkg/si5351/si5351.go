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

	pll              []*PLL
	fractionalOutput []*FractionalOutput
	integerOutput    []*IntegerOutput

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
		pll:              loadPLLs(bus),
		fractionalOutput: loadFractionalOutputs(bus),
		integerOutput:    loadIntegerOutputs(bus),
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

// PLL returns the PLL with the given index.
func (s *Si5351) PLL(pll PLLIndex) *PLL {
	return s.pll[pll]
}

// PLLA returns PLL A.
func (s *Si5351) PLLA() *PLL {
	return s.pll[PLLA]
}

// PLLB returns PLL B.
func (s *Si5351) PLLB() *PLL {
	return s.pll[PLLB]
}

// Output returns the output with the given index.
func (s *Si5351) Output(output OutputIndex) *Output {
	if output <= Clk5 {
		return &s.fractionalOutput[output].Output
	}
	return &s.integerOutput[output].Output
}

// Clk0 returns the output CLK0.
func (s *Si5351) Clk0() *FractionalOutput {
	return s.fractionalOutput[Clk0]
}

// Clk1 returns the output CLK1
func (s *Si5351) Clk1() *FractionalOutput {
	return s.fractionalOutput[Clk1]
}

// Clk2 returns the output CLK2
func (s *Si5351) Clk2() *FractionalOutput {
	return s.fractionalOutput[Clk2]
}

// Clk3 returns the output CLK3
func (s *Si5351) Clk3() *FractionalOutput {
	return s.fractionalOutput[Clk3]
}

// Clk4 returns the output CLK4
func (s *Si5351) Clk4() *FractionalOutput {
	return s.fractionalOutput[Clk4]
}

// Clk5 returns the output CLK5
func (s *Si5351) Clk5() *FractionalOutput {
	return s.fractionalOutput[Clk5]
}

// Clk6 returns the output CLK6
func (s *Si5351) Clk6() *IntegerOutput {
	return s.integerOutput[0]
}

// Clk7 returns the output CLK7
func (s *Si5351) Clk7() *IntegerOutput {
	return s.integerOutput[1]
}

// SetupPLLInputSource writes the input source configuration to the Si5351's register.
func (s *Si5351) SetupPLLInputSource(clkinInputDivider ClockDivider, pllASource, pllBSource PLLInputSource) error {
	value := byte((clkinInputDivider&0xF)<<4) |
		byte((pllASource&1)<<s.PLLA().Register.InputSourceOffset) |
		byte((pllBSource&1)<<s.PLLB().Register.InputSourceOffset)

	s.bus.WriteReg(RegPLLInputSource, value)

	if s.bus.Err() == nil {
		s.InputDivider = clkinInputDivider
		s.PLLA().InputSource = pllASource
		s.PLLB().InputSource = pllBSource
	}
	return s.bus.Err()
}

// SetupPLLRaw directly sets the frequency multiplier parameters for the given PLL and resets it.
func (s *Si5351) SetupPLLRaw(pll PLLIndex, a, b, c uint32) error {
	s.pll[pll].SetupMultiplier(FractionalRatio{A: a, B: b, C: c})
	s.pll[pll].Reset()
	return s.bus.Err()
}

// SetupMultisynthRaw directly sets the frequency divider and RDiv parameters for the Multisynth of the given output.
func (s *Si5351) SetupMultisynthRaw(output OutputIndex, a, b, c uint32, RDiv ClockDivider) error {
	if int(output) >= len(s.fractionalOutput) {
		return errors.New("only CLK0-CLK5 are currently supported")
	}

	s.fractionalOutput[output].SetupDivider(FractionalRatio{A: a, B: b, C: c})

	return s.bus.Err()
}

// SetupPLL sets the given PLL to the closest possible value of the given frequency and resets it.
func (s *Si5351) SetupPLL(pll PLLIndex, frequency Frequency) (Frequency, error) {
	multiplier := FindFractionalMultiplier(s.Crystal.Frequency(), frequency)

	s.pll[pll].SetupMultiplier(multiplier)
	s.pll[pll].Reset()

	return multiplier.Multiply(s.Crystal.Frequency()), s.bus.Err()
}

// PrepareOutputs prepares the given outputs for use with the given PLL and control parameters.
func (s *Si5351) PrepareOutputs(pll PLLIndex, invert bool, inputSource ClockInputSource, drive OutputDrive, outputs ...OutputIndex) error {
	for _, output := range outputs {
		s.Output(output).SetupControl(false, false, pll, invert, inputSource, drive)
	}
	return s.bus.Err()
}

// SetOutputFrequency sets the given output to the closest possible value of the given frequency that can be
// generated with the PLL the output is associated with. Set the frequency of the PLL first.
// The method returns the effective output frequency.
func (s *Si5351) SetOutputFrequency(output OutputIndex, frequency Frequency) (Frequency, error) {
	if int(output) >= len(s.fractionalOutput) {
		return 0, errors.New("only CLK0-CLK5 are currently supported")
	}

	o := s.fractionalOutput[output]
	pllFrequency := s.pll[o.PLL].Multiplier.Multiply(s.Crystal.Frequency())
	divider := FindFractionalDivider(pllFrequency, frequency)
	o.SetupDivider(divider)

	return divider.Divide(pllFrequency), s.bus.Err()
}

// SetOutputDivider sets the divider of the given output.
// The method returns the effective output frequency.
func (s *Si5351) SetOutputDivider(output OutputIndex, a, b, c uint32) (Frequency, error) {
	if int(output) >= len(s.fractionalOutput) {
		return 0, errors.New("only CLK0-CLK5 are currently supported")
	}

	o := s.fractionalOutput[output]
	pllFrequency := s.pll[o.PLL].Multiplier.Multiply(s.Crystal.Frequency())
	divider := FractionalRatio{A: a, B: b, C: c}
	o.SetupDivider(divider)

	return divider.Divide(pllFrequency), s.bus.Err()
}

// SetupQuadratureOutput sets up the given PLL and the given outputs to generate the closest possible value
// of the given frequency with a quadrature signal (90° phase shifted) on the second output.
// The method returns the effective PLL frequency and the effective output frequency.
func (s *Si5351) SetupQuadratureOutput(pll PLLIndex, phase, quadrature OutputIndex, frequency Frequency) (Frequency, Frequency, error) {
	if int(phase) >= len(s.fractionalOutput) || int(quadrature) >= len(s.fractionalOutput) {
		return 0, 0, errors.New("only CLK0-CLK5 are currently supported")
	}

	p := s.pll[pll]
	i := s.fractionalOutput[phase]
	q := s.fractionalOutput[quadrature]

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
