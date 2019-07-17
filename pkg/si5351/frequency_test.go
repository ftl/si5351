package si5351

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFractionalRatioBy4(t *testing.T) {
	div := FractionalRatio{By4: true}
	bytes := div.Bytes()
	assert.Equal(t, byte(0xC), bytes[2]&0xC)

	div.By4 = false
	bytes = div.Bytes()
	assert.Equal(t, byte(0), bytes[2]&0xC)
}

func TestFractionalRatioClockDivider(t *testing.T) {
	div := FractionalRatio{ClockDivider: ClockBy16}
	bytes := div.Bytes()
	assert.Equal(t, byte(0x40), bytes[2]&0x70)
}

func TestFindFractionalMultiplier(t *testing.T) {
	crystal := Crystal{BaseFrequency: Crystal25MHz, CorrectionPPM: 30}
	for f := 600; f <= 900; f++ {
		frequency := Frequency(f) * MHz
		t.Run(fmt.Sprintf("%d", frequency), func(t *testing.T) {
			t.Parallel()
			multiplier := FindFractionalMultiplier(crystal.Frequency(), frequency)
			actual := multiplier.Multiply(crystal.Frequency())
			assert.True(t, math.Abs(float64(frequency-actual)) < 13)
		})
	}
}

func TestFindFractionalDivider(t *testing.T) {
	pllFrequency := 900 * MHz
	for f := 1; f <= 150; f++ {
		frequency := Frequency(f) * MHz
		t.Run(fmt.Sprintf("%d", frequency), func(t *testing.T) {
			t.Parallel()
			divider := FindFractionalDivider(pllFrequency, frequency)
			actual := divider.Divide(pllFrequency)
			assert.True(t, math.Abs(float64(frequency-actual)) < 9)
		})
	}
}

func TestFindFractionalMultiplierWithIntegerDivider(t *testing.T) {
	refFrequency := 25 * MHz
	for f := 1000; f <= 150000; f += 10 {
		frequency := Frequency(f) * KHz
		t.Run(fmt.Sprintf("%d", frequency), func(t *testing.T) {
			t.Parallel()
			multiplier, divider := FindFractionalMultiplierWithIntegerDivider(refFrequency, frequency)
			actual := divider.Divide(multiplier.Multiply(refFrequency))
			assert.True(t, math.Abs(float64(frequency-actual)) < 3, "", actual)
		})
	}
}
