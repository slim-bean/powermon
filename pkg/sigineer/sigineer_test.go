package sigineer

import (
	"strconv"
	"testing"
)

func Test_testBinary(t *testing.T) {
	type args struct {
		i    uint64
		mask uint64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test match",
			args: args{
				i:    0x80,
				mask: 0b10000000,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testBinary(tt.args.i, tt.args.mask); got != tt.want {
				t.Errorf("testBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringInt(t *testing.T) {
	i, err := strconv.ParseUint("00001001", 2, 64)
	if err != nil {
		t.Fatal(err)
	}
	if i != 0b00001001 {
		t.Fatal("failed to parse binary string properly")
	}
}
