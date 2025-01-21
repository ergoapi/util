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
		{"v1 = v2", args{"1.0.0", "1.0.0"}, true},
		{"v1 > v2", args{"2.0.0", "1.0.0"}, false},
		{"v1 == v2", args{"2023.10.4", "2023.10.5"}, true},
		{"v1 typo", args{"1.0.0", ""}, false},
		{"v2 typo", args{"", "1.0.0"}, true},
		{"include v", args{"1.0.0", "v1.0.1"}, true},
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

func TestNext(t *testing.T) {
	type args struct {
		now   string
		major bool
		minor bool
		patch bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"next major", args{"1.0.0", true, false, false}, "2.0.0"},
		{"next minor", args{"1.0.0", false, true, false}, "1.1.0"},
		{"next patch", args{"1.0.0", false, false, true}, "1.0.1"},
		{"next major with prefix", args{"v1.0.0", true, false, false}, "v2.0.0"},
		{"next minor with prefix", args{"v1.0.0", false, true, false}, "v1.1.0"},
		{"next patch with prefix", args{"v1.0.0", false, false, true}, "v1.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Next(tt.args.now, tt.args.major, tt.args.minor, tt.args.patch); got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
