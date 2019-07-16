package cmd

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/ftl/si5351/pkg/si5351"
)

func parseFrequency(f string) (si5351.Frequency, error) {
	input := strings.ToLower(strings.TrimSpace(f))
	var magnitude si5351.Frequency
	switch {
	case strings.HasSuffix(input, "m"):
		magnitude = si5351.MHz
		input = input[:len(input)-1]
	case strings.HasSuffix(input, "k"):
		magnitude = si5351.KHz
		input = input[:len(input)-1]
	default:
		magnitude = si5351.Hz
	}
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return si5351.Frequency(value) * magnitude, nil
}

func parseRatio(s string) (a, b, c int, err error) {
	values := strings.Split(s, ",")
	if len(values) != 3 {
		err = errors.New("ratio must have three components: a,b,c")
	}
	if err == nil {
		a, err = strconv.Atoi(strings.TrimSpace(values[0]))
	}
	if err == nil {
		b, err = strconv.Atoi(strings.TrimSpace(values[1]))
	}
	if err == nil {
		c, err = strconv.Atoi(strings.TrimSpace(values[2]))
	}
	return
}

func parsePLL(s string) (si5351.PLLIndex, error) {
	switch strings.ToUpper(s) {
	case "A":
		return si5351.PLLA, nil
	case "B":
		return si5351.PLLB, nil
	default:
		return si5351.PLLA, errors.Errorf("invalid PLL %s, try A or B", s)
	}
}

func parseOutput(s string) (si5351.OutputIndex, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < int(si5351.Clk0) || i > int(si5351.Clk5) {
		return 0, errors.Errorf("invalid output %s, only outputs 0-5 supported", s)
	}
	return si5351.OutputIndex(i), nil
}

func toCrystalFrequency(f int) si5351.Frequency {
	switch f {
	case 27:
		return si5351.Crystal27MHz
	default:
		return si5351.Crystal25MHz
	}
}

func toCrystalLoad(l int) si5351.CrystalLoad {
	switch l {
	case 6:
		return si5351.CrystalLoad6PF
	case 8:
		return si5351.CrystalLoad8PF
	default:
		return si5351.CrystalLoad10PF
	}
}

func toOutputDrive(d int) si5351.OutputDrive {
	switch d {
	case 4:
		return si5351.OutputDrive4mA
	case 6:
		return si5351.OutputDrive6mA
	case 8:
		return si5351.OutputDrive8mA
	default:
		return si5351.OutputDrive2mA
	}
}
