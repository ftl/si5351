package si5351

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrystalFrequencyCorrection(t *testing.T) {
	crystal := Crystal{BaseFrequency: Crystal25MHz, CorrectionPPM: 30}

	assert.Equal(t, Frequency(25000750), crystal.Frequency())
}
