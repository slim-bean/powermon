package sigineer

import (
	"fmt"
	"strconv"
	"strings"
)

type Signeer struct {
	VoltIn          float64 `json:"v_in"`
	VoltFault       float64 `json:"v_fault"`
	VoltOut         float64 `json:"v_out"`
	CurrentPercent  float64 `json:"i_pct"`
	Frequency       float64 `json:"freq"`
	VoltBatt        float64 `json:"v_batt"`
	Temp            float64 `json:"temp"`
	UtilityFail     bool    `json:"util_fail"`
	BattLow         bool    `json:"batt_low"`
	AVR             bool    `json:"avr"`
	UPSFail         bool    `json:"ups_fail"`
	LineInteractive bool    `json:"line_int"`
	Testing         bool    `json:"testing"`
	Shutdown        bool    `json:"shutdown"`
}

func Parse(in string) (*Signeer, error) {
	s := &Signeer{}
	if len(in) == 0 {
		return nil, fmt.Errorf("empty string")
	}
	if in[0:1] != "(" {
		return nil, fmt.Errorf("expected start char to be ( instead it was %s", in[0:1])
	}
	in = in[1:]
	inSplit := strings.Split(in, " ")

	if len(inSplit) < 8 {
		return nil, fmt.Errorf("string had too few pars, expected 8, found %d", len(inSplit))
	}

	f, err := strconv.ParseFloat(inSplit[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[0])
	}
	s.VoltIn = f

	f, err = strconv.ParseFloat(inSplit[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[1])
	}
	s.VoltFault = f

	f, err = strconv.ParseFloat(inSplit[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[2])
	}
	s.VoltOut = f

	f, err = strconv.ParseFloat(inSplit[3], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[3])
	}
	s.CurrentPercent = f

	f, err = strconv.ParseFloat(inSplit[4], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[4])
	}
	s.Frequency = f

	f, err = strconv.ParseFloat(inSplit[5], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[5])
	}
	s.VoltBatt = f

	f, err = strconv.ParseFloat(inSplit[6], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert string to float: %s", inSplit[6])
	}
	s.Temp = f

	i, err := strconv.ParseUint(strings.TrimSuffix(inSplit[7], "\r"), 2, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the packet, failed to convert the binary string to int: %s", inSplit[7])
	}
	if testBinary(i, 0b10000000) {
		s.UtilityFail = true
	}
	if testBinary(i, 0b01000000) {
		s.BattLow = true
	}
	if testBinary(i, 0b00100000) {
		s.AVR = true
	}
	if testBinary(i, 0b00010000) {
		s.UPSFail = true
	}
	if testBinary(i, 0b00001000) {
		s.LineInteractive = true
	}
	if testBinary(i, 0b00000100) {
		s.Testing = true
	}
	if testBinary(i, 0b00000010) {
		s.Shutdown = true
	}
	return s, nil
}

func testBinary(i, mask uint64) bool {
	return i&mask == mask
}
