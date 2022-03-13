package eg4

import "fmt"

type State int64

const (
	IDLE State = iota
	CHARGING
	DISCHARGING
)

func (s State) String() string {
	switch s {
	case IDLE:
		return "idle"
	case CHARGING:
		return "charging"
	case DISCHARGING:
		return "discharging"
	default:
		return fmt.Sprintf("unknown state '%d'", s)
	}
}

type EG4 struct {
	Cell1  float64 `json:"cell_1"`
	Cell2  float64 `json:"cell_2"`
	Cell3  float64 `json:"cell_3"`
	Cell4  float64 `json:"cell_4"`
	Cell5  float64 `json:"cell_5"`
	Cell6  float64 `json:"cell_6"`
	Cell7  float64 `json:"cell_7"`
	Cell8  float64 `json:"cell_8"`
	Cell9  float64 `json:"cell_9"`
	Cell10 float64 `json:"cell_10"`
	Cell11 float64 `json:"cell_11"`
	Cell12 float64 `json:"cell_12"`
	Cell13 float64 `json:"cell_13"`
	Cell14 float64 `json:"cell_14"`
	Cell15 float64 `json:"cell_15"`
	Cell16 float64 `json:"cell_16"`

	BatteryAmps   float64 `json:"batt_a"`
	BatteryVolts  float64 `json:"batt_v"`
	BatterySOC    float64 `json:"batt_soc"`
	BatteryCycles int64   `json:"batt_cycles"`
	BatterySOH    float64 `json:"battery_soh"`

	Temp1   float64 `json:"temp_1"`
	Temp2   float64 `json:"temp_2"`
	Temp3   float64 `json:"temp_3"`
	Temp4   float64 `json:"temp_4"`
	MOSTemp float64 `json:"mos_temp"`
	EnvTemp float64 `json:"env_temp"`

	State State `json:"state"`
}

func Parse(b []byte) (*EG4, error) {
	eg4 := &EG4{}
	g := 4
	l := 5
	gb := 0x01
	lb := 0x10
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.Cell1 = convertCell(b[6], b[7])
	eg4.Cell2 = convertCell(b[8], b[9])
	eg4.Cell3 = convertCell(b[10], b[11])
	eg4.Cell4 = convertCell(b[12], b[13])
	eg4.Cell5 = convertCell(b[14], b[15])
	eg4.Cell6 = convertCell(b[16], b[17])
	eg4.Cell7 = convertCell(b[18], b[19])
	eg4.Cell8 = convertCell(b[20], b[21])
	eg4.Cell9 = convertCell(b[22], b[23])
	eg4.Cell10 = convertCell(b[24], b[25])
	eg4.Cell11 = convertCell(b[26], b[27])
	eg4.Cell12 = convertCell(b[28], b[29])
	eg4.Cell13 = convertCell(b[30], b[31])
	eg4.Cell14 = convertCell(b[32], b[33])
	eg4.Cell15 = convertCell(b[34], b[35])
	eg4.Cell16 = convertCell(b[36], b[37])

	g = 38
	l = 39
	gb = 0x02
	lb = 0x01
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.BatteryAmps = convertAmps(b[40], b[41])

	g = 42
	l = 43
	gb = 0x03
	lb = 0x01
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.BatterySOC = convertSOC(b[44], b[45])

	g = 50
	l = 51
	gb = 0x05
	lb = 0x06
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.Temp1 = convertTemp(b[53])
	eg4.Temp2 = convertTemp(b[55])
	eg4.Temp3 = convertTemp(b[57])
	eg4.Temp4 = convertTemp(b[59])
	eg4.MOSTemp = convertTemp(b[61])
	eg4.EnvTemp = convertTemp(b[63])

	g = 64
	l = 65
	gb = 0x06
	lb = 0x05
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.State = convertState(b[69])

	g = 76
	l = 77
	gb = 0x07
	lb = 0x01
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.BatteryCycles = convertCycles(b[78], b[79])

	g = 80
	l = 81
	gb = 0x08
	lb = 0x01
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.BatteryVolts = convertVolts(b[82], b[83])

	g = 84
	l = 85
	gb = 0x09
	lb = 0x01
	if b[g] != byte(gb) || b[l] != byte(lb) {
		return eg4, fmt.Errorf("expected b[%d] to be %d and b[%d] to be %d and instead they were 0x%02X and 0x%02X", g, gb, l, lb, b[g], b[l])
	}
	eg4.BatterySOH = convertSOH(b[86], b[87])

	return eg4, nil
}

func convertCell(high, low byte) float64 {
	return float64(uint16(high)<<8|(uint16(low))) / 1000
}

func convertAmps(high, low byte) float64 {
	return (30000 - float64(uint16(high)<<8|(uint16(low)))) / 100
}

func convertSOC(high, low byte) float64 {
	return float64(uint16(high)<<8|(uint16(low))) / 100
}

func convertTemp(low byte) float64 {
	return float64(low) - 50
}

func convertState(by byte) State {
	return State(by)
}

func convertCycles(high, low byte) int64 {
	return int64(uint16(high)<<8 | (uint16(low)))
}

func convertVolts(high, low byte) float64 {
	return float64(uint16(high)<<8|(uint16(low))) / 100
}

func convertSOH(high, low byte) float64 {
	return float64(uint16(high)<<8|(uint16(low))) / 100
}
