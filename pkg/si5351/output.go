package si5351

// OutputIndex indicates one of the output clocks.
type OutputIndex int

// The output clocks.
const (
	Clk0 OutputIndex = iota
	Clk1
	Clk2
	Clk3
	Clk4
	Clk5
	Clk6
	Clk7
)

// Output describes the properties common to all of the Si5351's output clocks.
type Output struct {
	Register    OutputRegister
	PowerDown   bool
	IntegerMode bool
	Invert      bool
	PLL         PLLIndex
	InputSource ClockInputSource
	Drive       OutputDrive

	bus Bus
}

// FractionalOutput represents an output that has a fractional frequency divider (CLK0-CLK5).
// Those outputs can also have a phase shift.
type FractionalOutput struct {
	Output
	FrequencyDivider FractionalRatio
	PhaseShift       uint8
}

// IntegerOutput represents an output that has an integer frequency divider (CLK6-CLK7).
type IntegerOutput struct {
	Output
	FrequencyDivider uint8
	RDiv             ClockDivider
}

// ClockInputSource describes the input source of an output clock.
type ClockInputSource uint8

// All possible input sources for an output clock.
const (
	ClockInputCrystal ClockInputSource = iota
	ClockInputClkin
	ClockInputReserved
	ClockInputMultisynth
)

// OutputDrive describes the drive strength of an Output.
type OutputDrive uint8

// All possible drive strengthes of an Output.
const (
	OutputDrive2mA OutputDrive = iota
	OutputDrive4mA
	OutputDrive6mA
	OutputDrive8mA
)

// OutputDisableState describes the state of an Output when it is disabled.
type OutputDisableState uint8

// All possible states a Clock can have when disabled.
const (
	OutputDisableLow OutputDisableState = iota
	OutputDisableHigh
	OutputDisableHighZ
	OutputDisableNever
)

// OutputRegister describes the registers used by a Clock.
type OutputRegister struct {
	Control            uint8
	DisableState       uint8
	DisableStateOffset uint8
	PhaseShift         uint8
	Divider            uint8
	DividerOffset      uint8
}

// FractionalOutputRegisters contains the register descriptions for all outputs with fractional dividers.
var FractionalOutputRegisters = []OutputRegister{
	{RegClk0Control, RegClk3_0DisableState, 0, RegClk0InitialPhaseOffset, RegMultisynth0Parameters, 0},
	{RegClk1Control, RegClk3_0DisableState, 2, RegClk1InitialPhaseOffset, RegMultisynth1Parameters, 0},
	{RegClk2Control, RegClk3_0DisableState, 4, RegClk2InitialPhaseOffset, RegMultisynth2Parameters, 0},
	{RegClk3Control, RegClk3_0DisableState, 6, RegClk3InitialPhaseOffset, RegMultisynth3Parameters, 0},
	{RegClk4Control, RegClk7_4DisableState, 0, RegClk4InitialPhaseOffset, RegMultisynth4Parameters, 0},
	{RegClk5Control, RegClk7_4DisableState, 2, RegClk5InitialPhaseOffset, RegMultisynth5Parameters, 0},
}

func loadFractionalOutputs(bus Bus) []*FractionalOutput {
	result := make([]*FractionalOutput, len(FractionalOutputRegisters))
	for i, register := range FractionalOutputRegisters {
		result[i] = &FractionalOutput{
			Output: Output{
				Register: register,
				bus:      bus,
			},
		}
	}
	return result
}

// IntegerOutputRegisters contains the register descriptions for all outputs with integer dividers.
var IntegerOutputRegisters = []OutputRegister{
	{RegClk6Control, RegClk7_4DisableState, 4, 0, RegMultisynth6Parameters, 0},
	{RegClk7Control, RegClk7_4DisableState, 6, 0, RegMultisynth7Parameters, 4},
}

func loadIntegerOutputs(bus Bus) []*IntegerOutput {
	result := make([]*IntegerOutput, len(IntegerOutputRegisters))
	for i, register := range IntegerOutputRegisters {
		result[i] = &IntegerOutput{
			Output: Output{
				Register: register,
				bus:      bus,
			},
		}
	}
	return result
}

// SetupControl writes the control register of the Output.
func (o *Output) SetupControl(powerDown bool, integerMode bool, pll PLLIndex, invert bool, inputSource ClockInputSource, drive OutputDrive) error {
	value := byte(pll<<5) | byte(inputSource<<2) | byte(drive)
	if powerDown {
		value |= (1 << 7)
	}
	if integerMode {
		value |= (1 << 6)
	}
	if invert {
		value |= (1 << 4)
	}

	o.bus.WriteReg(o.Register.Control, value)

	if o.bus.Err() == nil {
		o.PowerDown = powerDown
		o.IntegerMode = integerMode
		o.PLL = pll
		o.Invert = invert
		o.InputSource = inputSource
		o.Drive = drive
	}
	return o.bus.Err()
}

// SetPowerDown sets the power down flag of the Output and writes it to the output's control register.
func (o *Output) SetPowerDown(powerDown bool) error {
	value := byte(o.PLL<<5) | byte(o.InputSource<<2) | byte(o.Drive)
	if powerDown {
		value |= (1 << 7)
	}
	if o.IntegerMode {
		value |= (1 << 6)
	}
	if o.Invert {
		value |= (1 << 4)
	}

	o.bus.WriteReg(o.Register.Control, value)

	if o.bus.Err() == nil {
		o.PowerDown = powerDown
	}

	return o.bus.Err()
}

// SetPLL sets the PLL of the Output and writes it to the output's control register.
func (o *Output) SetPLL(pll PLLIndex) error {
	value := byte(pll<<5) | byte(o.InputSource<<2) | byte(o.Drive)
	if o.PowerDown {
		value |= (1 << 7)
	}
	if o.IntegerMode {
		value |= (1 << 6)
	}
	if o.Invert {
		value |= (1 << 4)
	}

	o.bus.WriteReg(o.Register.Control, value)

	if o.bus.Err() == nil {
		o.PLL = pll
	}

	return o.bus.Err()
}

// SetInvert sets the inversion flag of the Output and writes it to the output's control register.
func (o *Output) SetInvert(invert bool) error {
	value := byte(o.PLL<<5) | byte(o.InputSource<<2) | byte(o.Drive)
	if o.PowerDown {
		value |= (1 << 7)
	}
	if o.IntegerMode {
		value |= (1 << 6)
	}
	if invert {
		value |= (1 << 4)
	}

	o.bus.WriteReg(o.Register.Control, value)

	if o.bus.Err() == nil {
		o.Invert = invert
	}

	return o.bus.Err()
}

// SetInputSource sets the clock input source of the Output and writes it to the output's control register.
func (o *Output) SetInputSource(inputSource ClockInputSource) error {
	value := byte(o.PLL<<5) | byte(inputSource<<2) | byte(o.Drive)
	if o.PowerDown {
		value |= (1 << 7)
	}
	if o.IntegerMode {
		value |= (1 << 6)
	}
	if o.Invert {
		value |= (1 << 4)
	}

	o.bus.WriteReg(o.Register.Control, value)

	if o.bus.Err() == nil {
		o.InputSource = inputSource
	}

	return o.bus.Err()
}

// SetDrive sets the output drive strength of the Output and writes it to the output's control register.
func (o *Output) SetDrive(drive OutputDrive) error {
	value := byte(o.PLL<<5) | byte(o.InputSource<<2) | byte(drive)
	if o.PowerDown {
		value |= (1 << 7)
	}
	if o.IntegerMode {
		value |= (1 << 6)
	}
	if o.Invert {
		value |= (1 << 4)
	}

	o.bus.WriteReg(o.Register.Control, value)

	if o.bus.Err() == nil {
		o.Drive = drive
	}

	return o.bus.Err()
}

// SetupDivider writes the frequency divider into the registers.
func (o *FractionalOutput) SetupDivider(divider FractionalRatio) error {
	divider.WriteTo(o.bus.RegWriter(o.Register.Divider))

	if o.bus.Err() == nil {
		o.FrequencyDivider = divider
	}
	return o.bus.Err()
}

// SetupPhaseShift sets the phase shift of the Clock.
func (o *FractionalOutput) SetupPhaseShift(phaseShift uint8) error {
	o.bus.WriteReg(o.Register.PhaseShift, byte(phaseShift&0x7F))

	if o.bus.Err() == nil {
		o.PhaseShift = phaseShift
	}
	return o.bus.Err()
}
