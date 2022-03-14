package eg4

import (
	"reflect"
	"testing"
)

func Test_convertCell(t *testing.T) {
	type args struct {
		high byte
		low  byte
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "known good",
			args: args{
				high: 0x0C,
				low:  0xF3,
			},
			want: 3.315,
		},
		{
			name: "ignore high byte MSB",
			args: args{
				high: 0x8D,
				low:  0x9A,
			},
			want: 3.482,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertCell(tt.args.high, tt.args.low); got != tt.want {
				t.Errorf("convertCell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertAmps(t *testing.T) {
	type args struct {
		high byte
		low  byte
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "zero",
			args: args{
				high: 0x75,
				low:  0x30,
			},
			want: 0.0,
		},
		{
			name: "positive",
			args: args{
				high: 0x6E,
				low:  0x4C,
			},
			want: 17.64,
		},
		{
			name: "negative",
			args: args{
				high: 0x84,
				low:  0xA7,
			},
			want: -39.59,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertAmps(tt.args.high, tt.args.low); got != tt.want {
				t.Errorf("convertAmps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertSOC(t *testing.T) {
	type args struct {
		high byte
		low  byte
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "known good",
			args: args{
				high: 0x23,
				low:  0x14,
			},
			want: 89.80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSOC(tt.args.high, tt.args.low); got != tt.want {
				t.Errorf("convertSOC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertTemp(t *testing.T) {
	type args struct {
		low byte
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "known good",
			args: args{
				low: 0x43,
			},
			want: 17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTemp(tt.args.low); got != tt.want {
				t.Errorf("convertTemp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertState(t *testing.T) {
	type args struct {
		by byte
	}
	tests := []struct {
		name string
		args args
		want State
	}{
		{
			name: "idle",
			args: args{
				by: 0,
			},
			want: IDLE,
		},
		{
			name: "charging",
			args: args{
				by: 0x01,
			},
			want: CHARGING,
		},
		{
			name: "discharging",
			args: args{
				by: 0x02,
			},
			want: DISCHARGING,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertState(tt.args.by); got != tt.want {
				t.Errorf("convertState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertCycles(t *testing.T) {
	type args struct {
		high byte
		low  byte
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "cycle 1",
			args: args{
				high: 0x00,
				low:  0x01,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertCycles(tt.args.high, tt.args.low); got != tt.want {
				t.Errorf("convertCycles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertVolts(t *testing.T) {
	type args struct {
		high byte
		low  byte
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "known good",
			args: args{
				high: 0x15,
				low:  0x62,
			},
			want: 54.74,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertVolts(tt.args.high, tt.args.low); got != tt.want {
				t.Errorf("convertVolts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertSOH(t *testing.T) {
	type args struct {
		high byte
		low  byte
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "one hundred percent",
			args: args{
				high: 0x27,
				low:  0x10,
			},
			want: 100.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSOH(tt.args.high, tt.args.low); got != tt.want {
				t.Errorf("convertSOH() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *EG4
		wantErr bool
	}{
		{
			name: "known good",
			args: args{
				b: []byte{0x7E, 0x01, 0x01, 0x58, 0x01, 0x10, 0x0C, 0xB0, 0x0C, 0xA6, 0x0C, 0xA8, 0x0C, 0x95, 0x0C, 0x98, 0x0C, 0xB0, 0x0C, 0x94, 0x0C, 0xAF, 0x0C, 0x95, 0x0C, 0x96, 0x0C, 0xA8, 0x0C, 0x92, 0x0C, 0xA7, 0x0C, 0x94, 0x0C, 0xA7, 0x0C, 0x9F, 0x02, 0x01, 0x84, 0xA7, 0x03, 0x01, 0x23, 0x0A, 0x04, 0x01, 0x27, 0x10, 0x05, 0x06, 0x00, 0x41, 0x00, 0x41, 0x00, 0x41, 0x00, 0x40, 0x80, 0x43, 0x20, 0x42, 0x06, 0x05, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x01, 0x00, 0x01, 0x08, 0x01, 0x14, 0x33, 0x09, 0x01, 0x27, 0x10, 0x0A, 0x01, 0x00, 0x00, 0x96, 0x0D},
			},
			want: &EG4{
				Cell1:         3.248,
				Cell2:         3.238,
				Cell3:         3.240,
				Cell4:         3.221,
				Cell5:         3.224,
				Cell6:         3.248,
				Cell7:         3.220,
				Cell8:         3.247,
				Cell9:         3.221,
				Cell10:        3.222,
				Cell11:        3.240,
				Cell12:        3.218,
				Cell13:        3.239,
				Cell14:        3.220,
				Cell15:        3.239,
				Cell16:        3.231,
				BatteryAmps:   -39.59,
				BatteryVolts:  51.71,
				BatterySOC:    89.70,
				BatteryCycles: 1,
				BatterySOH:    100.00,
				Temp1:         15,
				Temp2:         15,
				Temp3:         15,
				Temp4:         14,
				MOSTemp:       17,
				EnvTemp:       16,
				State:         DISCHARGING,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
