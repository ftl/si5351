package si5351

// PLLIndex indicates one of both PLLs.
type PLLIndex int

// The PLLs.
const (
	PLLA PLLIndex = iota
	PLLB
)

// PLL represents a PLL of the Si5351.
type PLL struct {
	Register    PLLRegister
	InputSource PLLInputSource
	Multiplier  FractionalRatio

	bus Bus
}

// PLLInputSource describes the input source of a PLL.
type PLLInputSource uint8

// All possible input sources of a PLL.
const (
	PLLInputCrystal PLLInputSource = iota
	PLLInputClkin
)

// PLLRegister describes the registers used by a PLL.
type PLLRegister struct {
	Multiplier        uint8
	ResetOffset       uint8
	InputSourceOffset uint8
}

// PLLRegisters contains the register descriptions of all PLLs.
var PLLRegisters = []PLLRegister{
	{RegPLLAMultisynthParameters, 5, 2},
	{RegPLLBMultisynthParameters, 7, 3},
}

func loadPLLs(bus Bus) []*PLL {
	result := make([]*PLL, len(PLLRegisters))
	for i, register := range PLLRegisters {
		result[i] = &PLL{
			Register: register,
			bus:      bus,
		}
	}
	return result
}

// SetupMultiplier writes the frequency multiplier into the registers and resets the PLL.
func (p *PLL) SetupMultiplier(multiplier FractionalRatio) error {
	multiplier.WriteTo(p.bus.RegWriter(p.Register.Multiplier))

	if p.bus.Err() == nil {
		p.Multiplier = multiplier
	}
	return p.bus.Err()
}

// Reset the PLL.
func (p *PLL) Reset() error {
	p.bus.WriteReg(RegPLLReset, (1 << p.Register.ResetOffset))
	return p.bus.Err()
}
