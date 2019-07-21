package si5351

import (
	"io"
)

// Frequency represents a frequency in Hz
type Frequency uint32

// Frequency multipliers
const (
	Hz  Frequency = 1
	KHz Frequency = 1000
	MHz Frequency = 1000000
)

// ClockDivider represents a clock divider used at several places to divide a clock by a multiple of two.
type ClockDivider uint8

// The clock dividers.
const (
	ClockBy1 ClockDivider = iota
	ClockBy2
	ClockBy4
	ClockBy8
	ClockBy16
	ClockBy32
	ClockBy64
	ClockBy128
)

// Factor returns the factor of this divider in the denominator.
func (d ClockDivider) Factor() uint8 {
	return (1 << uint(d))
}

// FractionalRatio represents a fractional ratio used to configure the PLLs and the Multisynths.
type FractionalRatio struct {
	A            uint32
	B            uint32
	C            uint32
	ClockDivider ClockDivider
	By4          bool
}

// Encode produces the representation of the three parameters that represent the divider in the Si5351's registers.
func (d *FractionalRatio) Encode() (p1, p2, p3 uint32) {
	var fraction uint32
	if d.C == 0 {
		fraction = 0
		p3 = 1
	} else {
		fraction = uint32(128.0 * (float64(d.B) / float64(d.C)))
		p3 = d.C
	}
	p1 = 128*d.A + fraction - 512
	p2 = 128*d.B - d.C*fraction
	return
}

// Multiply this ration with the given base frequency.
func (d *FractionalRatio) Multiply(base Frequency) Frequency {
	if d.C == 0 {
		return base * Frequency(d.A)
	}
	return Frequency(float64(base) * (float64(d.A) + float64(d.B)/float64(d.C)))
}

// Divide the given frequency by this ratio.
func (d *FractionalRatio) Divide(base Frequency) Frequency {
	if d.C == 0 {
		return base / (Frequency(d.A) * Frequency(d.ClockDivider.Factor()))
	}
	return Frequency(float64(base) / ((float64(d.A) + (float64(d.B) / float64(d.C))) * float64(d.ClockDivider.Factor())))
}

// IsInteger indicates if this divider can be used in integer mode.
func (d *FractionalRatio) IsInteger() bool {
	return (d.A%2 == 0) && (d.B == 0)
}

// Bytes returns the representation of this divider in the Si5351's registers as bytes.
func (d *FractionalRatio) Bytes() []byte {
	p1, p2, p3 := d.Encode()
	bytes := []byte{
		byte((p3 & 0x0000FF00) >> 8),
		byte(p3 & 0x000000FF),
		byte((p1&0x00030000)>>16) | byte(d.ClockDivider<<4),
		byte((p1 & 0x0000FF00) >> 8),
		byte(p1 & 0x000000FF),
		byte((p3&0x000F0000)>>12) | byte((p2&0x000F0000)>>16),
		byte((p2 & 0x0000FF00) >> 8),
		byte(p2 & 0x000000FF),
	}
	if d.By4 {
		bytes[2] |= 0xC
	}
	return bytes
}

// WriteTo writes the register representation to the given writer.
func (d *FractionalRatio) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(d.Bytes())
	return int64(n), err
}

// FindFractionalMultiplier calculates a fractional ratio that allows to generate the given frequency from the given reference frequency.
func FindFractionalMultiplier(refFrequency, frequency Frequency) FractionalRatio {
	const (
		minA, maxA   = 15, 90
		defaultDenom = 0xFFFFF // 1048575 // (1 << 21) - 1 // 2000000
	)

	a := uint32(frequency / refFrequency)
	if a < minA {
		a = minA
	} else if a > maxA {
		a = maxA
	}
	b := uint32(float64(frequency%refFrequency) * (float64(defaultDenom) / float64(refFrequency)))
	c := uint32(defaultDenom)

	return FractionalRatio{A: a, B: b, C: c}
}

// FindFractionalDivider calculates a fractional ration that allows to generate the given frequency from the given reference frequency.
func FindFractionalDivider(refFrequency Frequency, frequency Frequency) FractionalRatio {
	const (
		minA, maxA   = 6, 1800
		defaultDenom = 0xFFFFF // 1048575 // (1<<21) - 1 // 2000000
	)

	a := uint32(refFrequency / frequency)
	if a < minA {
		a = minA
	} else if a > maxA {
		a = maxA
	}
	b := uint32(float64(refFrequency%frequency) * (float64(defaultDenom) / float64(frequency)))
	c := uint32(defaultDenom)

	return FractionalRatio{A: a, B: b, C: c}
}

// FindFractionalMultiplierWithIntegerDivider calculates a pair of ratios, where the divider is integer.
func FindFractionalMultiplierWithIntegerDivider(refFrequency Frequency, frequency Frequency) (multiplier, divider FractionalRatio) {
	const (
		minPLLFreq, maxPLLFreq = 600 * MHz, 900 * MHz
		maxFreq                = 150 * MHz
		minA, maxA             = 6, 126
	)

	pllFrequency := minPLLFreq
	a := (pllFrequency / frequency)

	for pllFrequency%frequency != 0 {
		a++
		pllFrequency = frequency * a
	}
	for (a%2 == 1) || (a < minA) {
		a++
	}

	clockDivider := ClockBy1
	for a > maxA {
		a = a >> 1
		clockDivider++
	}

	pllFrequency = frequency * a * Frequency(clockDivider.Factor())
	multiplier = FindFractionalMultiplier(refFrequency, pllFrequency)

	divider = FractionalRatio{A: uint32(a), B: 0, C: 1, ClockDivider: clockDivider}
	if clockDivider == ClockBy4 {
		divider.By4 = true
	}

	return
}
