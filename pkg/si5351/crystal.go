package si5351

// The standard crystal frequencies.
const (
	Crystal25MHz = 25 * MHz
	Crystal27MHz = 27 * MHz
)

// Crystal represents the reference Crystal of the si5351.
type Crystal struct {
	BaseFrequency Frequency
	Load          CrystalLoad
	CorrectionPPM int
}

// CrystalLoad represents the capacitve load of the Crystal.
type CrystalLoad byte

// The crystal load indicators.
const (
	CrystalLoad6PF  CrystalLoad = (1 << 6)
	CrystalLoad8PF  CrystalLoad = (2 << 6)
	CrystalLoad10PF CrystalLoad = (3 << 6)
)

// Frequency is the corrected frequency of this Crystal.
func (c Crystal) Frequency() Frequency {
	return Frequency(float64(c.BaseFrequency) + ((float64(c.CorrectionPPM) / 1000000.0) * float64(c.BaseFrequency)))
}
