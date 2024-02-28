package version

import (
	"testing"
)

func TestLTv2(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"v1 < v2", args{"1.0.0", "2.0.0"}, true},
		{"v1 = v2", args{"1.0.0", "1.0.0"}, false},
		{"v1 > v2", args{"2.0.0", "1.0.0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LTv2(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("LTv2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGTv2(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"v1 < v2", args{"1.0.0", "2.0.0"}, false},
		{"v1 = v2", args{"1.0.0", "1.0.0"}, false},
		{"v1 > v2", args{"2.0.0", "1.0.0"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GTv2(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("GTv2() = %v, want %v", got, tt.want)
			}
		})
	}
}
